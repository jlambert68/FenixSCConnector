package incomingPubSubMessages

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

// PullPubSubTestInstructionExecutionMessages
// Use GCP Client Library to subscribe to a PubSub-Topic
func PullPubSubTestInstructionExecutionMessages(accessTokenReceivedChannelPtr *chan bool) {
	projectID := common_config.GcpProject
	subID := generatePubSubTopicSubscriptionNameForExecutionStatusUpdates()

	accessTokenReceivedChannel := *accessTokenReceivedChannelPtr

	common_config.Logger.WithFields(logrus.Fields{
		"id": "e2695a29-3412-48ff-ab51-662c711fef44",
	}).Debug("Incoming 'PullPubSubTestInstructionExecutionMessages'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "e61fd7f6-95dd-4bbc-a7ae-ee8c4571174f",
	}).Debug("Outgoing 'PullPubSubTestInstructionExecutionMessages'")

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

	var pubSubClient *pubsub.Client
	var err error
	var opts []grpc.DialOption

	ctx := context.Background()

	// Should Service Account Key be used
	if len(common_config.LocalServiceAccountPath) != 0 {
		// Use Service Account
		pubSubClient, err = pubsub.NewClient(ctx, projectID)

	} else {

		// No Service Account

		// Add Access token via GCP login
		var returnMessageAckNack bool
		var returnMessageString string

		// When Connector is NOT running in GCP
		if common_config.ExecutionLocationForConnector != common_config.GCP {
			ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GenerateTokenForPubSub) //gcp.GetTokenFromWorkerForPubSub) //gcp.GenerateTokenForPubSub)
			if returnMessageAckNack == false {

				common_config.Logger.WithFields(logrus.Fields{
					"ID":                   "53c325cd-f1bb-4b2e-bed9-106ff5b00e94",
					"returnMessageAckNack": returnMessageAckNack,
					"returnMessageString":  returnMessageString,
				}).Error("Got some problem when generating GCP access token")

				return
			}
		}
	}

	//When running on GCP then use credential otherwise not
	if !true { //common_config.ExecutionLocationForWorker == common_config.GCP {

		var creds credentials.TransportCredentials
		creds = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}

		pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithGRPCDialOption(opts[0]))

	} else {
		pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithoutAuthentication())
	}

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "a606a500-7e76-43a5-876d-c51ff8a683ec",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return
	}
	defer pubSubClient.Close()
	/*
		var pubSubTopics []*pubsub.Topic

		// Get Topics
		var pubSubTopicIterator *pubsub.TopicIterator
		pubSubTopicIterator = pubSubClient.Topics(ctx)
		for {
			var pubSubTopic *pubsub.Topic
			pubSubTopic, err = pubSubTopicIterator.Next()
			if errors.Is(err, iterator.Done) {

				// Clear the error before leaving
				err = nil

				break
			}
			if err != nil {

				common_config.Logger.WithFields(logrus.Fields{
					"ID":  "2029f0b4-be98-4057-adf9-911147adfce1",
					"err": err,
				}).Error("Got some problem iterating the topics-response")

				break
			}
			pubSubTopics = append(pubSubTopics, pubSubTopic)
		}

		fmt.Println(pubSubTopics)

	*/

	// Get all topics
	var pubSubTopics []*pubsub.Topic

	// Get Topics
	var pubSubTopicIterator *pubsub.TopicIterator
	pubSubTopicIterator = pubSubClient.Topics(ctx)
	for {
		var pubSubTopic *pubsub.Topic
		pubSubTopic, err = pubSubTopicIterator.Next()
		if errors.Is(err, iterator.Done) {

			// Clear the error before leaving
			err = nil

			break
		}
		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "2029f0b4-be98-4057-adf9-911147adfce1",
				"err": err,
			}).Error("Got some problem iterating the topics-response")

			break
		}
		pubSubTopics = append(pubSubTopics, pubSubTopic)
	}
	common_config.Logger.WithFields(logrus.Fields{
		"ID":           "2029f0b4-be98-4057-adf9-911147adfce1",
		"pubSubTopics": pubSubTopics,
	}).Info("the topics-response")

	x, err := pubSubClient.Subscription(subID).IAM().Policy(ctx)
	fmt.Println(x, err)

	creds, err := google.FindDefaultCredentials(ctx, pubsub.ScopePubSub)
	if err != nil {
		// Handle error.
		fmt.Println(creds)
	}

	//client, err := pubsub.NewClient(ctx, projectID, option.WithCredentials(creds))
	//if err != nil {
	//	fmt.Println(err)
	//}

	// Set up en PubSub-client-subscription
	clientSubscription := pubSubClient.Subscription(subID)

	perms, err := clientSubscription.IAM().TestPermissions(ctx, []string{
		"pubsub.subscriptions.consume",
		"pubsub.subscriptions.update",
	})
	fmt.Println(perms)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	//ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()

	exists, err := clientSubscription.Exists(ctx)
	fmt.Println(exists, err)

	config, err := clientSubscription.Config(ctx)
	fmt.Println(config, err)

	var numberOfMessagesInPullResponse int

	for {

		numberOfMessagesInPullResponse, err = retrievePubSubMessagesViaRestApi(subID, gcp.Gcp.GetGcpAccessTokenForAuthorizedAccountsPubSub())

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "856533ec-5ba9-46ff-b8c5-af7f3a9da2ac",
				"err": err,
			}).Fatalln("PubSub receiver for TestInstructionExecutions ended, which is not intended")

		}

		// If there are more than zero messages then don't wait
		if numberOfMessagesInPullResponse == 0 {
			// Wait 15 seconds before looking for more PubSub-messages
			time.Sleep(15 * time.Second)
		}

	}
	/*
		var received int32


		err = clientSubscription.Receive(ctx, func(_ context.Context, pubSubMessage *pubsub.Message) {

			common_config.Logger.WithFields(logrus.Fields{
				"ID": "fe01e83c-1316-4927-9970-a65f9db93c13",
			}).Debug(fmt.Printf("Got message: %q", string(pubSubMessage.Data)))

			atomic.AddInt32(&received, 1)

			// Remove any unwanted characters
			// Remove '\n'
			var cleanedMessage string
			var cleanedMessageAsByteArray []byte
			var pubSubMessageAsString string

			pubSubMessageAsString = string(pubSubMessage.Data)
			cleanedMessage = strings.ReplaceAll(pubSubMessageAsString, "\n", "")

			// Replace '\"' with '"'
			cleanedMessage = strings.ReplaceAll(cleanedMessage, "\\\"", "\"")

			cleanedMessage = strings.ReplaceAll(cleanedMessage, " ", "")

			// Convert back into byte-array
			cleanedMessageAsByteArray = []byte(cleanedMessage)

			// Convert PubSub-message back into proto-message
			var processTestInstructionExecutionPubSubRequest fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionPubSubRequest
			err = protojson.Unmarshal(cleanedMessageAsByteArray, &processTestInstructionExecutionPubSubRequest)
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":                         "bb8e4c1c-12d9-4d19-b77c-165dd05fd4eb",
					"Error":                      err,
					"string(pubSubMessage.Data)": string(pubSubMessage.Data),
				}).Error("Something went wrong when converting 'PubSub-message into proto-message")

				// Drop this message, without sending 'Ack'
				return
			}

			// Convert into Message used by converter which is the message from reversed request service
			var processTestInstructionExecutionReveredRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest
			processTestInstructionExecutionReveredRequest = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest{
				ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
					processTestInstructionExecutionPubSubRequest.GetProtoFileVersionUsedByClient()),
				TestInstruction: &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest_TestInstructionExecutionMessage{
					TestInstructionExecutionUuid: processTestInstructionExecutionPubSubRequest.TestInstruction.GetTestInstructionExecutionUuid(),
					TestInstructionUuid:          processTestInstructionExecutionPubSubRequest.TestInstruction.GetTestInstructionUuid(),
					TestInstructionName:          processTestInstructionExecutionPubSubRequest.TestInstruction.GetTestInstructionName(),
					MajorVersionNumber:           processTestInstructionExecutionPubSubRequest.TestInstruction.GetMajorVersionNumber(),
					MinorVersionNumber:           processTestInstructionExecutionPubSubRequest.TestInstruction.GetMinorVersionNumber(),
					TestInstructionAttributes:    nil, // Converted below
				},
				TestData: &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest_TestDataMessage{
					TestDataSetUuid:           processTestInstructionExecutionPubSubRequest.TestData.GetTestDataSetUuid(),
					ManualOverrideForTestData: nil, // Converted below
				},
			}

			// Convert 'TestInstruction:TestInstructionAttributes'
			var tempTestInstructionAttributes []*fenixExecutionWorkerGrpcApi.
				ProcessTestInstructionExecutionReveredRequest_TestInstructionAttributeMessage

			// Loop 'TestInstructionAttributes' from PubSub-message
			for _, pubSubTestInstructionAttribute := range processTestInstructionExecutionPubSubRequest.TestInstruction.TestInstructionAttributes {
				var tempTestInstructionAttribute *fenixExecutionWorkerGrpcApi.
					ProcessTestInstructionExecutionReveredRequest_TestInstructionAttributeMessage
				tempTestInstructionAttribute = &fenixExecutionWorkerGrpcApi.
					ProcessTestInstructionExecutionReveredRequest_TestInstructionAttributeMessage{
					TestInstructionAttributeType: fenixExecutionWorkerGrpcApi.TestInstructionAttributeTypeEnum(
						pubSubTestInstructionAttribute.GetTestInstructionAttributeType()),
					TestInstructionAttributeUuid:     pubSubTestInstructionAttribute.GetTestInstructionAttributeUuid(),
					TestInstructionAttributeName:     pubSubTestInstructionAttribute.GetTestInstructionAttributeName(),
					AttributeValueAsString:           pubSubTestInstructionAttribute.GetAttributeValueAsString(),
					AttributeValueUuid:               pubSubTestInstructionAttribute.GetTestInstructionAttributeUuid(),
					TestInstructionAttributeTypeUuid: pubSubTestInstructionAttribute.GetTestInstructionAttributeTypeUuid(),
					TestInstructionAttributeTypeName: pubSubTestInstructionAttribute.GetTestInstructionAttributeTypeName(),
				}

				// Append to slice of 'TestInstructionAttributes'
				tempTestInstructionAttributes = append(tempTestInstructionAttributes, tempTestInstructionAttribute)
			}

			processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionAttributes = tempTestInstructionAttributes

			// Convert 'TestData:ManualOverrideForTestData'
			var tempManualOverrideForTestDataSlice []*fenixExecutionWorkerGrpcApi.
				ProcessTestInstructionExecutionReveredRequest_TestDataMessage_ManualOverrideForTestDataMessage

			// Loop 'TestInstructionAttributes' from PubSub-message
			for _, pubSubManualOverrideForTestData := range processTestInstructionExecutionPubSubRequest.TestData.ManualOverrideForTestData {
				var tempManualOverrideForTestDataMessage *fenixExecutionWorkerGrpcApi.
					ProcessTestInstructionExecutionReveredRequest_TestDataMessage_ManualOverrideForTestDataMessage
				tempManualOverrideForTestDataMessage = &fenixExecutionWorkerGrpcApi.
					ProcessTestInstructionExecutionReveredRequest_TestDataMessage_ManualOverrideForTestDataMessage{
					TestDataSetAttributeUuid:  pubSubManualOverrideForTestData.GetTestDataSetAttributeUuid(),
					TestDataSetAttributeName:  pubSubManualOverrideForTestData.GetTestDataSetAttributeName(),
					TestDataSetAttributeValue: pubSubManualOverrideForTestData.GetTestDataSetAttributeValue(),
				}

				// Append to slice of 'TestInstructionAttributes'
				tempManualOverrideForTestDataSlice = append(tempManualOverrideForTestDataSlice, tempManualOverrideForTestDataMessage)
			}

			processTestInstructionExecutionReveredRequest.TestData.ManualOverrideForTestData = tempManualOverrideForTestDataSlice

			// Call 'CA' backend to convert 'TestInstruction' into useful structure later to be used by FangEngine
			var tempProcessTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse
			var fangEngineRestApiMessageValues *restCallsToCAEngine.FangEngineRestApiMessageStruct
			_, tempProcessTestInstructionExecutionResponse, fangEngineRestApiMessageValues =
				messagesToExecutionWorkerServer.ConvertTestInstructionIntoFangEngineStructure(
					processTestInstructionExecutionReveredRequest)

			// Send 'ProcessTestInstructionExecutionPubSubRequest-response' back to worker over direct gRPC-call
			couldSend, returnMessage := connectorEngine.TestInstructionExecutionEngine.
				MessagesToExecutionWorkerObjectReference.
				SendConnectorProcessTestInstructionExecutionResponse(tempProcessTestInstructionExecutionResponse)

			if couldSend == false {
				common_config.Logger.WithFields(logrus.Fields{
					"ID":            "55820706-bd18-41a6-be0a-c7d3b649e0e2",
					"returnMessage": returnMessage,
				}).Error("Couldn't send response to Worker")

			} else {

				// Send TestInstruction to FangEngine using RestCall
				var finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage
				finalTestInstructionExecutionResultMessage = messagesToExecutionWorkerServer.SendTestInstructionToFangEngineUsingRestCall(
					fangEngineRestApiMessageValues, processTestInstructionExecutionReveredRequest)

				// Send 'ProcessTestInstructionExecutionReversedResponse' back to worker over direct gRPC-call
				couldSend, returnMessage := connectorEngine.TestInstructionExecutionEngine.MessagesToExecutionWorkerObjectReference.
					SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer(finalTestInstructionExecutionResultMessage)

				if couldSend == false {
					common_config.Logger.WithFields(logrus.Fields{
						"ID": "1ce93ee2-5542-4437-9c05-d7f9d19313fa",
						"finalTestInstructionExecutionResultMessage": finalTestInstructionExecutionResultMessage,
						"returnMessage": returnMessage,
					}).Error("Couldn't send response to Worker")

				} else {

					// Send 'Ack' back to PubSub-system that message has taken care of
					pubSubMessage.Ack()
				}
			}

		})
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "2410eaa0-dce7-458b-ad9b-28d53680f995",
				"err": err,
			}).Fatalln("PubSub receiver for TestInstructionExecutions ended, which is not intended")
		}


	*/
}
