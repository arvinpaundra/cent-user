package cmd

import (
	"context"
	"time"

	"github.com/arvinpaundra/cent/user/core"
	"github.com/spf13/cobra"
)

var pollerCmd = &cobra.Command{
	Use:   "poller",
	Short: "Start poller outbox",
	Run: func(cmd *cobra.Command, args []string) {
		wait := core.GracefulShutdown(context.Background(), 30*time.Second, map[string]func(ctx context.Context) error{})

		_ = <-wait
	},
}

func init() {
	rootCmd.AddCommand(pollerCmd)
}
