package main

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/connectorEngine"
	"FenixSCConnector/gRPCServer"
	"FenixSCConnector/incomingPubSubMessages"
	"FenixSCConnector/messagesToExecutionWorkerServer"
	"fmt"
	uuidGenerator "github.com/google/uuid"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"log"
	"os"
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

	// Create Unique Uuid for run time instance used as identification when communication with GuiExecutionServer
	common_config.ApplicationRunTimeUuid = uuidGenerator.New().String()
	fmt.Println("sharedCode.ApplicationRunTimeUuid: " + common_config.ApplicationRunTimeUuid)

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

	err := incomingPubSubMessages.PullMsgs(os.Stdout)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

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
