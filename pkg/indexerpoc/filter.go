package indexerpoc

import "github.com/ethereum/go-ethereum/common"

// Filter contains filtering options for blockchain data retrieval
type Filter struct {
	Addresses []common.Address
	Topics    [][]common.Hash
}

// NewFilter creates a filter for specific addresses and topics
func NewFilter(addressStrings []string, topicStrings [][]string) *Filter {
	addresses := make([]common.Address, 0, len(addressStrings))
	for _, addr := range addressStrings {
		addresses = append(addresses, common.HexToAddress(addr))
	}

	topics := make([][]common.Hash, 0, len(topicStrings))
	for _, topicList := range topicStrings {
		topicHashes := make([]common.Hash, 0, len(topicList))
		for _, topic := range topicList {
			topicHashes = append(topicHashes, common.HexToHash(topic))
		}
		topics = append(topics, topicHashes)
	}

	return &Filter{
		Addresses: addresses,
		Topics:    topics,
	}
}
