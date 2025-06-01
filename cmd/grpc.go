package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/arvinpaundra/cent/user/api/interceptor"
	grpcapp "github.com/arvinpaundra/cent/user/application/grpc"
	"github.com/arvinpaundra/cent/user/config"
	"github.com/arvinpaundra/cent/user/core"
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/arvinpaundra/cent/user/database/nosqlpkg"
	"github.com/arvinpaundra/cent/user/database/sqlpkg"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var grpcPort string

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadEnv(".", ".env", "env")

		pgsql := sqlpkg.NewPostgres()

		sqlpkg.NewConnection(pgsql)

		rdb := nosqlpkg.NewRedisDB()

		nosqlpkg.NewInMemoryConection(rdb)

		srv := grpc.NewServer(
			grpc.UnaryInterceptor(interceptor.CheckApiKey()),
		)

		grpcapp.Register(
			srv,
			sqlpkg.GetConnection(),
			nosqlpkg.GetInMemoryConnection(),
			validator.NewValidator(),
		)

		go func() {
			addr := fmt.Sprintf(":%s", grpcPort)

			listener, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("failed to listen: %s", err.Error())
			}

			err = srv.Serve(listener)
			if err != nil {
				log.Fatalf("failed to start grpc server: %s", err.Error())
			}
		}()

		wait := core.GracefulShutdown(context.Background(), 30*time.Second, map[string]func(ctx context.Context) error{
			"grpc": func(ctx context.Context) error {
				srv.GracefulStop()

				return nil
			},
			"postgres": func(ctx context.Context) error {
				return pgsql.Close()
			},
			"redis": func(ctx context.Context) error {
				return rdb.Close()
			},
		})

		_ = <-wait
	},
}

func init() {
	grpcCmd.Flags().StringVarP(&grpcPort, "port", "p", "8083", "bind server to port")
	rootCmd.AddCommand(grpcCmd)
}
