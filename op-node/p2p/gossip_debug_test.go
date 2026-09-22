package p2p

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/golang/snappy"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/clock"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/ptr"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
)

func debugFields(args []any) map[string]any {
	fields := make(map[string]any)
	for i := 0; i < len(args); i += 2 {
		fields[args[i].(string)] = args[i+1]
	}
	return fields
}

func TestGossipDebugRPCMetadata(t *testing.T) {
	d := newGossipDebug(testlog.Logger(t, log.LevelInfo), "self")
	rpc := &pubsub.RPC{RPC: pb.RPC{
		Subscriptions: []*pb.RPC_SubOpts{{Topicid: ptr.New("blocks"), Subscribe: ptr.New(true)}},
		Publish:       []*pb.Message{{Topic: ptr.New("blocks"), Data: []byte("body-must-not-be-logged")}},
		Control: &pb.ControlMessage{
			Ihave:     []*pb.ControlIHave{{TopicID: ptr.New("blocks"), MessageIDs: []string{"\x00\xff"}}},
			Iwant:     []*pb.ControlIWant{{MessageIDs: []string{"\x00\xff"}}},
			Graft:     []*pb.ControlGraft{{TopicID: ptr.New("blocks")}},
			Prune:     []*pb.ControlPrune{{TopicID: ptr.New("blocks"), Backoff: ptr.New(uint64(90))}},
			Idontwant: []*pb.ControlIDontWant{{MessageIDs: []string{"\x00\xff"}}},
		},
	}}
	before, err := rpc.Marshal()
	require.NoError(t, err)
	require.NoError(t, d.inspectRPC("remote", rpc))
	after, err := rpc.Marshal()
	require.NoError(t, err)
	require.Equal(t, before, after, "diagnostics must not modify the RPC")
	require.EqualValues(t, 1, d.counts[debugPublishRecv].Load())
	require.Positive(t, d.lastPublish.Load())

	rec := <-d.queue
	require.Equal(t, debugRPCRecv, rec.kind)
	fields := debugFields(rec.args)
	require.Equal(t, peer.ID("remote"), fields["peer"])
	items := fields["items"].([]debugRPCItem)
	require.Len(t, items, 7)
	require.EqualValues(t, 90, *items[5].BackoffSec)
	require.Equal(t, []string{"00ff"}, items[2].IDSample)
	encoded, err := json.Marshal(fields)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "body-must-not-be-logged")

	// Reusing/modifying the RPC after the hook returns must not change queued data.
	*rpc.Control.Prune[0].Backoff = 999
	rpc.Control.Ihave[0].MessageIDs[0] = "changed"
	require.EqualValues(t, 90, *items[5].BackoffSec)
	require.Equal(t, []string{"00ff"}, items[2].IDSample)
}

func TestGossipDebugMetadataBounds(t *testing.T) {
	rpc := &pubsub.RPC{RPC: pb.RPC{Control: &pb.ControlMessage{}}}
	for range 100 {
		rpc.Publish = append(rpc.Publish, &pb.Message{Topic: ptr.New(strings.Repeat("x", 10000)), Data: make([]byte, 1024)})
		rpc.Control.Ihave = append(rpc.Control.Ihave, &pb.ControlIHave{
			MessageIDs: []string{strings.Repeat("y", 10000), "a", "b", "c", "d"},
		})
	}
	fields := debugFields(rpcFields("remote", rpc))
	require.Equal(t, 100, fields["publish"])
	items := fields["items"].([]debugRPCItem)
	require.Len(t, items, 2*gossipDebugItems)
	require.Len(t, items[0].Topic, gossipDebugTextLimit+3)
	require.Len(t, items[gossipDebugItems].IDSample, gossipDebugItems)
	require.Len(t, items[gossipDebugItems].IDSample[0], 64*2+3)
	// Optional/nil protobuf entries and an absent Control are valid to inspect.
	require.NotPanics(t, func() {
		rpcFields("remote", &pubsub.RPC{RPC: pb.RPC{Publish: []*pb.Message{nil}, Subscriptions: []*pb.RPC_SubOpts{nil}}})
	})
}

func TestGossipDebugBackpressure(t *testing.T) {
	d := newGossipDebug(testlog.Logger(t, log.LevelInfo), "self")
	d.dataLimit = rate.NewLimiter(rate.Inf, 0)
	rpc := &pubsub.RPC{RPC: pb.RPC{Publish: []*pb.Message{{Data: []byte("payload")}}}}
	const workers, perWorker = 8, 100
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for range perWorker {
				if err := d.inspectRPC("remote", rpc); err != nil {
					t.Errorf("inspector rejected RPC: %v", err)
				}
			}
		})
	}
	wg.Wait()
	require.Len(t, d.queue, gossipDebugQueueSize)
	require.EqualValues(t, workers*perWorker, d.counts[debugRPCRecv].Load())
	require.EqualValues(t, workers*perWorker, d.counts[debugPublishRecv].Load())
	require.EqualValues(t, workers*perWorker-gossipDebugQueueSize, d.queueDropped.Load())

	d.dataLimit = rate.NewLimiter(0, 0)
	require.NoError(t, d.inspectRPC("remote", rpc))
	require.EqualValues(t, 1, d.rateLimited.Load())
	require.EqualValues(t, workers*perWorker+1, d.counts[debugPublishRecv].Load())
	// Even while data is rate limited, a PRUNE can use the independent control budget.
	<-d.queue
	require.NoError(t, d.inspectRPC("remote", &pubsub.RPC{RPC: pb.RPC{
		Control: &pb.ControlMessage{Prune: []*pb.ControlPrune{{Backoff: ptr.New(uint64(60))}}},
	}}))
	require.Len(t, d.queue, gossipDebugQueueSize)
	require.EqualValues(t, 1, d.rateLimited.Load())
}

