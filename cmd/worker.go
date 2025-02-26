package cmd

import (
	"brigade/pkg/config"
	"brigade/pkg/worker"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	plugin  string
	channel string
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start a worker that listens for messages on a channel and executes a plugin.",
	Long:  `The worker command starts a worker that listens for messages on a channel and executes a plugin.`,
	Run: func(cmd *cobra.Command, args []string) {
		var cfg *config.Config
		if cfgFile != "" {
			var err error
			cfg, err = config.Load(cfgFile)
			if err != nil {
				fmt.Println("Error loading config file:", err)
				return
			}
		}
		ctx := context.Background()
		workerClient := worker.NewWorker(cfg)
		workerClient.Start(ctx)

		<-ctx.Done()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	workerCmd.Flags().StringVarP(&plugin, "plugin", "p", "docker", "The name of the plugin to use")
	workerCmd.Flags().StringVarP(&channel, "channel", "c", "channel", "The name of the channel to use")
}
