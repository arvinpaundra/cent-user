package cmd

import (
	"context"
	"time"

	pollerapp "github.com/arvinpaundra/cent/user/application/poller"
	"github.com/arvinpaundra/cent/user/config"
	"github.com/arvinpaundra/cent/user/core"
	"github.com/arvinpaundra/cent/user/core/messaging"
	"github.com/arvinpaundra/cent/user/core/poller"
	"github.com/arvinpaundra/cent/user/database/sqlpkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pollerCmd = &cobra.Command{
	Use:   "poller",
	Short: "Start poller event outbox",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadEnv(".", ".env", "env")

		pgsql := sqlpkg.NewPostgres()

		sqlpkg.NewConnection(pgsql)

		p := poller.NewPoller().SetBaseDelay(5 * time.Second).SetMaxDelay(3 * time.Hour)

		nc := messaging.NewNats(viper.GetString("NATS_URL"))

		go pollerapp.StartWorker(p, sqlpkg.GetConnection(), nc.GetConnection())

		wait := core.GracefulShutdown(context.Background(), 30*time.Second, map[string]func(ctx context.Context) error{
			"poller": func(_ context.Context) error {
				return p.Close()
			},
			"postgres": func(_ context.Context) error {
				return pgsql.Close()
			},
		})

		_ = <-wait
	},
}

func init() {
	rootCmd.AddCommand(pollerCmd)
}