type debugBlockingWriter struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func (w *debugBlockingWriter) Write(b []byte) (int, error) {
	w.once.Do(func() { close(w.entered) })
	<-w.release
	return len(b), nil
}

func TestGossipDebugSlowLogger(t *testing.T) {
	mn := mocknet.New()
	t.Cleanup(func() { require.NoError(t, mn.Close()) })
	h, err := mn.GenPeer()
	require.NoError(t, err)
	w := &debugBlockingWriter{entered: make(chan struct{}), release: make(chan struct{})}
	d := newGossipDebug(log.NewLogger(slog.NewTextHandler(w, nil)), h.ID())
	d.start(context.Background(), h, func(peer.ID) (float64, error) { return 0, nil })
	t.Cleanup(func() {
		d.stop()
		close(w.release)
		select {
		case <-d.done:
		case <-time.After(5 * time.Second):
			t.Error("diagnostic worker did not stop")
		}
	})
	select {
	case <-w.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("diagnostic worker did not start")
	}
	returned := make(chan struct{})
	go func() {
		for range 1000 {
			_ = d.inspectRPC("remote", &pubsub.RPC{})
		}
		d.stop()
		close(returned)
	}()
	select {
	case <-returned:
	case <-time.After(5 * time.Second):
		t.Fatal("receive hook or shutdown blocked on logging")
	}
}

func TestGossipDebugValidationReasons(t *testing.T) {
	for _, tc := range []struct {
		name   string
		data   []byte
		result pubsub.ValidationResult
		reason string
	}{
		{"compression", []byte{0xff}, pubsub.ValidationReject, "invalid_snappy_length"},
		{"size", snappy.Encode(nil, []byte{0}), pubsub.ValidationReject, "undersized_payload"},
		{"signer", snappy.Encode(nil, make([]byte, 66)), pubsub.ValidationIgnore, "sequencer_address_unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := testlog.Logger(t, log.LevelCrit)
			d := newGossipDebug(logger, "self")
			base := &mockGossipSetupConfigurablesWithThreshold{threshold: time.Minute}
			cfg := &rollup.Config{L2ChainID: big.NewInt(100)}
			runCfg := &testutils.MockRuntimeConfig{}
			plain := BuildBlocksValidator(logger, cfg, runCfg, eth.BlockV4, base, clock.SystemClock)
			debug := BuildBlocksValidator(logger, cfg, runCfg, eth.BlockV4, &gossipDebugConfig{base, d}, clock.SystemClock)
			msg := &pubsub.Message{Message: &pb.Message{Data: tc.data, Topic: ptr.New("blocks")}, ID: "id"}
			require.Equal(t, tc.result, plain(context.Background(), "remote", msg))
			require.Equal(t, tc.result, debug(context.Background(), "remote", msg))
			rec := <-d.queue
			require.Equal(t, tc.reason, debugFields(rec.args)["reason"])
		})
	}
}

type debugGossipHandler func(context.Context, peer.ID, *eth.ExecutionPayloadEnvelope) error

func (fn debugGossipHandler) OnUnsafeL2Payload(ctx context.Context, from peer.ID, payload *eth.ExecutionPayloadEnvelope) error {
	return fn(ctx, from, payload)
}

func TestGossipDebugHandlerPassthrough(t *testing.T) {
	for _, expectedErr := range []error{nil, errors.New("engine unavailable")} {
		d := newGossipDebug(testlog.Logger(t, log.LevelInfo), "self")
		payload := &eth.ExecutionPayloadEnvelope{ExecutionPayload: &eth.ExecutionPayload{BlockNumber: 123}}
		ctx := context.Background()
		calls := 0
		handler := &debugGossipIn{debug: d, next: debugGossipHandler(func(gotCtx context.Context, from peer.ID, got *eth.ExecutionPayloadEnvelope) error {
			calls++
			require.Same(t, payload, got)
			require.Equal(t, ctx, gotCtx)
			require.Equal(t, peer.ID("remote"), from)
			return expectedErr
		})}
		require.Equal(t, expectedErr, handler.OnUnsafeL2Payload(ctx, "remote", payload))
		require.Equal(t, 1, calls)
		if expectedErr == nil {
			require.EqualValues(t, 1, d.counts[debugHandlerDone].Load())
			require.Positive(t, d.lastHandled.Load())
		} else {
			require.EqualValues(t, 1, d.counts[debugHandlerError].Load())
			require.Zero(t, d.lastHandled.Load())
		}
	}
	require.Equal(t, float64(-1), debugAge(time.Now(), 0))
}
