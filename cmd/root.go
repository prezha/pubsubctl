package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

const version = "0.6.5-beta"

var (
	ctx    context.Context
	cancel context.CancelFunc

	client *pubsub.Client

	cloud bool
	host  string = "localhost:8085"

	projectID string = "default"
	topicID   string
	subID     string

	timeout time.Duration
)

// init sets global flags - for root command and all subcommands
func init() {
	rootCmd.PersistentFlags().BoolVar(&cloud, "cloud", false, "use cloud pubsub instead of the emulator")

	// if set, use $PUBSUB_EMULATOR_HOST as default value for host flag
	if env := os.Getenv("PUBSUB_EMULATOR_HOST"); env != "" {
		host = env
	}
	rootCmd.PersistentFlags().StringVar(&host, "host", host, "[address:port] of the emulator host, defaulting to PUBSUB_EMULATOR_HOST environment variable value (if set), ignored if 'cloud' flag is also set")

	// if set, use $PUBSUB_PROJECT_ID as default value for projectID flag
	if env := os.Getenv("PUBSUB_PROJECT_ID"); env != "" {
		projectID = env
	}
	rootCmd.PersistentFlags().StringVarP(&projectID, "project", "p", projectID, "pubsub project, defaulting to PUBSUB_PROJECT_ID environment variable value (if set)")

	rootCmd.PersistentFlags().StringVarP(&topicID, "topic", "t", "", "pubsub topic")
	rootCmd.PersistentFlags().StringVarP(&subID, "subscription", "s", "", "pubsub subscription")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 5*time.Second, "time to wait for command execution (value <=0 disables timeout)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "pubsubctl",
	Short: "pubsubctl",
	Long: fmt.Sprintf(`pubsubctl v%s
	pubsubctl is a basic Google Cloud Platform Pub/Sub [Emulator] CLI`, version),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cloud {
			if err := os.Unsetenv("PUBSUB_EMULATOR_HOST"); err != nil {
				log.Fatalf("cannot unset environment variable \"PUBSUB_EMULATOR_HOST\": %v", err)
			}
		} else {
			if err := os.Setenv("PUBSUB_EMULATOR_HOST", host); err != nil {
				log.Fatalf("cannot set environment variable \"PUBSUB_EMULATOR_HOST\": %v", err)
			}
		}

		// ignore default timeout for polling
		if poll && !cmd.Flags().Changed("timeout") {
			timeout = 0
		}

		ctx = context.Background()
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			// note: 'defer cancel()' goes to PersistentPostRun as we'll use the context in subcommands
		}

		var err error
		client, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("cannot initialise pubsub client: %v", err)
		}
		// note: 'defer client.Close()' goes to PersistentPostRun as we'll use the client in subcommands

		if _, err := subscriptions(ctx, client); err != nil {
			log.Fatalf("cannot connect to pubsub emulator (check if it's listening on %q): %v", host, err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cancel != nil {
			cancel()
		}
		if client != nil {
			client.Close()
		}
	},
}
