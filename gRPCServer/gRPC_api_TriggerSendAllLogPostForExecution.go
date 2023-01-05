package gRPCServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/connectorEngine"
	"context"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// TriggerSendAllLogPostForExecution
// Trigger Connector to inform Worker of all log posts that have been produced for an execution
func (s *fenixExecutionConnectorGrpcServicesServer) TriggerSendAllLogPostForExecution(ctx context.Context, triggerTestInstructionExecutionResultMessage *fenixExecutionConnectorGrpcApi.TriggerTestInstructionExecutionResultMessage) (ackNackResponse *fenixExecutionConnectorGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "8b84024c-1abe-40e9-abb6-671b7001a769",
	}).Debug("Incoming 'gRPCServer - TriggerSendAllLogPostForExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "b678cc80-bfce-4ea0-8b20-13035026aa0b",
	}).Debug("Outgoing 'gRPCServer - TriggerSendAllLogPostForExecution'")

	// Calling system
	userId := "External Trigger"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectConnectorProtoFileVersion(userId, triggerTestInstructionExecutionResultMessage.ProtoFileVersionUsedByCaller)
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Send Message on CommandChannel to be able to send Result back to Fenix Execution Server
	channelCommand := connectorEngine.ChannelCommandStruct{
		ChannelCommand: connectorEngine.ChannelCommandTriggerSendAllLogPostForExecution,
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
