package common_config

import (
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
)

// IsCallerUsingCorrectConnectorProtoFileVersion ********************************************************************************************************************
// Check if Caller  is using correct proto-file version
func IsCallerUsingCorrectConnectorProtoFileVersion(callingClientUuid string, usedProtoFileVersion fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum) (returnMessage *fenixExecutionConnectorGrpcApi.AckNackResponse) {

	var callerUseCorrectProtoFileVersion bool
	var protoFileExpected fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum
	var protoFileUsed fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum

	protoFileUsed = usedProtoFileVersion
	protoFileExpected = fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(GetHighestConnectorProtoFileVersion())

	// Check if correct proto files is used
	if protoFileExpected == protoFileUsed {
		callerUseCorrectProtoFileVersion = true
	} else {
		callerUseCorrectProtoFileVersion = false
	}

	// Check if Client is using correct proto files version
	if callerUseCorrectProtoFileVersion == false {
		// Not correct proto-file version is used

		// Set Error codes to return message
		var errorCodes []fenixExecutionConnectorGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionConnectorGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionConnectorGrpcApi.ErrorCodesEnum_ERROR_WRONG_PROTO_FILE_VERSION
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixExecutionConnectorGrpcApi.AckNackResponse{
			AckNack:                         false,
			Comments:                        "Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "'",
			ErrorCodes:                      errorCodes,
			ProtoFileVersionUsedByConnector: protoFileExpected,
		}

		return returnMessage

	} else {
		return nil
	}

}

// GetHighestConnectorProtoFileVersion
// Get the highest highestConnectorProtoFileVersion for Connector-gRPC-api
func GetHighestConnectorProtoFileVersion() int32 {

	// Check if there already is a 'highestConnectorProtoFileVersion' saved, if so use that one
	if highestConnectorProtoFileVersion != -1 {
		return highestConnectorProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestConnectorProtoFileVersion = maxValue

	return highestConnectorProtoFileVersion
}

// IsCallerUsingCorrectWorkerProtoFileVersion ********************************************************************************************************************
// Check if Caller  is using correct proto-file version, used when Testing locally
func IsCallerUsingCorrectWorkerProtoFileVersion(callingClientUuid string, usedProtoFileVersion fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum) (returnMessage *fenixExecutionWorkerGrpcApi.AckNackResponse) {

	var callerUseCorrectProtoFileVersion bool
	var protoFileExpected fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum
	var protoFileUsed fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum

	protoFileUsed = usedProtoFileVersion
	protoFileExpected = fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(GetHighestExecutionWorkerProtoFileVersion())

	// Check if correct proto files is used
	if protoFileExpected == protoFileUsed {
		callerUseCorrectProtoFileVersion = true
	} else {
		callerUseCorrectProtoFileVersion = false
	}

	// Check if Client is using correct proto files version
	if callerUseCorrectProtoFileVersion == false {
		// Not correct proto-file version is used

		// Set Error codes to return message
		var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum
		var errorCode fenixExecutionWorkerGrpcApi.ErrorCodesEnum

		errorCode = fenixExecutionWorkerGrpcApi.ErrorCodesEnum_ERROR_WRONG_PROTO_FILE_VERSION
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "'",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: protoFileExpected,
		}

		return returnMessage

	} else {
		return nil
	}

}

// GetHighestExecutionWorkerProtoFileVersion
// Get the highest GetHighestExecutionWorkerProtoFileVersion for Execution Worker
func GetHighestExecutionWorkerProtoFileVersion() int32 {

	// Check if there already is a 'highestExecutionWorkerProtoFileVersion' saved, if so use that one
	if highestExecutionWorkerProtoFileVersion != -1 {
		return highestExecutionWorkerProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestExecutionWorkerProtoFileVersion = maxValue

	return highestExecutionWorkerProtoFileVersion
}
