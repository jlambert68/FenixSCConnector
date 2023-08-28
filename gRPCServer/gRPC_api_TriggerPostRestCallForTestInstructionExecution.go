package gRPCServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/messagesToExecutionWorkerServer"
	"FenixSCConnector/restCallsToCAEngine"
	"context"
	"encoding/json"
	fenixExecutionConnectorGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionConnectorGrpcApi/go_grpc_api"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	_ "github.com/jlambert68/FenixTestInstructionsDataAdmin/Domains"
	"github.com/sirupsen/logrus"
)

// TriggerPostRestCallForTestInstructionExecution
// Trigger Connector to do RestCall with message sent with this request, used for testing
func (s *fenixExecutionConnectorGrpcServicesServer) TriggerPostRestCallForTestInstructionExecution(ctx context.Context, processTestInstructionExecutionReveredRequest *fenixExecutionConnectorGrpcApi.ProcessTestInstructionExecutionReveredRequest) (triggerPostRestCallForTestInstructionExecutionResponse *fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse, err error) {

	s.logger.WithFields(logrus.Fields{
		"id": "69274850-6595-4664-81b3-86bc21904c9d",
	}).Debug("Incoming 'gRPCServer - TriggerPostRestCallForTestInstructionExecution'")

	defer s.logger.WithFields(logrus.Fields{
		"id": "c9bf391c-1d64-4a7b-a8dd-1a95c85548d5",
	}).Debug("Outgoing 'gRPCServer - TriggerPostRestCallForTestInstructionExecution'")

	// Calling system
	userId := "External Trigger"

	// Check if Client is using correct proto files version
	var ackNackResponse *fenixExecutionWorkerGrpcApi.AckNackResponse
	ackNackResponse = common_config.IsCallerUsingCorrectWorkerProtoFileVersion(userId, fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(processTestInstructionExecutionReveredRequest.ProtoFileVersionUsedByClient))
	if ackNackResponse != nil {

		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         ackNackResponse.AckNack,
				Comments:                        ackNackResponse.Comments,
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(ackNackResponse.ProtoFileVersionUsedByClient),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// *********************************************
	//Convert fenixExecutionConnectorGrpcApi.ProcessTestInstructionExecutionReveredRequest into json
	responseBodydata, err := json.Marshal(processTestInstructionExecutionReveredRequest)
	if err != nil {
		// Problem when converting into json
		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         false,
				Comments:                        "Problem when converting 'processTestInstructionExecutionReveredRequest' into json: " + err.Error(),
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// Convert json into fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest
	var processTestInstructionExecutionReveredRequestWorkerVersion fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest
	err = json.Unmarshal(responseBodydata, &processTestInstructionExecutionReveredRequestWorkerVersion)
	if err != nil {
		// Problem when converting into json
		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         false,
				Comments:                        "Problem when converting 'json for processTestInstructionExecutionReveredRequest' into Worker-version: " + err.Error(),
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// *********************************************

	// Convert 'TestInstruction' into useful structure later to be used by FangEngine
	var processTestInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse
	var fangEngineRestApiMessageValues *restCallsToCAEngine.FangEngineRestApiMessageStruct
	processTestInstructionExecutionReversedResponse, _, fangEngineRestApiMessageValues = messagesToExecutionWorkerServer.ConvertTestInstructionIntoFangEngineStructure(&processTestInstructionExecutionReveredRequestWorkerVersion)

	// *********************************************

	//Convert fenixExecutionWorkerGrpcApi.processTestInstructionExecutionReversedResponse into json
	responseBodydata, err = json.Marshal(processTestInstructionExecutionReversedResponse)
	if err != nil {
		// Problem when converting into json
		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         false,
				Comments:                        "Problem when converting 'processTestInstructionExecutionReversedResponse' into json: " + err.Error(),
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// Convert json into fenixExecutionConnectorGrpcApi.processTestInstructionExecutionReversedResponse
	var processTestInstructionExecutionReversedResponseConnectorVersion fenixExecutionConnectorGrpcApi.ProcessTestInstructionExecutionReversedResponse
	err = json.Unmarshal(responseBodydata, &processTestInstructionExecutionReversedResponseConnectorVersion)
	if err != nil {
		// Problem when converting into json
		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         false,
				Comments:                        "Problem when converting 'json for processTestInstructionExecutionReversedResponse' into connector-version: " + err.Error(),
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// *********************************************

	// Send TestInstruction to FangEngine using RestCall
	var finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage
	finalTestInstructionExecutionResultMessage = messagesToExecutionWorkerServer.SendTestInstructionToFangEngineUsingRestCall(fangEngineRestApiMessageValues, &processTestInstructionExecutionReveredRequestWorkerVersion)

	// *********************************************

	//Convert fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage into json
	responseBodydata, err = json.Marshal(finalTestInstructionExecutionResultMessage)
	if err != nil {
		// Problem when converting into json
		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         false,
				Comments:                        "Problem when converting 'finalTestInstructionExecutionResultMessage' into json: " + err.Error(),
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// Convert json into fenixExecutionConnectorGrpcApi.FinalTestInstructionExecutionResultMessage
	var finalTestInstructionExecutionResultMessageConnectorVersion fenixExecutionConnectorGrpcApi.FinalTestInstructionExecutionResultMessage
	err = json.Unmarshal(responseBodydata, &finalTestInstructionExecutionResultMessageConnectorVersion)
	if err != nil {
		// Problem when converting into json
		triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
			AckNackResponse: &fenixExecutionConnectorGrpcApi.AckNackResponse{
				AckNack:                         false,
				Comments:                        "Problem when converting 'json for finalTestInstructionExecutionResultMessage' into connector-version: " + err.Error(),
				ErrorCodes:                      nil,
				ProtoFileVersionUsedByConnector: fenixExecutionConnectorGrpcApi.CurrentFenixExecutionConnectorProtoFileVersionEnum(common_config.GetHighestConnectorProtoFileVersion()),
			},
			ProcessTestInstructionExecutionReversedResponse: nil,
			FinalTestInstructionExecutionResultMessage:      nil,
		}

		return triggerPostRestCallForTestInstructionExecutionResponse, nil
	}

	// *********************************************

	// Generate response
	triggerPostRestCallForTestInstructionExecutionResponse = &fenixExecutionConnectorGrpcApi.TriggerPostRestCallForTestInstructionExecutionResponse{
		AckNackResponse: nil,
		ProcessTestInstructionExecutionReversedResponse: &processTestInstructionExecutionReversedResponseConnectorVersion,
		FinalTestInstructionExecutionResultMessage:      &finalTestInstructionExecutionResultMessageConnectorVersion,
	}

	return triggerPostRestCallForTestInstructionExecutionResponse, nil

}
