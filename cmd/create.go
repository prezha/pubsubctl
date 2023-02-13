package cmd

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create topic/subscription",
	Long: `create the specified topic and optionally subscription (if given) under the pubsub project/topic, eg:
	pubsubctl create [--project=<projectID>] --topic=<topicID> [--subscription=<subID>]`,
	Run: func(cmd *cobra.Command, args []string) {
		topic, err := createTopic(ctx, client, topicID)
		if err != nil {
			log.Fatalf("cannot create topic %q: %v", topicID, err)
		}
		defer topic.Stop()
		log.Printf("created %q", topic.String())

		if subID != "" {
			sub, err := createSubscription(ctx, client, topic, subID)
			if err != nil {
				log.Fatalf("cannot create subscription %q: %v", subID, err)
			}
			log.Printf("created %q", sub.String())
		}
	},
}

// createTopic returns exiting topic or creates new one.
func createTopic(ctx context.Context, client *pubsub.Client, id string) (*pubsub.Topic, error) {
	if id == "" {
		return nil, fmt.Errorf("nil topic")
	}

	topic := client.Topic(id)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		topic, err = client.CreateTopic(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return topic, nil
}

// createSubscription returns exiting subscription for topic or creates new one.
func createSubscription(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, id string) (*pubsub.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("nil subscription")
	}

	sub := client.Subscription(id)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}

	// delete existing subscription if its config cannot be fetched or it's not linked to same topic name
	if ok {
		cfg, err := sub.Config(ctx)
		if err != nil || cfg.Topic.String() != topic.String() {
			if err := delete(ctx, client, sub.String()); err != nil {
				return nil, fmt.Errorf("cannot delete existing subscription: %v", err)
			}
		}
	}

	sub, err = client.CreateSubscription(ctx, id, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		return nil, err
	}

	return sub, nil
}
