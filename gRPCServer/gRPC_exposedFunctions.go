package gRPCServer

import (
	"github.com/sirupsen/logrus"
)

// InitiategRPCObject - Initiate local logger object
func (fenixExecutionConnectorGrpcObject *FenixExecutionConnectorGrpcObjectStruct) InitiategRPCObject(logger *logrus.Logger) {

	fenixExecutionConnectorGrpcObject.logger = logger
	//fenixExecutionConnectorGrpcObject.CommandChannelReference = commandChannelReference

}

/*
// InitiateLocalObject - Initiate local 'ExecutionConnectorGrpcObject'
func (fenixExecutionConnectorGrpcObject *FenixExecutionConnectorGrpcObjectStruct) InitiateLocalObject(inFenixExecutionConnectorGrpcObject *FenixExecutionConnectorGrpcObjectStruct) {

	fenixExecutionConnectorGrpcObject.ExecutionConnectorGrpcObject = inFenixExecutionConnectorGrpcObject
}


*/
