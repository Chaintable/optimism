package p2p

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-service/testlog"
)

type debugTestConn struct{ network.Conn }

func (debugTestConn) ID() string          { return "connection" }
func (debugTestConn) RemotePeer() peer.ID { return "remote" }

type debugTestStream struct {
	network.Stream
	err      error
	written  []byte
	closes   int
	resets   int
	deadline time.Time
}

func (*debugTestStream) ID() string            { return "stream" }
func (*debugTestStream) Conn() network.Conn    { return debugTestConn{} }
func (*debugTestStream) Protocol() protocol.ID { return pubsub.GossipSubID_v12 }
func (*debugTestStream) Stat() network.Stats   { return network.Stats{Direction: network.DirOutbound} }
func (*debugTestStream) Read(b []byte) (int, error) {
	return copy(b, "abc"), io.EOF
}
func (s *debugTestStream) Write(b []byte) (int, error) {
	s.written = append(s.written, b...)
	return 2, s.err
}
func (s *debugTestStream) Close() error { s.closes++; return s.err }
func (s *debugTestStream) Reset() error { s.resets++; return s.err }
func (s *debugTestStream) SetReadDeadline(t time.Time) error {
	s.deadline = t
	return s.err
}

func TestGossipDebugStreamPassthrough(t *testing.T) {
	d := newGossipDebug(testlog.Logger(t, log.LevelInfo), "self")
	original := &debugTestStream{err: errors.New("write interrupted")}
	s := (&gossipDebugHost{debug: d}).wrap(original)
	b := make([]byte, 8)
	n, err := s.Read(b)
	require.Equal(t, 3, n)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, "abc", string(b[:n]))
	n, err = s.Write([]byte("hello"))
	require.Equal(t, 2, n)
	require.Same(t, original.err, err)
	require.Equal(t, []byte("hello"), original.written)
	require.Same(t, original.err, s.Close())
	require.Same(t, original.err, s.Reset())
	require.Equal(t, 1, original.closes)
	require.Equal(t, 1, original.resets)
	deadline := time.Now()
	require.Same(t, original.err, s.SetReadDeadline(deadline))
	require.Equal(t, deadline, original.deadline)
	require.EqualValues(t, 3, d.streamRead.Load())
	require.EqualValues(t, 2, d.streamWrite.Load())
	require.EqualValues(t, 1, d.counts[debugStreamReadError].Load())
	require.EqualValues(t, 1, d.counts[debugStreamWriteError].Load())
	<-d.queue // open
	require.Equal(t, "EOF", debugFields((<-d.queue).args)["err"])
}

// Exercise the actual pubsub hooks and host adapter together, using an in-memory
// network. No TCP listeners or external peers are involved.
func TestGossipDebugPubSub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	mn := mocknet.New()
	t.Cleanup(func() { require.NoError(t, mn.Close()) })
	h1, err := mn.GenPeer()
	require.NoError(t, err)
	h2, err := mn.GenPeer()
	require.NoError(t, err)
	d := newGossipDebug(testlog.Logger(t, log.LevelInfo), h1.ID())
	h := &gossipDebugHost{Host: h1, debug: d}
	ps1, err := pubsub.NewGossipSub(ctx, h, d.options()...)
	require.NoError(t, err)
	ps2, err := pubsub.NewGossipSub(ctx, h2, pubsub.WithFloodPublish(true))
	require.NoError(t, err)
	t1, err := ps1.Join("blocks")
	require.NoError(t, err)
	t2, err := ps2.Join("blocks")
	require.NoError(t, err)
	s1, err := t1.Subscribe()
	require.NoError(t, err)
	s2, err := t2.Subscribe()
	require.NoError(t, err)
	t.Cleanup(func() {
		s1.Cancel()
		s2.Cancel()
		require.NoError(t, t1.Close())
		require.NoError(t, t2.Close())
	})
	require.NoError(t, mn.LinkAll())
	require.NoError(t, mn.ConnectAllButSelf())
	require.Eventually(t, func() bool { return len(t1.ListPeers()) == 1 && len(t2.ListPeers()) == 1 }, 5*time.Second, 10*time.Millisecond)
	require.NoError(t, t2.Publish(ctx, []byte("block payload")))
	msg, err := s1.Next(ctx)
	require.NoError(t, err)
	require.Equal(t, []byte("block payload"), msg.Data)
	require.Positive(t, d.counts[debugRPCRecv].Load())
	require.Positive(t, d.counts[debugPublishRecv].Load())
	require.Positive(t, d.counts[debugDeliver].Load())
	require.Positive(t, d.counts[debugStreamOpen].Load())
	require.Positive(t, d.streamRead.Load())
	require.Positive(t, d.streamWrite.Load())
}
