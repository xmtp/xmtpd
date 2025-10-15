package environments

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"
)

//go:embed anvil.json
var envAnvil []byte

//go:embed testnet.json
var envTestnet []byte

//go:embed testnet-staging.json
var envTestnetStaging []byte

type SmartContractEnvironment string

const (
	Anvil          SmartContractEnvironment = "anvil"
	Testnet        SmartContractEnvironment = "testnet"
	TestnetStaging SmartContractEnvironment = "testnet-staging"
)

func (s *SmartContractEnvironment) UnmarshalFlag(value string) error {
	validChoices := []string{string(Anvil), string(Testnet), string(TestnetStaging)}
	if !slices.Contains(validChoices, value) {
		joined := strings.Join(validChoices, ", ")
		return fmt.Errorf("invalid environment: %s (valid choices: %s)", value, joined)
	}

	*s = SmartContractEnvironment(value)
	return nil
}

func GetEnvironmentConfig(env SmartContractEnvironment) ([]byte, error) {
	switch env {
	case Anvil:
		return envAnvil, nil
	case Testnet:
		return envTestnet, nil
	case TestnetStaging:
		return envTestnetStaging, nil
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
}
