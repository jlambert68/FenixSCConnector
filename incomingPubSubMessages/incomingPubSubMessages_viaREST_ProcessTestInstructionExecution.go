package incomingPubSubMessages

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

// PullPubSubTestInstructionExecutionMessagesGcpRestApi
// Use GCP RestApi to subscribe to a PubSub-Topic
func PullPubSubTestInstructionExecutionMessagesGcpRestApi(accessTokenReceivedChannelPtr *chan bool) {

	accessTokenReceivedChannel := *accessTokenReceivedChannelPtr

	common_config.Logger.WithFields(logrus.Fields{
		"id": "e2695a29-3412-48ff-ab51-662c711fef44",
	}).Debug("Incoming 'PullPubSubTestInstructionExecutionMessagesGcpRestApi'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "e61fd7f6-95dd-4bbc-a7ae-ee8c4571174f",
	}).Debug("Outgoing 'PullPubSubTestInstructionExecutionMessagesGcpRestApi'")

	// Before Starting PubSub-receiver secure that an access token has been received
	for {
		var accessTokenReceived bool
		accessTokenReceived = <-accessTokenReceivedChannel

		if accessTokenReceived == true {
			// Continue when we got an access token
			break
		} else {

		}

	}
	/*
		// Add Access token via GCP login
		var returnMessageAckNack bool
		var returnMessageString string

		// When Connector is NOT running in GCP
		if common_config.ExecutionLocationForConnector != common_config.GCP {
			_, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(context.Background(), gcp.GetTokenForGrpcAndPubSub) //gcp.GetTokenFromWorkerForPubSub) //gcp.GenerateTokenForPubSub)
			if returnMessageAckNack == false {

				common_config.Logger.WithFields(logrus.Fields{
					"ID":                   "c96f4e4a-93b8-4175-ac2e-5b4eacd89a8f",
					"returnMessageAckNack": returnMessageAckNack,
					"returnMessageString":  returnMessageString,
				}).Error("Got some problem when generating GCP access token")

				return
			}
		}
	*/

	// Generate Subscription name to use
	var subID string
	subID = generatePubSubTopicSubscriptionNameForExecutionStatusUpdates()

	// Create a loop to be able to have a continuous PubSub Subscription Engine
	var numberOfMessagesInPullResponse int
	var err error
	var returnAckNack bool
	var returnMessage string
	var ctx context.Context

	ctx = context.Background()

	for {

		// Generate a new token is needed
		_, returnAckNack, returnMessage = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokenForPubSub)
		if returnAckNack == false {

			// Set to zero because we need some waiting time
			numberOfMessagesInPullResponse = 0

			common_config.Logger.WithFields(logrus.Fields{
				"id":            "4d4f1144-a905-4b3c-8d71-ef533eea514c",
				"returnMessage": returnMessage,
			}).Debug("Problem when generating a new token. Waiting some time before next try")

		} else {

			// Pull a certain number of messages from Subscription
			numberOfMessagesInPullResponse, err = retrievePubSubMessagesViaRestApi(subID, gcp.Gcp.GetGcpAccessTokenForAuthorizedAccountsPubSub())

			if err != nil {

				common_config.Logger.WithFields(logrus.Fields{
					"ID":  "32cdeb33-26a0-480d-98f9-ce06d13bb8aa",
					"err": err,
				}).Fatalln("PubSub receiver for TestInstructionExecutions ended, which is not intended")

			}

			// If there are more than zero messages then don't wait
			if numberOfMessagesInPullResponse == 0 {
				// Wait 15 seconds before looking for more PubSub-messages
				time.Sleep(5 * time.Second)
			}
		}

	}

}
