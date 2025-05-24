package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arvinpaundra/cent/user/api/route"
	"github.com/arvinpaundra/cent/user/config"
	"github.com/arvinpaundra/cent/user/core"
	"github.com/arvinpaundra/cent/user/core/validator"
	sqlpkg "github.com/arvinpaundra/cent/user/database"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var port string

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Start rest server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadEnv(".", ".env", "env")

		sqlpkg.NewConnection(sqlpkg.NewPostgres())

		g := gin.New()

		route.New(g, sqlpkg.GetConnection(), validator.NewValidator())

		srv := http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: g,
		}

		go func() {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("failed to start server: %s", err.Error())
			}
		}()

		wait := core.GracefulShutdown(context.Background(), 30*time.Second, map[string]func(ctx context.Context) error{
			"rest-server": func(_ context.Context) error {
				return srv.Close()
			},
			"postgres": func(_ context.Context) error {
				db, err := sqlpkg.GetConnection().DB()
				if err != nil {
					return err
				}

				return db.Close()
			},
		})

		_ = <-wait
	},
}

func init() {
	restCmd.Flags().StringVarP(&port, "port", "p", "8080", "bind server to port. default: 8080")
	rootCmd.AddCommand(restCmd)
}
