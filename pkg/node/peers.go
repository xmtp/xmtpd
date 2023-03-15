package node

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/xmtp/xmtpd/pkg/zap"
	"go.uber.org/zap/zapcore"
)

// persistentPeers maintains a connection to the given set of peer addresses.
// The libp2p package does not attempt to reconnect to disconnected peers or
// maintain a minimum number of connections to peers, so the pattern of a
// "persistent peer" is useful in the case where there is no DHT peer discovery
// mechanism being used that would otherwise maintain minimum connectivity, or
// in the case where you would like to always connect to a specific peer.
type persistentPeers struct {
	host host.Host
}

func newPersistentPeers(ctx context.Context, log *zap.Logger, host host.Host, addrs []string) (*persistentPeers, error) {
	p := &persistentPeers{
		host: host,
	}

	peers := make([]peer.AddrInfo, 0, len(addrs))
	for _, addr := range addrs {
		maddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, errors.Wrap(err, "parsing persistent peer address")
		}
		peer, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			return nil, errors.Wrap(err, "getting persistent peer address info")
		}
		if peer == nil {
			return nil, fmt.Errorf("persistent peer address info is nil: %s", addr)
		}
		if peer.ID == host.ID() {
			continue
		}
		peers = append(peers, *peer)
	}

	if log.Level() == zapcore.DebugLevel {
		// Log connected peers periodically.
		go func() {
			for {
				peers := p.connectedPeers()
				ids := make([]peer.ID, 0, len(peers))
				for id := range peers {
					ids = append(ids, peer.ID(id))
				}
				log.Debug("connected peers", zap.Int("total_peers", len(peers)), zap.PeerIDs("peers", ids...))

				ticker := time.NewTicker(30 * time.Second)
				defer ticker.Stop()
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}()
	}

	// Connect to p2p persistent peers.
	for _, peer := range peers {
		peer := peer
		go func() {
			for {
				peers := p.connectedPeers()
				if _, ok := peers[peer.ID.Pretty()]; !ok {
					err := backoff.Retry(func() error {
						ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
						defer cancel()
						log.Debug("connecting to persistent peer", zap.PeerID("peer", peer.ID))
						err := p.host.Connect(ctx, peer)
						if err != nil {
							log.Debug("error connecting to persistent peer", zap.Error(err), zap.PeerID("peer", peer.ID))
						}
						return err
					}, backoff.NewExponentialBackOff())
					if err != nil {
						log.Error("error connecting to persistent peer", zap.Error(err))
					} else {
						log.Info("connected to persistent peer", zap.PeerID("peer", peer.ID))
					}
				}

				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}()
	}

	return p, nil
}

func (p *persistentPeers) connectedPeers() map[string]*peer.AddrInfo {
	peers := map[string]*peer.AddrInfo{}
	for _, conn := range p.host.Network().Conns() {
		peers[conn.RemotePeer().Pretty()] = &peer.AddrInfo{
			ID:    conn.RemotePeer(),
			Addrs: []multiaddr.Multiaddr{conn.RemoteMultiaddr()},
		}
	}
	return peers
}
