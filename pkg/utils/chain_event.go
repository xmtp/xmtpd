package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetEventSig extracts the event signature from an ABI by event name.
func GetEventSig(abi *abi.ABI, eventName string) (string, error) {
	event, ok := abi.Events[eventName]
	if !ok {
		return "", fmt.Errorf("event %s not found", eventName)
	}
	return event.Sig, nil
}

// GetEventTopic extracts the event topic (the hash of the signature) from an ABI by event name.
func GetEventTopic(abi *abi.ABI, eventName string) (common.Hash, error) {
	sig, err := GetEventSig(abi, eventName)
	if err != nil {
		return common.Hash{}, err
	}
	hashed := crypto.Keccak256Hash([]byte(sig))
	return hashed, nil
}
