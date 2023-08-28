package connectorEngine

import (
	"FenixSCConnector/messagesToExecutionWorkerServer"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

type TestInstructionExecutionEngineStruct struct {
	logger                                                               *logrus.Logger
	CommandChannelReference                                              *ExecutionEngineChannelType
	MessagesToExecutionWorkerObjectReference                             *messagesToExecutionWorkerServer.MessagesToExecutionWorkerObjectStruct
	ongoingTimerOrConnectionForCallingWorkerForTestInstructionsToExecute bool
}

var TestInstructionExecutionEngine TestInstructionExecutionEngineStruct

// ExecutionEngineCommandChannel, which is references by all parts of the Connector
var ExecutionEngineCommandChannel ExecutionEngineChannelType

type ExecutionEngineChannelType chan ChannelCommandStruct

type ChannelCommandType uint8

const (
	ChannelCommandSendAreYouAliveToFenixWorkerServer ChannelCommandType = iota
	ChannelCommandTriggerReportProcessingCapability
	ChannelCommandTriggerReportCompleteTestInstructionExecutionResult
	ChannelCommandTriggerReportCurrentTestInstructionExecutionResult
	ChannelCommandTriggerSendAllLogPostForExecution
	ChannelCommandTriggerRequestForTestInstructionExecutionToProcess
	ChannelCommandTriggerRequestForTestInstructionExecutionToProcessIn5Seconds
)

type ChannelCommandStruct struct {
	ChannelCommand                                        ChannelCommandType
	ReportCompleteTestInstructionExecutionResultParameter ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct
}

// ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct
// Parameter used when to forward the final execution result for a TestInstruction
type ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct struct {
	TriggerTestInstructionExecutionResultMessage *fenixExecutionConnectorGrpcApi.TriggerTestInstructionExecutionResultMessage
}
