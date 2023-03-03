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
)

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

	// Log connected peers periodically.
	go func() {
		for {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				peers := p.connectedPeers()
				addrs := make([]string, 0, len(peers))
				for addr := range peers {
					addrs = append(addrs, addr)
				}
				log.Info("total connected peers", zap.Int("total_peers", len(peers)), zap.Strings("peers", addrs))
			}
		}
	}()

	// Connect to p2p persistent peers.
	for _, peer := range peers {
		peer := peer
		go func() {
			for {
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					peers := p.connectedPeers()
					if _, ok := peers[peer.ID.Pretty()]; ok {
						continue
					}
					err := backoff.Retry(func() error {
						ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
						defer cancel()
						log.Info("connecting to persistent peer", zap.String("peer", peer.ID.Pretty()))
						return p.host.Connect(ctx, peer)
					}, backoff.NewExponentialBackOff())
					if err != nil {
						log.Error("error connecting to persistent peer", zap.Error(err))
					}
					log.Info("connected to persistent peer", zap.String("peer", peer.ID.Pretty()))
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
