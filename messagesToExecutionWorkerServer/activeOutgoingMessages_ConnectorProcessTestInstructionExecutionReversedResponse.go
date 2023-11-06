package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer
// When a TestInstructionExecution has been received by the Connector, via reversed streaming , then it informs the Worker
// that the Connector will execute it
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer(
	processTestInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse) (
	bool, string) {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "7deef335-37fb-462c-978c-5a97a52c207f",
		"processTestInstructionExecutionReversedResponse": processTestInstructionExecutionReversedResponse,
	}).Debug("Incoming 'SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer'")

	common_config.Logger.WithFields(logrus.Fields{
		"id": "f05c825b-16cf-4cc0-8e7a-37e375c24d17",
	}).Debug("Outgoing 'SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer'")

	var ctx context.Context
	var returnMessageAckNack bool
	var returnMessageString string

	ctx = context.Background()

	// Set up connection to Server
	ctx, err := toExecutionWorkerObject.SetConnectionToFenixExecutionWorkerServer(ctx)
	if err != nil {
		return false, err.Error()
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"ID": "209c1aaa-b5b6-4d4a-a04c-e3b328ac1eaf",
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

	// slice with sleep time, in milliseconds, between each attempt to do gRPC-call to Worker
	var sleepTimeBetweenGrpcCallAttempts []int
	sleepTimeBetweenGrpcCallAttempts = []int{100, 200, 300, 300, 500, 500, 1000, 1000, 1000, 1000} // Total: 5.9 seconds

	// Do multiple attempts to do gRPC-call to Execution Worker, when it fails
	var numberOfgRPCCallAttempts int
	var gRPCCallAttemptCounter int
	numberOfgRPCCallAttempts = len(sleepTimeBetweenGrpcCallAttempts)
	gRPCCallAttemptCounter = 0

	for {

		returnMessage, err := fenixExecutionWorkerGrpcClient.ConnectorProcessTestInstructionExecutionReversedResponse(ctx, processTestInstructionExecutionReversedResponse)

		// Add to counter for how many gRPC-call-attempts to Worker that have been done
		gRPCCallAttemptCounter = gRPCCallAttemptCounter + 1

		// Shouldn't happen
		if err != nil {

			// Only return the error after last attempt
			if gRPCCallAttemptCounter >= numberOfgRPCCallAttempts {

				common_config.Logger.WithFields(logrus.Fields{
					"ID":    "bb37e04d-2154-47df-8eca-ea076a132a59",
					"error": err,
				}).Error("Problem to do gRPC-call to Fenix Execution Worker for 'SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer'")

				return false, err.Error()
			}

			// Sleep for some time before retrying to connect
			time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenGrpcCallAttempts[gRPCCallAttemptCounter-1]))

		} else if returnMessage.AckNack == false {
			// FenixTestDataSyncServer couldn't handle gPRC call
			common_config.Logger.WithFields(logrus.Fields{
				"ID":                        "7763f7d1-9a5e-4407-b97b-0737455c6e54",
				"Message from Fenix Worker": returnMessage.Comments,
			}).Error("Problem to do gRPC-call to Worker for 'SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer'")

			return false, err.Error()
		} else {

			common_config.Logger.WithFields(logrus.Fields{
				"ID": "b48ae8cc-a145-4527-b417-b3bb815824fc",
				"processTestInstructionExecutionReversedResponse": processTestInstructionExecutionReversedResponse,
			}).Debug("Response regarding that worker received a TestInstruction to execute was successfully sent back to worker")

			return returnMessage.AckNack, returnMessage.Comments

		}

	}
}
