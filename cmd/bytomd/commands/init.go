package commands

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	cfg "github.com/bytom-gm/config"
)

var initFilesCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize blockchain",
	Run:   initFiles,
}

func init() {
	initFilesCmd.Flags().String("chain_id", config.ChainID, "Select [gm-testnet] or [solonet]")

	RootCmd.AddCommand(initFilesCmd)
}

func initFiles(cmd *cobra.Command, args []string) {
	configFilePath := path.Join(config.RootDir, "config.toml")
	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		log.WithField("config", configFilePath).Info("Already exists config file.")
		return
	}

	switch config.ChainID {
	case "gm-testnet":
		cfg.EnsureRoot(config.RootDir, config.ChainID)
	default:
		cfg.EnsureRoot(config.RootDir, "solonet")
	}

	log.WithField("config", configFilePath).Info("Initialized bytom")
}
