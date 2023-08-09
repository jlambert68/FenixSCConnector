package incomingPubSubMessages

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"sync/atomic"
	"time"
)

func PullMsgs(w io.Writer) error {
	projectID := common_config.GcpProject
	subID := "testinstruction-execution-sub"

	var pubSubClient *pubsub.Client
	var err error
	var opts []grpc.DialOption

	ctx := context.Background()

	// Add Access token
	var returnMessageAckNack bool
	var returnMessageString string

	ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokenForPubSub)
	if returnMessageAckNack == false {
		return errors.New(returnMessageString)
	}

	//When running on GCP then use credential otherwise not
	if true { //common_config.ExecutionLocationForWorker == common_config.GCP {

		var creds credentials.TransportCredentials
		creds = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}

		pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithGRPCDialOption(opts[0]))

	} else {

		pubSubClient, err = pubsub.NewClient(ctx, projectID)
	}

	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer pubSubClient.Close()

	clientSubscription := pubSubClient.Subscription(subID)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var received int32
	err = clientSubscription.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		atomic.AddInt32(&received, 1)
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("clientSubscription.Receive: %w", err)
	}
	fmt.Fprintf(w, "Received %d messages\n", received)

	return nil
}
