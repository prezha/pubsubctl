package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete topic and/or subscription",
	Long: `delete specified object with the given path path, eg:
	pubsubctl delete [--project=<projectID>] --topic=<topicID>
	or
	pubsubctl delete [--project=<projectID>] --topic=<topicID> --subscription=<subID>
	or
	pubsubctl delete [--project=<projectID>] --subscription=<subID>
	or
	pubsubctl delete <path>`,
	Run: func(cmd *cobra.Command, args []string) {
		if topicID == "" && subID == "" && len(args) == 0 {
			log.Fatalf("topic, subscription or path must be specified")
		}

		if topicID != "" {
			path := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
			if err := delete(ctx, client, path); err != nil {
				log.Fatalf("cannot delete topic %q: %v", path, err)
			}
			log.Printf("topic %q deleted", path)
		}

		if subID != "" {
			path := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID)
			if err := delete(ctx, client, path); err != nil {
				log.Fatalf("cannot delete subscription %q: %v", path, err)
			}
			log.Printf("subscription %q deleted", path)
		}

		if len(args) > 0 {
			path := args[0]
			if err := delete(ctx, client, path); err != nil {
				log.Fatalf("cannot delete path %q: %v", path, err)
			}
			log.Printf("path %q deleted", path)
		}
	},
}

// delete deletes object with given path.
// ref: https://cloud.google.com/pubsub/docs/admin#delete_a_topic
// ref: https://cloud.google.com/pubsub/docs/create-subscription#delete_subscription
// ref: https://cloud.google.com/pubsub/docs/create-subscription#detach_a_subscription_from_a_topic
func delete(ctx context.Context, client *pubsub.Client, path string) error {
	p := strings.Split(path, "/")
	if len(p) != 4 {
		return fmt.Errorf("invalid path %q", path)
	}

	obj := p[2]
	id := p[3]
	switch obj {
	case "topics":
		t := client.Topic(id)
		defer t.Stop()
		err := t.Delete(ctx)
		return err
	case "subscriptions":
		s := client.Subscription(id)
		return s.Delete(ctx)
	default:
		return fmt.Errorf("invalid path %q", path)
	}
}
