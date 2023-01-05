package gRPCServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/connectorEngine"
	"context"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// TriggerReportProcessingCapability
// Trigger Connector to inform Execution Worker of Clients capability to execute requests in parallell, serial or no processing at all(right now)
func (s *fenixExecutionConnectorGrpcServicesServer) TriggerReportProcessingCapability(ctx context.Context, emptyParameter *fenixExecutionConnectorGrpcApi.EmptyParameter) (ackNackResponse *fenixExecutionConnectorGrpcApi.AckNackResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "d85d5be5-33e8-4b8e-9577-50e4b84df389",
	}).Debug("Incoming 'gRPCServer - TriggerReportProcessingCapability'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "0a46c193-d37a-40bc-8c7b-43c1c2e02898",
	}).Debug("Outgoing 'gRPCServer - TriggerReportProcessingCapability'")

	// Calling system
	userId := "External Trigger"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectConnectorProtoFileVersion(userId, emptyParameter.ProtoFileVersionUsedByCaller)
	if returnMessage != nil {

		return returnMessage, nil
	}

	// Send Message on CommandChannel to be able to send Result back to Fenix Execution Server
	channelCommand := connectorEngine.ChannelCommandStruct{
		ChannelCommand: connectorEngine.ChannelCommandTriggerReportProcessingCapability,
		ReportCompleteTestInstructionExecutionResultParameter: connectorEngine.ChannelCommandSendReportCompleteTestInstructionExecutionResultToFenixExecutionServerStruct{
			TriggerTestInstructionExecutionResultMessage: &fenixExecutionConnectorGrpcApi.TriggerTestInstructionExecutionResultMessage{}},
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
