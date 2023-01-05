package main

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/connectorEngine"
	"FenixSCConnector/gRPCServer"
	"FenixSCConnector/messagesToExecutionWorkerServer"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		fenixExecutionConnectorObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend GrpcServer Server
		fenixExecutionConnectorObject.GrpcServer.StopGrpcServer()

	}
}

func fenixExecutionConnectorMain() {

	// Set up BackendObject
	fenixExecutionConnectorObject = &fenixExecutionConnectorObjectStruct{
		logger:     nil,
		GrpcServer: &gRPCServer.FenixExecutionConnectorGrpcObjectStruct{},
		TestInstructionExecutionEngine: connectorEngine.TestInstructionExecutionEngineStruct{
			MessagesToExecutionWorkerObjectReference: &messagesToExecutionWorkerServer.MessagesToExecutionWorkerObjectStruct{
				//GcpAccessToken: nil,
			},
		},
	}

	// Init logger
	//fenixExecutionConnectorObject.InitLogger(loggerFileName)
	fenixExecutionConnectorObject.logger = common_config.Logger

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Initiate CommandChannel
	connectorEngine.ExecutionEngineCommandChannel = make(chan connectorEngine.ChannelCommandStruct)

	// Start ChannelCommand Engine
	fenixExecutionConnectorObject.TestInstructionExecutionEngine.CommandChannelReference = &connectorEngine.ExecutionEngineCommandChannel
	fenixExecutionConnectorObject.TestInstructionExecutionEngine.InitiateTestInstructionExecutionEngineCommandChannelReader(connectorEngine.ExecutionEngineCommandChannel)

	// Initiate  gRPC-server
	fenixExecutionConnectorObject.GrpcServer.InitiategRPCObject(fenixExecutionConnectorObject.logger)

	// Create Message for CommandChannel to connect to Worker to be able to get TestInstructions to Execute
	triggerTestInstructionExecutionResultMessage := &fenixExecutionConnectorGrpcApi.TriggerTestInstructionExecutionResultMessage{}
	channelCommand := connectorEngine.ChannelCommandStruct{
		ChannelCommand: connectorEngine.ChannelCommandTriggerRequestForTestInstructionExecutionToProcess,
		ReportCompleteTestInstructionExecutionResultParameter: connectorEngine.ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct{
			TriggerTestInstructionExecutionResultMessage: triggerTestInstructionExecutionResultMessage},
	}

	// Send message on channel
	connectorEngine.ExecutionEngineCommandChannel <- channelCommand

	// Start Backend GrpcServer-server
	fenixExecutionConnectorObject.GrpcServer.InitGrpcServer(fenixExecutionConnectorObject.logger)

}
