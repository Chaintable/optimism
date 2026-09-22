package p2p

import (
	"context"
	"encoding/hex"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"golang.org/x/time/rate"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/ptr"
)

const (
	gossipDebugInterval  = 30 * time.Second
	gossipDebugQueueSize = 128
	gossipDebugItems     = 4
	gossipDebugTextLimit = 128
)

type gossipDebugKind uint8

const (
	debugRPCRecv gossipDebugKind = iota
	debugRPCSend
	debugRPCDrop
	debugPublishRecv
	debugAddPeer
	debugRemovePeer
	debugJoin
	debugLeave
	debugGraft
	debugPrune
	debugValidate
	debugValidationAccept
	debugValidationReject
	debugValidationIgnore
	debugDeliver
	debugReject
	debugDuplicate
	debugThrottle
	debugUndeliverable
	debugHandlerStart
	debugHandlerDone
	debugHandlerError
	debugConnected
	debugDisconnected
	debugStreamOpen
	debugStreamOpenError
	debugStreamReadError
	debugStreamWriteError
	debugStreamClose
	debugStreamReset
	debugStreamHandlerExit
	debugKindCount
)

var gossipDebugNames = [debugKindCount]string{
	"rpc_recv", "rpc_send", "rpc_drop", "publish_recv", "peer_add", "peer_remove",
	"topic_join", "topic_leave", "mesh_graft", "mesh_prune", "validate_start",
	"validation_accept", "validation_reject", "validation_ignore", "deliver",
	"reject", "duplicate", "throttle", "undeliverable", "handler_start",
	"handler_done", "handler_error", "connected", "disconnected",
	"stream_open", "stream_open_error", "stream_read_error", "stream_write_error",
	"stream_close", "stream_reset", "stream_handler_exit",
}

type gossipDebugRecord struct {
	kind gossipDebugKind
	at   time.Time
	args []any
}

// gossipDebug never writes logs from a libp2p/validator callback. Its bounded
// queue and separate control/data budgets keep slow logging off the receive
// path. Counters and last-receive times are updated even when details are lost.
// Queued records contain only bounded copies of metadata, never message bodies.
type gossipDebug struct {
	log          log.Logger
	queue        chan gossipDebugRecord
	controlLimit *rate.Limiter
	dataLimit    *rate.Limiter
	counts       [debugKindCount]atomic.Uint64
	rateLimited  atomic.Uint64
	queueDropped atomic.Uint64
	lastPublish  atomic.Int64
	lastDeliver  atomic.Int64
	lastHandled  atomic.Int64
	streamRead   atomic.Uint64
	streamWrite  atomic.Uint64
	cancel       context.CancelFunc
	done         chan struct{}
}

var _ pubsub.RawTracer = (*gossipDebug)(nil)

type gossipDebugConfig struct {
	GossipSetupConfigurables
	debug *gossipDebug
}

func newGossipDebug(logger log.Logger, self peer.ID) *gossipDebug {
	return &gossipDebug{
		log:          logger.New("diagnostic", "gossip", "self", self),
		queue:        make(chan gossipDebugRecord, gossipDebugQueueSize),
		controlLimit: rate.NewLimiter(2, 32),
		dataLimit:    rate.NewLimiter(8, 16),
		done:         make(chan struct{}),
	}
}

func (d *gossipDebug) options() []pubsub.Option {
	return []pubsub.Option{pubsub.WithAppSpecificRpcInspector(d.inspectRPC), pubsub.WithRawTracer(d)}
}

func (d *gossipDebug) record(kind gossipDebugKind, data bool, fields func() []any) {
	d.counts[kind].Add(1)
	limiter := d.controlLimit
	if data {
		limiter = d.dataLimit
	}
	if !limiter.Allow() {
		d.rateLimited.Add(1)
		return
	}
	rec := gossipDebugRecord{kind: kind, at: time.Now(), args: fields()}
	select {
	case d.queue <- rec:
	default:
		d.queueDropped.Add(1)
	}
}

func debugText(s string) string {
	if len(s) > gossipDebugTextLimit {
		return s[:gossipDebugTextLimit] + "..."
	}
	// Do not retain the RPC's backing storage after returning from a callback.
	return strings.Clone(s)
}

