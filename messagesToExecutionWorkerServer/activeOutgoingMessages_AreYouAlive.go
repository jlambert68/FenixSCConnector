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
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) SendAreYouAliveToFenixExecutionServer() (bool, string) {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "5792072c-20a9-490b-a7cf-8c4f80979552",
	}).Debug("Incoming 'SendAreYouAliveToFenixExecutionServer'")

	common_config.Logger.WithFields(logrus.Fields{
		"id": "353930b1-5c6f-4826-955c-19f543e2ab85",
	}).Debug("Outgoing 'SendAreYouAliveToFenixExecutionServer'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	ctx = context.Background()

	// Set up connection to Server
	ctx, err := toExecutionWorkerObject.SetConnectionToFenixExecutionWorkerServer(ctx)
	if err != nil {
		return false, err.Error()
	}

	// Create the message with all test data to be sent to Fenix
	emptyParameter := &fenixExecutionWorkerGrpcApi.EmptyParameter{

		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(
			common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		toExecutionWorkerObject.Logger.WithFields(logrus.Fields{
			"ID": "c5ba19bd-75ff-4366-818d-745d4d7f1a52",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP &&
		common_config.GCPAuthentication == true {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(
			ctx, gcp.GetTokenForGrpcAndPubSub)
		if returnMessageAckNack == false {
			return false, returnMessageString
		}

	}

	// Do the gRPC-call
	//md2 := MetadataFromHeaders(headers)
	//myctx := metadata.NewOutgoingContext(ctx, md2)

	returnMessage, err := fenixExecutionWorkerGrpcClient.ConnectorAreYouAlive(ctx, emptyParameter)

	// Shouldn't happen
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":    "818aaf0b-4112-4be4-97b9-21cc084c7b8b",
			"error": err,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendAreYouAliveToFenixExecutionServer'")

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		common_config.Logger.WithFields(logrus.Fields{
			"ID":                                  "2ecbc800-2fb6-4e88-858d-a421b61c5529",
			"Message from Fenix Execution Server": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to FenixExecutionServer for 'SendAreYouAliveToFenixExecutionServer'")

		return false, err.Error()
	}

	return returnMessage.AckNack, returnMessage.Comments

}
