package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "white-elephant",
	Short: "TODO",
	Long:  `Also TODO`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		l, err := zap.NewDevelopment()
		if err != nil {
			return errors.Wrap(err, "creating zap logger")
		}

		zap.ReplaceGlobals(l)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func RootContext() context.Context {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		zap.L().Info("received SIGTERM")
		cancel()
	}()

	return ctx
}
