package cmd

import (
	"net"

	"github.com/chrisseto/white-elephant/pkg/migrations"
	"github.com/chrisseto/white-elephant/pkg/server"
	"github.com/cockroachdb/errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	listenAddress string
	dbURI         string
)

func init() {
	serverCmd.Flags().StringVarP(&listenAddress, "listen", "l", ":6969", "the address to listen on")
	serverCmd.Flags().StringVarP(&dbURI, "database", "d", "postgresql://root@localhost:4445/defaultdb", "the database URI")
}

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "TODO",
	Long:  `Also TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := RootContext()

		db, err := sqlx.Open("pgx", dbURI)
		if err != nil {
			return errors.Wrap(err, "connecting to DB")
		}

		defer db.Close()

		if err := migrations.Up(ctx, db.DB); err != nil {
			return err
		}

		lis, err := net.Listen("tcp", listenAddress)
		if err != nil {
			return errors.Wrapf(err, "listening on %s", listenAddress)
		}

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			s := server.New(db)
			return s.Serve(ctx, lis)
		})

		zap.L().Info("listening", zap.String("address", listenAddress))

		return group.Wait()
	},
}