func debugID(s string) string {
	if len(s) > 64 {
		return hex.EncodeToString([]byte(s[:64])) + "..."
	}
	return hex.EncodeToString([]byte(s))
}

type debugRPCItem struct {
	Type       string   `json:"type"`
	Topic      string   `json:"topic,omitempty"`
	Subscribe  *bool    `json:"subscribe,omitempty"`
	Bytes      int      `json:"bytes,omitempty"`
	IDs        int      `json:"ids,omitempty"`
	IDSample   []string `json:"id_sample,omitempty"`
	BackoffSec *uint64  `json:"backoff_seconds,omitempty"`
	Peers      int      `json:"peers,omitempty"`
}

func debugIDSample(ids []string) []string {
	out := make([]string, 0, min(len(ids), gossipDebugItems))
	for _, id := range ids[:min(len(ids), gossipDebugItems)] {
		out = append(out, debugID(id))
	}
	return out
}

// rpcFields observes decoded GossipSub RPC metadata before subscription and
// score filtering. It does not decompress, hash, retain, or log publish bodies.
func rpcFields(p peer.ID, rpc *pubsub.RPC) []any {
	c := rpc.GetControl()
	items := make([]debugRPCItem, 0, gossipDebugItems)
	for _, sub := range rpc.Subscriptions[:min(len(rpc.Subscriptions), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "subscription", Topic: debugText(sub.GetTopicid()), Subscribe: ptr.New(sub.GetSubscribe())})
	}
	for _, msg := range rpc.Publish[:min(len(rpc.Publish), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "publish", Topic: debugText(msg.GetTopic()), Bytes: len(msg.GetData())})
	}
	for _, item := range c.GetIhave()[:min(len(c.GetIhave()), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "ihave", Topic: debugText(item.GetTopicID()), IDs: len(item.GetMessageIDs()), IDSample: debugIDSample(item.GetMessageIDs())})
	}
	for _, item := range c.GetIwant()[:min(len(c.GetIwant()), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "iwant", IDs: len(item.GetMessageIDs()), IDSample: debugIDSample(item.GetMessageIDs())})
	}
	for _, item := range c.GetGraft()[:min(len(c.GetGraft()), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "graft", Topic: debugText(item.GetTopicID())})
	}
	for _, item := range c.GetPrune()[:min(len(c.GetPrune()), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "prune", Topic: debugText(item.GetTopicID()), BackoffSec: ptr.New(item.GetBackoff()), Peers: len(item.GetPeers())})
	}
	for _, item := range c.GetIdontwant()[:min(len(c.GetIdontwant()), gossipDebugItems)] {
		items = append(items, debugRPCItem{Type: "idontwant", IDs: len(item.GetMessageIDs()), IDSample: debugIDSample(item.GetMessageIDs())})
	}
	return []any{
		"peer", p, "subscriptions", len(rpc.Subscriptions), "publish", len(rpc.Publish),
		"ihave", len(c.GetIhave()), "iwant", len(c.GetIwant()), "graft", len(c.GetGraft()),
		"prune", len(c.GetPrune()), "idontwant", len(c.GetIdontwant()), "items", items,
	}
}

func (d *gossipDebug) rpc(kind gossipDebugKind, p peer.ID, rpc *pubsub.RPC) {
	// Mixed control/publish RPCs use the control budget so steady block traffic
	// cannot consume the budget reserved for subscription and mesh evidence.
	c := rpc.GetControl()
	hasControl := len(rpc.Subscriptions)+len(c.GetIhave())+len(c.GetIwant())+len(c.GetGraft())+len(c.GetPrune())+len(c.GetIdontwant()) > 0
	d.record(kind, len(rpc.Publish) > 0 && !hasControl,
		func() []any { return rpcFields(p, rpc) })
}

func (d *gossipDebug) inspectRPC(p peer.ID, rpc *pubsub.RPC) error {
	d.counts[debugPublishRecv].Add(uint64(len(rpc.Publish)))
	if len(rpc.Publish) > 0 {
		d.lastPublish.Store(time.Now().UnixNano())
	}
	d.rpc(debugRPCRecv, p, rpc)
	// An inspector error would reject the entire RPC. Diagnostics always allow it.
	return nil
}

