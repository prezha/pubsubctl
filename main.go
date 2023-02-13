// ref :https://cloud.google.com/pubsub/docs/emulator
// ref: https://cloud.google.com/go/docs/reference/cloud.google.com/go/pubsub/latest#emulator
// ref: https://cloud.google.com/pubsub/docs/publish-receive-messages-client-library
// ref: https://pkg.go.dev/cloud.google.com/go/pubsub#section-readme
// ref: https://cloud.google.com/pubsub/docs/admin#go
// ref: https://cloud.google.com/pubsub/docs/reference/error-codes

// build:
//   GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o pubsubctl-amd64
//   GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o pubsubctl-arm64

package main

import (
	"github.com/prezha/pubsubctl/cmd"
)

func main() {
	cmd.Execute()
}
