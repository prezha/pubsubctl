# pubsubctl
pubsubctl is a basic Google Pub/Sub Emulator/Cloud CLI

## quick start - basic _emulator_ usage
1. start pubsub emulator in a separate session (to monitor logs)
 - using the gCloud Docker image:
> docker run --rm -ti -p 127.0.0.1:8085:8085 gcr.io/google.com/cloudsdktool/google-cloud-cli:latest gcloud beta emulators pubsub start --host-port=0.0.0.0:8085 --project=my-project --log-http --verbosity=debug --user-output-enabled
 - using the Google Cloud CLI:
> gcloud beta emulators pubsub start --host-port=127.0.0.1:8085 --project=my-project --log-http --verbosity=debug --user-output-enabled
2. set the environment variables
> export PUBSUB_EMULATOR_HOST=127.0.0.1:8085

> export PUBSUB_PROJECT_ID=my-project
3. create topic and subscription
> pubsubctl create -t my-topic -s my-sub
4. publish a message to topic
> pubsubctl publish -t my-topic -m "my message"
5. receive message from subscription
> pubsubctl receive -s my-sub

## getting help
```
$ ./pubsubctl --help
pubsubctl v0.6.5-beta
        pubsubctl is a basic Google Cloud Platform Pub/Sub [Emulator] CLI

Usage:
  pubsubctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      create topic/subscription
  delete      delete topic and/or subscription
  help        Help about any command
  list        list topics and/or subscriptions
  publish     publish message
  receive     receive message
  test        test pubsub emulator

Flags:
      --cloud                 use cloud pubsub instead of the emulator
  -h, --help                  help for pubsubctl
      --host string           [address:port] of the emulator host, defaulting to PUBSUB_EMULATOR_HOST environment variable value (if set), ignored if 'cloud' flag is also set (default "localhost:8085")
  -p, --project string        pubsub project, defaulting to PUBSUB_PROJECT_ID environment variable value (if set) (default "default")
  -s, --subscription string   pubsub subscription
      --timeout duration      time to wait for command execution (value <=0 disables timeout) (default 5s)
  -t, --topic string          pubsub topic

Use "pubsubctl [command] --help" for more information about a command.
```

## build from source
 - linux
> GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o pubsubctl
 - macos
> GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o pubsubctl

## references
 - https://cloud.google.com/pubsub/docs/overview
 - https://cloud.google.com/pubsub/docs/emulator
 - https://cloud.google.com/sdk/gcloud/reference/beta/emulators/pubsub