// RecvRPC is intentionally empty: inspectRPC observes inbound RPCs earlier,
// including ones later dropped by the subscription filter or router score gate.
func (d *gossipDebug) RecvRPC(*pubsub.RPC) {}
func (d *gossipDebug) SendRPC(rpc *pubsub.RPC, p peer.ID) {
	d.rpc(debugRPCSend, p, rpc)
}
func (d *gossipDebug) DropRPC(rpc *pubsub.RPC, p peer.ID) {
	d.rpc(debugRPCDrop, p, rpc)
}

func (d *gossipDebug) peerEvent(kind gossipDebugKind, p peer.ID, topic string) {
	d.record(kind, false, func() []any { return []any{"peer", p, "topic", debugText(topic)} })
}

func (d *gossipDebug) AddPeer(p peer.ID, proto protocol.ID) {
	d.record(debugAddPeer, false, func() []any { return []any{"peer", p, "protocol", debugText(string(proto))} })
}
func (d *gossipDebug) RemovePeer(p peer.ID)              { d.peerEvent(debugRemovePeer, p, "") }
func (d *gossipDebug) Join(topic string)                 { d.peerEvent(debugJoin, "", topic) }
func (d *gossipDebug) Leave(topic string)                { d.peerEvent(debugLeave, "", topic) }
func (d *gossipDebug) Graft(p peer.ID, t string)         { d.peerEvent(debugGraft, p, t) }
func (d *gossipDebug) Prune(p peer.ID, t string)         { d.peerEvent(debugPrune, p, t) }
func (d *gossipDebug) ThrottlePeer(p peer.ID)            { d.peerEvent(debugThrottle, p, "") }
func (d *gossipDebug) ValidateMessage(m *pubsub.Message) { d.message(debugValidate, m, "") }
func (d *gossipDebug) RejectMessage(m *pubsub.Message, reason string) {
	d.message(debugReject, m, reason)
}
func (d *gossipDebug) DuplicateMessage(m *pubsub.Message) { d.message(debugDuplicate, m, "") }
func (d *gossipDebug) UndeliverableMessage(m *pubsub.Message) {
	d.message(debugUndeliverable, m, "subscriber_queue_full")
}
func (d *gossipDebug) DeliverMessage(m *pubsub.Message) {
	d.lastDeliver.Store(time.Now().UnixNano())
	d.message(debugDeliver, m, "")
}

func (d *gossipDebug) message(kind gossipDebugKind, msg *pubsub.Message, reason string) {
	d.record(kind, kind != debugReject && kind != debugUndeliverable, func() []any {
		return []any{"peer", msg.ReceivedFrom, "topic", debugText(msg.GetTopic()),
			"message_id", debugID(msg.ID), "bytes", len(msg.GetData()), "reason", debugText(reason)}
	})
}

func payloadFields(envelope *eth.ExecutionPayloadEnvelope) []any {
	if envelope == nil || envelope.ExecutionPayload == nil {
		return nil
	}
	p := envelope.ExecutionPayload
	return []any{"block", uint64(p.BlockNumber), "hash", p.BlockHash, "timestamp", uint64(p.Timestamp)}
}

func (d *gossipDebug) validation(p peer.ID, msg *pubsub.Message, result pubsub.ValidationResult, reason string, envelope *eth.ExecutionPayloadEnvelope) {
	kind := debugValidationReject
	if result == pubsub.ValidationAccept {
		kind = debugValidationAccept
	} else if result == pubsub.ValidationIgnore {
		kind = debugValidationIgnore
	}
	d.record(kind, result == pubsub.ValidationAccept, func() []any {
		return append([]any{"peer", p, "topic", debugText(msg.GetTopic()), "message_id", debugID(msg.ID),
			"result", validationResultString(result), "reason", reason}, payloadFields(envelope)...)
	})
}

type debugGossipIn struct {
	next  GossipIn
	debug *gossipDebug
}

func (in *debugGossipIn) OnUnsafeL2Payload(ctx context.Context, from peer.ID, envelope *eth.ExecutionPayloadEnvelope) error {
	fields := func() []any { return append([]any{"peer", from}, payloadFields(envelope)...) }
	in.debug.record(debugHandlerStart, true, fields)
	err := in.next.OnUnsafeL2Payload(ctx, from, envelope)
	if err != nil {
		in.debug.record(debugHandlerError, false, func() []any { return append(fields(), "err", debugText(err.Error())) })
	} else {
		in.debug.lastHandled.Store(time.Now().UnixNano())
		in.debug.record(debugHandlerDone, true, fields)
	}
	return err
}

