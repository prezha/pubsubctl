package cmd

import (
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/cobra"
)

var (
	testTopic        = "pubsubctl-test-topic"
	testSubscription = "pubsubctl-test-subscription"
	testMessage      = "pubsubctl test message"
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringVarP(&testMessage, "message", "m", testMessage, "pubsub message to send")
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test pubsub emulator",
	Long: `publish a test message to a test pubsub topic and receive it via test subscription (automatically creates project/topic/subscription if missing), eg:
	pubsubctl test [--project=<projectID>] [--topic=<topicID>] [--subscription=<subID>] [--message=<msg>]`,
	Run: func(cmd *cobra.Command, args []string) {
		if topicID == "" {
			topicID = testTopic
		}
		if subID == "" {
			subID = testSubscription
		}
		var m *pubsub.Message
		defer func() {
			if err := testCleanup(topicID, subID, m, peek); err != nil {
				log.Printf("cannot cleanup after test: %v", err)
			} else {
				log.Println("test cleanup successfull")
			}
			log.Println("test ended.")
		}()

		log.Println("test started...")

		topic, err := createTopic(ctx, client, topicID)
		if err != nil {
			log.Fatalf("cannot create topic %q: %v", topicID, err)
		}
		defer topic.Stop()
		log.Printf("created topic %q", client.Topic(topicID).String())

		sub, err := createSubscription(ctx, client, topic, subID)
		if err != nil {
			log.Fatalf("cannot create subscription %q: %v", subID, err)
		}
		log.Printf("created subscription %q", client.Subscription(subID).String())

		if _, err := publish(ctx, client, topic, testMessage); err != nil {
			log.Fatalf("cannot publish message: %v", err)
		}
		log.Printf("message:\n%q\npublished to topic %q", testMessage, client.Topic(topicID).String())

		m, err = receive(ctx, client, sub, poll, peek)
		if err != nil {
			log.Fatalf("test failed: %v", err)
		}
		log.Printf("message received from subscription %q", client.Subscription(subID).String())
		if !poll {
			if m == nil {
				log.Printf("test failed: received 'nil' message")
				return
			}
			if string(m.Data) != testMessage {
				log.Printf("test failed: expected: %q, received: %q", testMessage, m.Data)
			}
		}
		msg := "!"
		if m != nil {
			msg = fmt.Sprintf("%q", m.Data)
		}
		log.Printf("test passed: received expected message:\n%s", msg)
	},
}

// testCleanup removes pubsubctlTestTopic and pubsubctlTestSubscription and also acks msg.
func testCleanup(topicID string, subID string, msg *pubsub.Message, peek bool) error {
	if !peek && msg != nil {
		msg.Ack()
		log.Println("message acked")
	}

	if subID == testSubscription {
		if err := delete(ctx, client, client.Subscription(subID).String()); err != nil {
			return err
		}
		log.Printf("deleted subscription %q", client.Subscription(subID).String())
	}

	if topicID == testTopic {
		if err := delete(ctx, client, client.Topic(topicID).String()); err != nil {
			return err
		}
		log.Printf("deleted topic %q", client.Topic(topicID).String())
	}

	return nil
}
