package connectorEngine

import (
	"FenixSCConnector/common_config"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// Channel reader which is used for reading out commands to CommandEngine
func (executionEngine *TestInstructionExecutionEngineStruct) startCommandChannelReader() {

	var incomingChannelCommand ChannelCommandStruct

	for {
		// Wait for incoming command over channel
		incomingChannelCommand = <-*executionEngine.CommandChannelReference

		switch incomingChannelCommand.ChannelCommand {

		case ChannelCommandTriggerRequestForTestInstructionExecutionToProcess:
			executionEngine.initiateConnectorRequestForProcessTestInstructionExecution()

		case ChannelCommandTriggerRequestForTestInstructionExecutionToProcessIn5Seconds:
			executionEngine.initiateConnectorRequestForProcessTestInstructionExecutionInXSeconds(1 * 1)

		// No other command is supported
		default:
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                     "6bf37452-da99-4e7e-aa6a-4627b05d1bdb",
				"incomingChannelCommand": incomingChannelCommand,
			}).Fatalln("Unknown command in CommandChannel for Worker Engine")
		}
	}

}

// Call Worker to get TestInstructions to Execute, which is done as a message stream in the response from the Worker
func (executionEngine *TestInstructionExecutionEngineStruct) initiateConnectorRequestForProcessTestInstructionExecution() {

	// Call RequestForProcessTestInstructionExecution with parameter set to zero sleep before do the gPRC-call
	executionEngine.initiateConnectorRequestForProcessTestInstructionExecutionInXSeconds(0)

}

// Call Worker in X seconds, due to some connection error, to get TestInstructions to Execute, which is done as a message stream in the response from the Worker
func (executionEngine *TestInstructionExecutionEngineStruct) initiateConnectorRequestForProcessTestInstructionExecutionInXSeconds(waitTimeInSeconds int) {

	// Only trigger time of there is none ongoing
	if executionEngine.ongoingTimerOrConnectionForCallingWorkerForTestInstructionsToExecute == true {
		return
	}

	// Run it as a go-routine
	go func() {

		// Set that there is an ongoing timer
		executionEngine.ongoingTimerOrConnectionForCallingWorkerForTestInstructionsToExecute = true

		// Wait x minutes/second before triggering
		sleepDuration := time.Duration(waitTimeInSeconds) * time.Second
		time.Sleep(sleepDuration)

		// Call Worker to get TestInstructions to Execute, but only if Worker shouldn't be side steped
		if common_config.TurnOffCallToWorker == false {
			executionEngine.MessagesToExecutionWorkerObjectReference.InitiateConnectorRequestForProcessTestInstructionExecution()

			executionEngine.ongoingTimerOrConnectionForCallingWorkerForTestInstructionsToExecute = false

			// Create Message for CommandChannel to retry to connect in 1 second
			triggerTestInstructionExecutionResultMessage := &fenixExecutionConnectorGrpcApi.TriggerTestInstructionExecutionResultMessage{}
			channelCommand := ChannelCommandStruct{
				ChannelCommand: ChannelCommandTriggerRequestForTestInstructionExecutionToProcessIn5Seconds,
				ReportCompleteTestInstructionExecutionResultParameter: ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct{
					TriggerTestInstructionExecutionResultMessage: triggerTestInstructionExecutionResultMessage},
			}

			// Send message on channel
			*executionEngine.CommandChannelReference <- channelCommand
		}

	}()

}

// Check ongoing executions  for TestInstructions for change in status that should be propagated to other places
func (executionEngine *TestInstructionExecutionEngineStruct) checkOngoingExecutionsForTestInstructions() {

}

// SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer
// Forward the final result of a TestInstructionExecution done by domains own execution engine
func (executionEngine *TestInstructionExecutionEngineStruct) SendReportCompleteTestInstructionExecutionResultToFenixExecutionServer(channelCommand ChannelCommandStruct) {
	/*
		var finalTestInstructionExecutionResultMessageFromExecutionWorker *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage
		var finalTestInstructionExecutionResultMessageToExecutionServer *fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage

		// Convert message into Worker-message-structure-type
		finalTestInstructionExecutionResultMessageFromExecutionWorker = channelCommand.ReportCompleteTestInstructionExecutionResultParameter.TriggerTestInstructionExecutionResultMessage

		// Convert from Worker-message into ExecutionServer-message
		finalTestInstructionExecutionResultMessageToExecutionServer = &fenixExecutionServerGrpcApi.FinalTestInstructionExecutionResultMessage{
			ClientSystemIdentification: &fenixExecutionServerGrpcApi.ClientSystemIdentificationMessage{
				DomainUuid:                   finalTestInstructionExecutionResultMessageFromExecutionWorker.ClientSystemIdentification.DomainUuid,
				ProtoFileVersionUsedByClient: fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum(finalTestInstructionExecutionResultMessageFromExecutionWorker.ClientSystemIdentification.ProtoFileVersionUsedByClient),
			},
			TestInstructionExecutionUuid:   finalTestInstructionExecutionResultMessageFromExecutionWorker.TestInstructionExecutionUuid,
			TestInstructionExecutionStatus: fenixExecutionServerGrpcApi.TestInstructionExecutionStatusEnum(finalTestInstructionExecutionResultMessageFromExecutionWorker.TestInstructionExecutionStatus),
		}

		// Send the result using a go-routine to be able to process next command on command-queue
		go func() {
			sendResult, errorMessage := executionEngine.MessagesToExecutionWorkerObjectReference.SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer(finalTestInstructionExecutionResultMessageToExecutionServer)

			if sendResult == false {
				executionEngine.logger.WithFields(logrus.Fields{
					"id":             "e9aae7c6-8a14-4da2-8001-2029d5bbac8d",
					"errorMessage":   errorMessage,
					"channelCommand": channelCommand,
				}).Error("Couldn't do gRPC-call to Execution Server ('SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer')")
			}
		}()


	*/
}
