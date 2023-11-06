package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendConnectorProcessTestInstructionExecutionResponse
// When a TestInstructionExecution has been received by the Connector, via PubSub , then it informs the Worker
// that the Connector will execute it
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) SendConnectorProcessTestInstructionExecutionResponse(
	processTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse) (
	bool, string) {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "b7fbd3b3-8ea9-43e7-a79b-15f494e58e19",
		"processTestInstructionExecutionResponse": processTestInstructionExecutionResponse,
	}).Debug("Incoming 'SendConnectorProcessTestInstructionExecutionResponse'")

	common_config.Logger.WithFields(logrus.Fields{
		"id": "3d140810-b580-489c-b1ae-a4b0520888cb",
	}).Debug("Outgoing 'SendConnectorProcessTestInstructionExecutionResponse'")

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
			"ID": "34c7efe8-e1ab-4b6a-a945-59727f730a2e",
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

		returnMessage, err := fenixExecutionWorkerGrpcClient.ConnectorProcessTestInstructionExecutionResponse(ctx, processTestInstructionExecutionResponse)

		// Add to counter for how many gRPC-call-attempts to Worker that have been done
		gRPCCallAttemptCounter = gRPCCallAttemptCounter + 1

		// Shouldn't happen
		if err != nil {

			// Only return the error after last attempt
			if gRPCCallAttemptCounter >= numberOfgRPCCallAttempts {

				common_config.Logger.WithFields(logrus.Fields{
					"ID":    "bb37e04d-2154-47df-8eca-ea076a132a59",
					"error": err,
				}).Error("Problem to do gRPC-call to Fenix Execution Worker for 'SendConnectorProcessTestInstructionExecutionResponse'")

				return false, err.Error()
			}

			// Sleep for some time before retrying to connect
			time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenGrpcCallAttempts[gRPCCallAttemptCounter-1]))

		} else if returnMessage.AckNack == false {
			// Couldn't handle gPRC call
			common_config.Logger.WithFields(logrus.Fields{
				"ID":                        "b3174be3-af16-4ec6-a45e-ac54c8a06b53",
				"Message from Fenix Worker": returnMessage.Comments,
			}).Error("Problem to do gRPC-call to Worker for 'SendConnectorProcessTestInstructionExecutionResponse'")

			return false, returnMessage.Comments
		} else {

			common_config.Logger.WithFields(logrus.Fields{
				"ID": "1c8e6eb3-272c-4305-b966-69c7852bd60d",
				"processTestInstructionExecutionResponse": processTestInstructionExecutionResponse,
			}).Debug("Response regarding that Connector received a TestInstruction to execute was successfully sent back to Worker")

			return returnMessage.AckNack, returnMessage.Comments

		}

	}
}
