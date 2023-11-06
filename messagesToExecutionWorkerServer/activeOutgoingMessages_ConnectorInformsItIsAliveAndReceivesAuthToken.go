package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendAreYouAliveToFenixExecutionServer - Ask Execution Connector to check if Worker is up and running
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) SendConnectorInformsItIsAlive(
	connectorIsReadyMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage) {

	/*
		common_config.Logger.WithFields(logrus.Fields{
			"id": "dc761d3f-f85f-4b0e-a06d-755e9b8dd352",
		}).Debug("Incoming 'SendConnectorInformsItIsAlive'")

		common_config.Logger.WithFields(logrus.Fields{
			"id": "a682cce6-4e88-4613-8d14-f579c994b4bf",
		}).Debug("Outgoing 'SendConnectorInformsItIsAlive'")
	*/

	// Before exiting
	defer func() {

	}()

	var ctx context.Context
	var returnMessageAckNack bool

	ctx = context.Background()

	// Set up connection to Server
	ctx, err := toExecutionWorkerObject.SetConnectionToFenixExecutionWorkerServer(ctx)
	if err != nil {
		return
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		/*
			toExecutionWorkerObject.Logger.WithFields(logrus.Fields{
				"ID": "c6fdb82a-6956-4943-b08c-1f6b5164531f",
			}).Debug("Running Defer Cancel function")
		*/
		cancel()

	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP &&
		common_config.GCPAuthentication == true {

		// Add Access token
		ctx, returnMessageAckNack, _ = gcp.Gcp.GenerateGCPAccessToken(
			ctx, gcp.GetTokenForGrpcAndPubSub)
		if returnMessageAckNack == false {
			return
		}

	}

	// Do the gRPC-call
	//md2 := MetadataFromHeaders(headers)
	//myctx := metadata.NewOutgoingContext(ctx, md2)

	var connectorIsReadyResponseMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyResponseMessage
	connectorIsReadyResponseMessage, err = fenixExecutionWorkerGrpcClient.ConnectorInformsItIsAlive(
		ctx, connectorIsReadyMessage)

	// Shouldn't happen
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":    "41cc0850-93c2-4e57-8baf-11144840e601",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixExecutionWorker for 'SendConnectorInformsItIsAlive'")

		return

	} else if connectorIsReadyResponseMessage.AckNackResponse.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		common_config.Logger.WithFields(logrus.Fields{
			"ID":                                  "6fcf35a5-6a8f-4b3c-a2a0-e00c9d594c73",
			"Message from Fenix Execution Server": connectorIsReadyResponseMessage.AckNackResponse.Comments,
		}).Error("Problem to do gRPC-call to FenixExecutionWorker for 'SendConnectorInformsItIsAlive'")

		return
	}

	// Store Access token to be used when doing PubSub-subscriptions
	gcp.Gcp.GcpAccessTokenFromWorkerToBeUsedWithPubSub = connectorIsReadyResponseMessage.GetPubSubAuthorizationToken()

	return

}
