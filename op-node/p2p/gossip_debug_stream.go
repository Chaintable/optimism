package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// Only GossipSub receives this host adapter. Reads/writes and all connection,
// deadline, close and reset operations still go to the original implementation.
// These bytes are decrypted application-stream bytes, not TCP/Noise/Yamux or ping.
type gossipDebugHost struct {
	host.Host
	debug *gossipDebug
}

func (h *gossipDebugHost) NewStream(ctx context.Context, p peer.ID, ids ...protocol.ID) (network.Stream, error) {
	s, err := h.Host.NewStream(ctx, p, ids...)
	if err != nil {
		h.debug.record(debugStreamOpenError, false, func() []any { return []any{"peer", p, "err", debugText(err.Error())} })
		return s, err
	}
	return h.wrap(s), nil
}

func (h *gossipDebugHost) wrap(s network.Stream) *gossipDebugStream {
	out := &gossipDebugStream{Stream: s, debug: h.debug}
	out.event(debugStreamOpen, nil)
	return out
}

func (h *gossipDebugHost) handler(next network.StreamHandler) network.StreamHandler {
	return func(s network.Stream) {
		wrapped := h.wrap(s)
		defer wrapped.event(debugStreamHandlerExit, nil)
		next(wrapped)
	}
}

func (h *gossipDebugHost) SetStreamHandler(id protocol.ID, next network.StreamHandler) {
	h.Host.SetStreamHandler(id, h.handler(next))
}

func (h *gossipDebugHost) SetStreamHandlerMatch(id protocol.ID, match func(protocol.ID) bool, next network.StreamHandler) {
	h.Host.SetStreamHandlerMatch(id, match, h.handler(next))
}

type gossipDebugStream struct {
	network.Stream
	debug *gossipDebug
}

func (s *gossipDebugStream) event(kind gossipDebugKind, err error) {
	s.debug.record(kind, false, func() []any {
		fields := []any{"peer", s.Conn().RemotePeer(), "connection", debugText(s.Conn().ID()),
			"stream", debugText(s.ID()), "protocol", debugText(string(s.Protocol())), "direction", s.Stat().Direction.String()}
		if err != nil {
			fields = append(fields, "err", debugText(err.Error()))
		}
		return fields
	})
}

func (s *gossipDebugStream) Read(b []byte) (int, error) {
	n, err := s.Stream.Read(b)
	s.debug.streamRead.Add(uint64(n))
	if err != nil {
		s.event(debugStreamReadError, err) // EOF also identifies a graceful remote half-close
	}
	return n, err
}

func (s *gossipDebugStream) Write(b []byte) (int, error) {
	n, err := s.Stream.Write(b)
	s.debug.streamWrite.Add(uint64(n))
	if err != nil {
		s.event(debugStreamWriteError, err)
	}
	return n, err
}

func (s *gossipDebugStream) Close() error {
	err := s.Stream.Close()
	s.event(debugStreamClose, err)
	return err
}

func (s *gossipDebugStream) Reset() error {
	err := s.Stream.Reset()
	s.event(debugStreamReset, err)
	return err
}
