package environments

import (
	_ "embed"
	"fmt"
)

//go:embed anvil.json
var envAnvil []byte

//go:embed mainnet.json
var envMainnet []byte

//go:embed testnet.json
var envTestnet []byte

//go:embed testnet-staging.json
var envTestnetStaging []byte

//go:embed testnet-dev.json
var envTestnetDev []byte

type SmartContractEnvironment string

const (
	Anvil          SmartContractEnvironment = "anvil"
	Mainnet        SmartContractEnvironment = "mainnet"
	Testnet        SmartContractEnvironment = "testnet"
	TestnetStaging SmartContractEnvironment = "testnet-staging"
	TestnetDev     SmartContractEnvironment = "testnet-dev"
)

func (s *SmartContractEnvironment) UnmarshalFlag(value string) error {
	switch value {
	case string(Anvil),
		string(Mainnet),
		string(Testnet),
		string(TestnetStaging),
		string(TestnetDev):
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
	case Mainnet:
		return envMainnet, nil
	case Testnet:
		return envTestnet, nil
	case TestnetStaging:
		return envTestnetStaging, nil
	case TestnetDev:
		return envTestnetDev, nil
	default:
		return nil, fmt.Errorf("unknown environment: %s", env)
	}
}
