package cmd

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

var (
	poll bool
	peek bool
)

func init() {
	rootCmd.AddCommand(receiveCmd)
	receiveCmd.PersistentFlags().BoolVar(&poll, "poll", false, "poll pubsub for new messages")
	receiveCmd.PersistentFlags().BoolVar(&peek, "peek", false, "receive and nack pubsub message")
}

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "receive message",
	Long: `receive a message from the specified exisiting pubsub subscription,
	optionally, it can poll (ie, keep listening for new messages) and/or peek (ie, nacks the message after receiving and returns), eg:
	pubsubctl receive [--project=<projectID>] --subscription=<subID> [--poll] [--peek]`,
	Run: func(cmd *cobra.Command, args []string) {
		if poll && peek {
			log.Println("warning: using \"peek\" with \"poll\" can cause infinite loop!")
		}
		msg, err := receive(ctx, client, client.Subscription(subID), poll, peek)
		if err != nil {
			log.Fatalf("cannot receive message: %v", err)
		}
		if msg != nil {
			log.Printf("message received:\n%q", msg.Data)
		}
	},
}

// receive receives message(s) from subscription and returns any error occurred.
// If poll, keeps listening for new messages and prints them, otherwise, it returns first message received.
// If peek, nacks the message after receiving and returns.
// note: if there's only one consumer, combining --poll and --peek creates an infinite loop.
func receive(ctx context.Context, client *pubsub.Client, sub *pubsub.Subscription, poll bool, peek bool) (*pubsub.Message, error) {
	if sub == nil || sub.ID() == "" {
		return nil, fmt.Errorf("nil subscription")
	}
	ok, e := sub.Exists(ctx)
	if e != nil {
		return nil, fmt.Errorf("querying subscriptions: %v", e)
	}
	if !ok {
		return nil, fmt.Errorf("subscription %q does not exist", sub.String())
	}

	cctx, ccancel := context.WithCancel(ctx)
	defer ccancel()

	messages := make(chan *pubsub.Message)
	failures := make(chan error)
	var err error
	go func() {
		err = sub.Receive(cctx, func(_ context.Context, m *pubsub.Message) {
			if peek {
				m.Nack()
			} else {
				m.Ack()
			}
			messages <- m
			if !poll {
				ccancel()
			}
		})
		failures <- fmt.Errorf("receive: %v; context: %v", err, cctx.Err())
		close(failures)
		close(messages)
	}()

	if poll {
		for m := range messages {
			if m != nil {
				log.Printf("message received:\n%q", m.Data)
			}
		}
		return nil, err
	}

	select {
	case m := <-messages:
		return m, nil
	case err := <-failures:
		return nil, err
	}
}
