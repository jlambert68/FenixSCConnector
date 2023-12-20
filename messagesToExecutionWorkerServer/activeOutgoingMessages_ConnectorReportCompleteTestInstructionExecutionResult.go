package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer - When a TestInstruction has been fully executed the Client use this to inform the results of the execution result to the Worker (who the forward the message to the Execution Server)
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer(finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage) (bool, string) {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "db3419cd-f18b-4efa-b417-0d44cbf613e8",
		"finalTestInstructionExecutionResultMessage": finalTestInstructionExecutionResultMessage,
	}).Debug("Incoming 'SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer'")

	common_config.Logger.WithFields(logrus.Fields{
		"id": "11816322-9f48-4cfb-9329-49ab3a887d9d",
	}).Debug("Outgoing 'SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	ctx = context.Background()

	// Set up connection to WorkerServer, if that is not already done
	if toExecutionWorkerObject.connectionToWorkerInitiated == false {
		_, err := toExecutionWorkerObject.SetConnectionToFenixExecutionWorkerServer(ctx)
		if err != nil {
			return false, err.Error()
		}
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"ID": "5f02b94f-b07d-4bd7-9607-89cf712824c9",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP && common_config.GCPAuthentication == true {

		// Add Access token
		ctx, returnMessageAckNack, returnMessageString = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GetTokenForGrpcAndPubSub)
		if returnMessageAckNack == false {
			return false, returnMessageString
		}

	}

	// Creates a new temporary client only to be used for this call
	var tempFenixExecutionWorkerGrpcClient fenixExecutionWorkerGrpcApi.FenixExecutionWorkerConnectorGrpcServicesClient
	tempFenixExecutionWorkerGrpcClient = fenixExecutionWorkerGrpcApi.NewFenixExecutionWorkerConnectorGrpcServicesClient(
		remoteFenixExecutionWorkerServerConnection)

	// Do gRPC-call
	returnMessage, err := tempFenixExecutionWorkerGrpcClient.ConnectorReportCompleteTestInstructionExecutionResult(
		ctx, finalTestInstructionExecutionResultMessage)

	// Shouldn't happen
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":    "ebe601e0-14b9-42c5-8f8f-960acec80433",
			"error": err,
		}).Error("Problem to do gRPC-call to Fenix Execution Worker for 'SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer'")

		// Set that a new connection needs to be done next time
		toExecutionWorkerObject.connectionToWorkerInitiated = false

		return false, err.Error()

	} else if returnMessage.AckNack == false {
		// FenixTestDataSyncServer couldn't handle gPRC call
		common_config.Logger.WithFields(logrus.Fields{
			"ID":                               "02fe4ebe-c439-41f3-97be-a9f128bc56aa",
			"Message from Fenix Worker Server": returnMessage.Comments,
		}).Error("Problem to do gRPC-call to Worker Server for 'SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer'")

		return false, returnMessage.Comments
	}

	return returnMessage.AckNack, returnMessage.Comments

}
