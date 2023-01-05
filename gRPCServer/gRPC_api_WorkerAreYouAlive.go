package gRPCServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/messagesToExecutionWorkerServer"
	"fmt"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// AreYouAlive - *********************************************************************
// Anyone can check if Fenix Execution Worker server is alive with this service, should be used to check serves for Connector
func (s *fenixExecutionConnectorGrpcServicesServer) WorkerAreYouAlive(ctx context.Context, emptyParameter *fenixExecutionConnectorGrpcApi.EmptyParameter) (*fenixExecutionConnectorGrpcApi.AckNackResponse, error) {

	s.logger.WithFields(logrus.Fields{
		"id": "dabd04a3-5357-4904-aceb-3493fa7396b6",
	}).Debug("Incoming 'gRPCServer - ConnectorAreYouAlive'")

	s.logger.WithFields(logrus.Fields{
		"id": "b9003ecf-b686-429b-b603-261f78e9c787",
	}).Debug("Outgoing 'gRPCServer - ConnectorAreYouAlive'")

	// Current user
	userID := "gRPC-api doesn't support UserId"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsCallerUsingCorrectConnectorProtoFileVersion(userID, fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(emptyParameter.ProtoFileVersionUsedByCaller))
	if returnMessage != nil {

		// Exiting
		return returnMessage, nil
	}

	// Set up instance to use for execution gPRC
	var fenixExecutionWorkerObject *messagesToExecutionWorkerServer.MessagesToExecutionWorkerObjectStruct
	fenixExecutionWorkerObject = &messagesToExecutionWorkerServer.MessagesToExecutionWorkerObjectStruct{
		Logger: s.logger,
		//GcpAccessToken: nil,
	}

	response, responseMessage := fenixExecutionWorkerObject.SendAreYouAliveToFenixExecutionServer()

	// Create Error Codes
	var errorCodes []fenixExecutionConnectorGrpcApi.ErrorCodesEnum

	ackNackResponseMessage := &fenixExecutionConnectorGrpcApi.AckNackResponse{
		AckNack:                         response,
		Comments:                        fmt.Sprintf("The response from Worker is '%s'", responseMessage),
		ErrorCodes:                      errorCodes,
		ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
	}

	return ackNackResponseMessage, nil

}
