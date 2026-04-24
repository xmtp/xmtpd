package environments

import (
	_ "embed"
	"fmt"
)

//go:embed anvil.json
var envAnvil []byte

//go:embed testnet.json
var envTestnet []byte

type SmartContractEnvironment string

const (
	Anvil   SmartContractEnvironment = "anvil"
	Testnet SmartContractEnvironment = "testnet"
)

func (s *SmartContractEnvironment) UnmarshalFlag(value string) error {
	switch value {
	case string(Anvil), string(Testnet):
		*s = SmartContractEnvironment(value)
		return nil
	default:
		// do not advertise staging in the options, keep it as a hidden option
		return fmt.Errorf("unknown environment type: %s (valid choices: testnet, mainnet)", value)
	}
}

func GetEnvironmentConfig(env SmartContractEnvironment) ([]byte, error) {
	switch env {
	case Anvil:
		return envAnvil, nil
	case Testnet:
		return envTestnet, nil
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
}
