package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list topics and/or subscriptions",
	Long: `list topics and/or subscriptions under the specified pubsub project, eg:
	pubsubctl list [--project=<projectID>] [topics | subscriptions]`,
	Run: func(cmd *cobra.Command, args []string) {
		obj := ""
		if len(args) > 0 {
			obj = args[0]
		}
		if obj != "" {
			out, err := list(ctx, client, obj)
			if err != nil {
				log.Fatalf("cannot list %q: %v", obj, err)
			}
			log.Printf("%s under project %q: %s", obj, projectID, out)
			return
		}

		out, err := list(ctx, client, "topics")
		if err != nil {
			log.Fatalf("cannot list topics: %v", err)
		}
		log.Printf("%q project's topics: %s", projectID, out)

		out, err = list(ctx, client, "subscriptions")
		if err != nil {
			log.Fatalf("cannot list subscriptions: %v", err)
		}
		log.Printf("%q project's subscriptions: %s", projectID, out)
	},
}

func list(ctx context.Context, client *pubsub.Client, obj string) (string, error) {
	switch obj {
	case "topics":
		t, err := topics(ctx, client)
		if err != nil {
			return "", err
		}
		b, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "subscriptions":
		s, err := subscriptions(ctx, client)
		if err != nil {
			return "", err
		}
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("invalid list option %q, valid ones are: \"topics\" or \"subscriptions\"", obj)
	}
}

func topics(ctx context.Context, client *pubsub.Client) (map[string][]string, error) {
	it := client.Topics(ctx)
	list := map[string][]string{}
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot iterate pubsub topics: %v", err)
		}
		subs, err := subIter(t.Subscriptions(ctx))
		if err != nil {
			return nil, fmt.Errorf("cannot iterate pubsub subscriptions of topic %q: %v", t.String(), err)
		}

		list[t.String()] = subs
	}
	return list, nil
}

func subscriptions(ctx context.Context, client *pubsub.Client) ([]string, error) {
	it := client.Subscriptions(ctx)
	return subIter(it)
}

func subIter(it *pubsub.SubscriptionIterator) ([]string, error) {
	list := []string{}
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot iterate pubsub subscriptions: %v", err)
		}
		list = append(list, s.String())
	}
	return list, nil
}
