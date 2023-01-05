package gRPCServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/connectorEngine"
	"context"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// TriggerReportCurrentTestInstructionExecutionResult
// Trigger Connector to inform Worker of the current execution results for an execution
func (s *fenixExecutionConnectorGrpcServicesServer) TriggerReportCurrentTestInstructionExecutionResult(ctx context.Context, triggerTestInstructionExecutionResultMessage *fenixExecutionConnectorGrpcApi.TriggerTestInstructionExecutionResultMessage) (ackNackResponse *fenixExecutionConnectorGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "515de0f2-62b9-4c65-bd39-e256960b1409",
	}).Debug("Incoming 'gRPCServer - TriggerReportCurrentTestInstructionExecutionResult'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "96730b58-04fe-42a7-945e-f112870b430a",
	}).Debug("Outgoing 'gRPCServer - TriggerReportCurrentTestInstructionExecutionResult'")

	// Calling system
	userId := "External Trigger"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectConnectorProtoFileVersion(userId, triggerTestInstructionExecutionResultMessage.ProtoFileVersionUsedByCaller)
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Send Message on CommandChannel to be able to send Result back to Fenix Execution Server
	channelCommand := connectorEngine.ChannelCommandStruct{
		ChannelCommand: connectorEngine.ChannelCommandTriggerReportCurrentTestInstructionExecutionResult,
		ReportCompleteTestInstructionExecutionResultParameter: connectorEngine.ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct{
			TriggerTestInstructionExecutionResultMessage: triggerTestInstructionExecutionResultMessage},
	}

	*s.CommandChannelReference <- channelCommand

	// Generate response
	ackNackResponse = &fenixExecutionConnectorGrpcApi.AckNackResponse{
		AckNack:                         true,
		Comments:                        "",
		ErrorCodes:                      nil,
		ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
	}

	return ackNackResponse, nil

}
