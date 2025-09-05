package commands

import (
	"fmt"
	"log"
	"os"

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
	globalConfigFile  string
	globalLogLevel    string
	globalLogEncoding string
	globalPrivateKey  string
	globalRpcURL      string
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
		StringVarP(&globalRpcURL, "rpc-url", "r", "", "RPC URL to use")

	if err := viper.BindPFlag("rpc-url", rootCmd.PersistentFlags().Lookup("rpc-url")); err != nil {
		return err
	}

	if err := viper.BindEnv("rpc-url", "RPC_URL"); err != nil {
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

	return nil
}

func cliLogger() (*zap.Logger, error) {
	l, _, err := utils.BuildLogger(config.LogOptions{
		LogLevel:    viper.GetString("log-level"),
		LogEncoding: viper.GetString("log-encoding"),
	})
	if err != nil || l == nil {
		return nil, fmt.Errorf("could not build logger: %s", err)
	}

	return l, nil
}

func init() {
	// Add a hidden command to generate completion scripts
	rootCmd.AddCommand(&cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion script",
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
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
