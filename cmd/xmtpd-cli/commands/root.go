// Package commands implements the CLI commands based on cobra.
package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/xmtp/xmtpd/pkg/config/environments"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "xmtpd-cli",
	Short: "xmtpd-cli is a CLI to manage the XMTP Network",
	Long:  `xmtpd-cli is a CLI to manage the XMTP Network`,
}

var (
	environment         string
	globalConfigFile    string
	globalLogLevel      string
	globalLogEncoding   string
	globalPrivateKey    string
	globalSettlementURL string
	globalAppURL        string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := configureRootCmd()
	if err != nil {
		log.Fatalf("could not configure root command: %s", err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatalf("could not execute root command: %s", err)
	}
}

func configureRootCmd() error {
	err := registerGlobalFlags()
	if err != nil {
		return err
	}

	rootCmd.AddCommand(
		keyManagementCmd(),
		nodeRegistryCmd(),
		rateRegistryCmd(),
		appChainCmd(),
		settlementChainCmd(),
		generateCmd(),
		fundsCmd(),
		paramsCmd(),
		versionCmd(),
	)

	return nil
}

func registerGlobalFlags() error {
	rootCmd.PersistentFlags().
		StringVarP(&globalLogLevel, "log-level", "l", "info", "set logging level. Available levels: debug, info, warn, error, fatal, panic")

	if err := viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		return err
	}
	if err := viper.BindEnv("log-level", "LOG_LEVEL"); err != nil {
		return err
	}

	rootCmd.PersistentFlags().
		StringVarP(&globalLogEncoding, "log-encoding", "e", "console", "set log encoding. Available encodings: console, json")

	if err := viper.BindPFlag("log-encoding", rootCmd.PersistentFlags().Lookup("log-encoding")); err != nil {
		return err
	}
	if err := viper.BindEnv("log-encoding", "LOG_ENCODING"); err != nil {
		return err
	}

	rootCmd.PersistentFlags().
		StringVarP(&globalConfigFile, "config-file", "c", "", "path to the config file")

	if err := viper.BindPFlag("config-file", rootCmd.PersistentFlags().Lookup("config-file")); err != nil {
		return err
	}
	if err := viper.BindEnv("config-file", "XMTPD_CONTRACTS_CONFIG_FILE_PATH"); err != nil {
		return err
	}

	rootCmd.PersistentFlags().
		StringVarP(&environment, "environment", "", "", "Deployed environment to load contracts config for")

	if err := viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("environment")); err != nil {
		return err
	}
	if err := viper.BindEnv("environment", "XMTPD_CONTRACTS_ENVIRONMENT"); err != nil {
		return err
	}

	rootCmd.PersistentFlags().
		StringVarP(&globalPrivateKey, "private-key", "p", "", "private key to use")

	if err := viper.BindPFlag("private-key", rootCmd.PersistentFlags().Lookup("private-key")); err != nil {
		return err
	}
	if err := viper.BindEnv("private-key", "PRIVATE_KEY"); err != nil {
		return err
	}

	rootCmd.PersistentFlags().
		StringVar(&globalSettlementURL, "settlement-rpc-url", "", "Settlement chain RPC URL")
	if err := viper.BindPFlag("settlement-rpc-url", rootCmd.PersistentFlags().Lookup("settlement-rpc-url")); err != nil {
		return err
	}
	if err := viper.BindEnv("settlement-rpc-url", "SETTLEMENT_RPC_URL"); err != nil {
		return err
	}

	rootCmd.PersistentFlags().
		StringVar(&globalAppURL, "app-rpc-url", "", "App chain RPC URL")
	if err := viper.BindPFlag("app-rpc-url", rootCmd.PersistentFlags().Lookup("app-rpc-url")); err != nil {
		return err
	}
	if err := viper.BindEnv("app-rpc-url", "APP_RPC_URL"); err != nil {
		return err
	}

	return nil
}

func cliLogger() (*zap.Logger, error) {
	l, _, err := utils.BuildLogger(config.LogOptions{
		LogLevel:    viper.GetString("log-level"),
		LogEncoding: viper.GetString("log-encoding"),
	})
	if err != nil || l == nil {
		return nil, fmt.Errorf("could not build logger: %w", err)
	}
	return l, nil
}

func resolveSettlementRPCURL() (string, error) {
	if v := viper.GetString("settlement-rpc-url"); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("missing settlement RPC URL: set --settlement-rpc-url")
}

func resolveAppRPCURL() (string, error) {
	if v := viper.GetString("app-rpc-url"); v != "" {
		return v, nil
	}
	return "", fmt.Errorf("missing app RPC URL: set --app-rpc-url")
}

func resolveConfig(configFile string, environment string) (*config.ContractsOptions, error) {
	if environment != "" && configFile != "" {
		return nil, fmt.Errorf(
			"--environment cannot be used with --config-file",
		)
	}

	if environment != "" {
		var env environments.SmartContractEnvironment
		err := env.UnmarshalFlag(environment)
		if err != nil {
			return nil, err
		}

		data, err := environments.GetEnvironmentConfig(
			env,
		)
		if err != nil {
			return nil, err
		}

		var cfg config.ChainConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}

		var resultOptions config.ContractsOptions

		config.FillConfigFromJSON(&resultOptions, &cfg)

		return &resultOptions, nil
	}

	if configFile == "" {
		return nil, fmt.Errorf("config-file or environment is required")
	}

	contracts, err := config.ContractOptionsFromEnv(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not load config from file: %w", err)
	}

	return contracts, nil
}

func init() {
	// Add a hidden command to generate completion scripts
	rootCmd.AddCommand(&cobra.Command{
		Use:          "completion [bash|zsh|fish|powershell]",
		Short:        "Generate shell completion script",
		ValidArgs:    []string{"bash", "zsh", "fish", "powershell"},
		Args:         cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Hidden:       true,
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				_ = rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				_ = rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				_ = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	})
}
