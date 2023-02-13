package cmd

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

var (
	message string
)

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.PersistentFlags().StringVarP(&message, "message", "m", "DON'T PANIC!", "pubsub message to send")
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish message",
	Long: `publish the specified message to the pubsub topic (automatically creates topic if missing), eg:
	pubsubctl publish [--project=<projectID>] --topic=<topicID> [--message=<msg>]`,
	Run: func(cmd *cobra.Command, args []string) {
		topic, err := createTopic(ctx, client, topicID)
		if err != nil {
			log.Fatalf("cannot publish to topic %q: %v", topicID, err)
		}
		defer topic.Stop()

		if _, err := publish(ctx, client, topic, message); err != nil {
			log.Fatalf("cannot publish message: %v", err)
		}
		log.Println("message published")
	},
}

// publish publishes message to topic and waits for delivery, returning server-generated message id or any error occurred.
func publish(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, message string) (string, error) {
	if message == "" {
		return "", fmt.Errorf("nil message")
	}

	if topic == nil || topic.ID() == "" {
		return "", fmt.Errorf("nil topic")
	}

	res := topic.Publish(ctx, &pubsub.Message{Data: []byte(message)})
	id, err := res.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot publish message: %v", err)
	}

	return id, nil
}
