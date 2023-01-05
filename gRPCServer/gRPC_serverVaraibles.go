package gRPCServer

import (
	"FenixSCConnector/connectorEngine"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type FenixExecutionConnectorGrpcObjectStruct struct {
	logger *logrus.Logger
	//ExecutionConnectorGrpcObject *FenixExecutionConnectorGrpcObjectStruct
	//CommandChannelReference *connectorEngine.ExecutionEngineChannelType
}

// gRPCServer variables
var (
	fenixExecutionConnectorGrpcServer *grpc.Server
	//registerFenixExecutionConnectorGrpcServicesServer       *grpc.Server
	//registerFenixExecutionConnectorWorkerGrpcServicesServer *grpc.Server
	lis net.Listener
)

// gRPCServer Server type
type fenixExecutionConnectorGrpcServicesServer struct {
	logger                  *logrus.Logger
	CommandChannelReference *connectorEngine.ExecutionEngineChannelType
	fenixExecutionConnectorGrpcApi.UnimplementedFenixExecutionConnectorGrpcServicesServer
}