func (d *gossipDebug) connection(kind gossipDebugKind, conn network.Conn) {
	d.record(kind, false, func() []any {
		return []any{"peer", conn.RemotePeer(), "connection", debugText(conn.ID()),
			"remote", debugText(conn.RemoteMultiaddr().String()), "direction", conn.Stat().Direction.String()}
	})
}

func (d *gossipDebug) start(parent context.Context, h host.Host, score func(peer.ID) (float64, error)) {
	ctx, cancel := context.WithCancel(parent)
	d.cancel = cancel
	notifier := &network.NotifyBundle{
		ConnectedF:    func(_ network.Network, conn network.Conn) { d.connection(debugConnected, conn) },
		DisconnectedF: func(_ network.Network, conn network.Conn) { d.connection(debugDisconnected, conn) },
	}
	h.Network().Notify(notifier)
	go func() {
		defer close(d.done)
		defer h.Network().StopNotify(notifier)
		ticker := time.NewTicker(gossipDebugInterval)
		defer ticker.Stop()
		d.log.Info("Gossip diagnostics enabled", "interval", gossipDebugInterval, "queue_capacity", cap(d.queue))
		for {
			select {
			case <-ctx.Done():
				return
			case rec := <-d.queue:
				d.log.Info("Gossip diagnostic event", append([]any{"event", gossipDebugNames[rec.kind], "observed_at", rec.at.UTC()}, rec.args...)...)
			case <-ticker.C:
				d.snapshot(h, score)
			}
		}
	}()
}

func (d *gossipDebug) stop() { d.cancel() }

func debugAge(now time.Time, unixNano int64) float64 {
	if unixNano == 0 {
		return -1 // never observed; zero would incorrectly imply a recent message
	}
	return now.Sub(time.Unix(0, unixNano)).Seconds()
}

// Snapshots use public host/stream APIs. They cannot observe short-lived streams
// between samples, Noise/Yamux control frames, or the remote peer's mesh state.
func (d *gossipDebug) snapshot(h host.Host, score func(peer.ID) (float64, error)) {
	now := time.Now()
	counts := make(map[string]uint64, debugKindCount)
	for kind, name := range gossipDebugNames {
		counts[name] = d.counts[kind].Load()
	}
	conns := h.Network().Conns()
	d.log.Info("Gossip diagnostic summary", "counts", counts,
		"last_publish_age_seconds", debugAge(now, d.lastPublish.Load()),
		"last_deliver_age_seconds", debugAge(now, d.lastDeliver.Load()),
		"last_handled_age_seconds", debugAge(now, d.lastHandled.Load()),
		"detail_rate_limited", d.rateLimited.Load(), "detail_queue_dropped", d.queueDropped.Load(),
		"gossip_stream_read_bytes", d.streamRead.Load(), "gossip_stream_write_bytes", d.streamWrite.Load(),
		"connections", len(conns), "connection_samples", min(len(conns), gossipDebugItems))
	for _, conn := range conns[:min(len(conns), gossipDebugItems)] {
		type streamInfo struct {
			ID        string `json:"id"`
			Protocol  string `json:"protocol"`
			Direction string `json:"direction"`
		}
		streams := conn.GetStreams()
		samples := make([]streamInfo, 0, min(len(streams), gossipDebugItems))
		for _, s := range streams[:min(len(streams), gossipDebugItems)] {
			samples = append(samples, streamInfo{debugText(s.ID()), debugText(string(s.Protocol())), s.Stat().Direction.String()})
		}
		p := conn.RemotePeer()
		value, err := score(p)
		d.log.Info("Gossip diagnostic connection", "peer", p, "connection", debugText(conn.ID()),
			"remote", debugText(conn.RemoteMultiaddr().String()), "direction", conn.Stat().Direction.String(),
			"security", conn.ConnState().Security, "muxer", conn.ConnState().StreamMultiplexer,
			"latency", h.Peerstore().LatencyEWMA(p), "score", value, "score_known", err == nil,
			"stream_count", len(streams), "stream_samples", samples)
	}
}
