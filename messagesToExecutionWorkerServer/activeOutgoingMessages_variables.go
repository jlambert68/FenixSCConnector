package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/gcp"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type MessagesToExecutionWorkerObjectStruct struct {
	Logger *logrus.Logger
	//GcpAccessToken *oauth2.Token
	Gcp gcp.GcpObjectStruct
	//CommandChannelReference *connectorEngine.ExecutionEngineChannelType
	connectionToWorkerInitiated bool
}

// Variables used for contacting Fenix Execution Worker Server
var (
	remoteFenixExecutionWorkerServerConnection *grpc.ClientConn
	FenixExecutionWorkerAddressToDial          string
	fenixExecutionWorkerGrpcClient             fenixExecutionWorkerGrpcApi.FenixExecutionWorkerConnectorGrpcServicesClient
)
