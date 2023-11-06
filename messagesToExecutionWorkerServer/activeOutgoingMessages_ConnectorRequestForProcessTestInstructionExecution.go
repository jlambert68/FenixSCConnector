package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"FenixSCConnector/restCallsToCAEngine"
	"context"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/Domains"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// InitiateConnectorRequestForProcessTestInstructionExecution
// This gPRC-methods is used when a Execution Connector needs to have its TestInstruction assignments using reverse streaming
// Execution Connector opens the gPRC-channel and assignments are then streamed back to Connector from Worker
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) InitiateConnectorRequestForProcessTestInstructionExecution() {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "c8e7cbdb-46bd-4545-a472-056fff940365",
	}).Debug("Incoming 'InitiateConnectorRequestForProcessTestInstructionExecution'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "be16c2a2-4443-4e55-8ad1-9c8478a75e12",
	}).Debug("Outgoing 'InitiateConnectorRequestForProcessTestInstructionExecution'")

	// Exit if Worker shouldn't be called
	if common_config.TurnOffCallToWorker == true {
		common_config.Logger.WithFields(logrus.Fields{
			"id": "fe86de5d-2b12-423f-a04b-549461816127",
		}).Debug("Execution Worker shouldn't be called, exit call procedure")

		return
	}

	var ctx context.Context
	var returnMessageAckNack bool

	ctx = context.Background()

	// Set up connection to Server
	ctx, err := toExecutionWorkerObject.SetConnectionToFenixExecutionWorkerServer(ctx)
	if err != nil {
		return
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background()) //, 30*time.Second)
	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"ID": "5f02b94f-b07d-4bd7-9607-89cf712824c9",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP && common_config.GCPAuthentication == true {

		// Add Access token
		ctx, returnMessageAckNack, _ = gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GetTokenForGrpcAndPubSub)
		if returnMessageAckNack == false {
			return
		}

	}

	// Set up call parameter
	emptyParameter := &fenixExecutionWorkerGrpcApi.EmptyParameter{
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion())}

	// Start up streamClient from Worker server
	streamClient, err := fenixExecutionWorkerGrpcClient.ConnectorRequestForProcessTestInstructionExecution(ctx, emptyParameter)

	// Couldn't connect to Worker
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "d9ab0434-1121-4e2e-95e7-3e1cc99656b0",
			"err": err,
		}).Error("Couldn't open streamClient from Worker Server. Will wait 1 second and try again")

		return
	}

	// Local channel to decide when Server stopped sending
	done := make(chan bool)

	// Run streamClient receiver as a go-routine
	go func() {
		for {
			processTestInstructionExecutionReveredRequest, err := streamClient.Recv()
			if err == io.EOF {
				done <- true //close(done)
				return
			}
			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"ID":  "3439f49f-d7d5-477e-9a6b-cfa5ed355bfe",
					"err": err,
				}).Error("Got some error when receiving TestInstructionExecutionsRequests from Worker, reconnect in 1 second")

				done <- true //close(done)
				return

			}

			// Check if message counts as a "keep Alive message, message is 'nil
			if processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionName == "KeepAlive" {
				// Is a keep alive message
				common_config.Logger.WithFields(logrus.Fields{
					"ID": "08b86c8d-81ba-4664-8cb5-8e53140dc870",
					"processTestInstructionExecutionReveredRequest": processTestInstructionExecutionReveredRequest,
				}).Debug("'Keep alive' message received from Worker")

			} else {
				// Is a standard TestInstruction to execute by Connector backend
				common_config.Logger.WithFields(logrus.Fields{
					"ID": "d1ea4370-3e8e-4d2b-9626-a193213e091a",
					"processTestInstructionExecutionReveredRequest": processTestInstructionExecutionReveredRequest,
				}).Debug("Receive TestInstructionExecution from Worker")

				// Send response and start processing TestInstruction in parallell
				go func() {

					// Call 'CA' backend to convert 'TestInstruction' into useful structure later to be used by FangEngine
					var fangEngineRestApiMessageValues *restCallsToCAEngine.FangEngineRestApiMessageStruct
					var processTestInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse // Used When Reversed stream-response is used
					processTestInstructionExecutionReversedResponse, _, fangEngineRestApiMessageValues =
						ConvertTestInstructionIntoFangEngineStructure(processTestInstructionExecutionReveredRequest)

					// Send 'ProcessTestInstructionExecutionReversedResponse' back to worker over direct gRPC-call
					couldSend, returnMessage := toExecutionWorkerObject.
						SendConnectorProcessTestInstructionExecutionReversedResponseToFenixWorkerServer(processTestInstructionExecutionReversedResponse)

					if couldSend == false {
						common_config.Logger.WithFields(logrus.Fields{
							"ID":            "95dddb21-0895-4016-9cb5-97ab4568f30b",
							"returnMessage": returnMessage,
						}).Error("Couldn't send response to Worker")

					} else {

						// Send TestInstruction to FangEngine using RestCall
						var finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage
						finalTestInstructionExecutionResultMessage = SendTestInstructionToFangEngineUsingRestCall(fangEngineRestApiMessageValues, processTestInstructionExecutionReveredRequest)

						// Send 'ProcessTestInstructionExecutionReversedResponse' back to worker over direct gRPC-call
						couldSend, returnMessage := toExecutionWorkerObject.SendReportCompleteTestInstructionExecutionResultToFenixWorkerServer(finalTestInstructionExecutionResultMessage)

						if couldSend == false {
							common_config.Logger.WithFields(logrus.Fields{
								"ID": "95dddb21-0895-4016-9cb5-97ab4568f30b",
								"finalTestInstructionExecutionResultMessage": finalTestInstructionExecutionResultMessage,
								"returnMessage": returnMessage,
							}).Error("Couldn't send response to Worker")
						}
					}
				}()

			}

		}
	}()

	// Server stopped sending so reconnect again in 5 seconds
	<-done

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "0b5fdb7c-91aa-4dfc-b587-7b6cef83d224",
	}).Debug("Server stopped sending so reconnect again in 5 seconds")

}

// Call 'CA' backend to convert 'TestInstruction' into useful structure later to be used by FangEngine
func ConvertTestInstructionIntoFangEngineStructure(
	processTestInstructionExecutionReveredRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest) (
	processTestInstructionExecutionReversedResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse, // Used When Reversed stream-response is used
	processTestInstructionExecutionResponse *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse, // Used when PubSub-requests are used
	fangEngineRestApiMessageValues *restCallsToCAEngine.FangEngineRestApiMessageStruct) {

	fangEngineRestApiMessageValues, err := restCallsToCAEngine.ConvertTestInstructionIntoFangEngineRestCallMessage(processTestInstructionExecutionReveredRequest)

	// Generate response depending on if the 'TestInstruction' could be converted into useful FangEngine-information or not

	// Create correct response message structure depending on if PubSub are used or not
	if common_config.UsePubSubToReceiveMessagesFromWorker == true {

		// Response from PubSub-request
		if err != nil {
			// Couldn't convert into FangEngine-messageType
			timeAtDurationEnd := time.Now()

			// Generate response message to Worker, that conversion didn't work out
			processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
				AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
					AckNack:                      false,
					Comments:                     err.Error(),
					ErrorCodes:                   nil,
					ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
				},
				TestInstructionExecutionUuid:   processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionExecutionUuid,
				ExpectedExecutionDuration:      timestamppb.New(timeAtDurationEnd),
				TestInstructionCanBeReExecuted: true,
			}
		} else {
			// Generate duration for Execution:: TODO This is only for test and should be done in another way later
			rand.Seed(time.Now().UnixNano())
			min := 180
			max := 200
			myRandomNumber := rand.Intn(max-min+1) + min

			executionDuration := time.Second * time.Duration(myRandomNumber)
			timeAtDurationEnd := time.Now().Add(executionDuration)

			// Generate OK response message to Worker
			processTestInstructionExecutionResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionResponse{
				AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
					AckNack:                      true,
					Comments:                     "",
					ErrorCodes:                   nil,
					ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
				},
				TestInstructionExecutionUuid:   processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionExecutionUuid,
				ExpectedExecutionDuration:      timestamppb.New(timeAtDurationEnd),
				TestInstructionCanBeReExecuted: false,
			}
		}

		return nil, processTestInstructionExecutionResponse, fangEngineRestApiMessageValues

	} else {

		// Response for reversed streaming
		if err != nil {
			// Couldn't convert into FangEngine-messageType
			timeAtDurationEnd := time.Now()

			// Generate response message to Worker, that conversion didn't work out
			processTestInstructionExecutionReversedResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse{
				AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
					AckNack:                      false,
					Comments:                     err.Error(),
					ErrorCodes:                   nil,
					ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
				},
				TestInstructionExecutionUuid:   processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionExecutionUuid,
				ExpectedExecutionDuration:      timestamppb.New(timeAtDurationEnd),
				TestInstructionCanBeReExecuted: true,
			}
		} else {
			// Generate duration for Execution:: TODO This is only for test and should be done in another way later
			rand.Seed(time.Now().UnixNano())
			min := 180
			max := 200
			myRandomNumber := rand.Intn(max-min+1) + min

			executionDuration := time.Second * time.Duration(myRandomNumber)
			timeAtDurationEnd := time.Now().Add(executionDuration)

			// Generate OK response message to Worker
			processTestInstructionExecutionReversedResponse = &fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReversedResponse{
				AckNackResponse: &fenixExecutionWorkerGrpcApi.AckNackResponse{
					AckNack:                      true,
					Comments:                     "",
					ErrorCodes:                   nil,
					ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
				},
				TestInstructionExecutionUuid:   processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionExecutionUuid,
				ExpectedExecutionDuration:      timestamppb.New(timeAtDurationEnd),
				TestInstructionCanBeReExecuted: false,
			}
		}

		return processTestInstructionExecutionReversedResponse, nil, fangEngineRestApiMessageValues
	}
}

func SendTestInstructionToFangEngineUsingRestCall(fangEngineRestApiMessageValues *restCallsToCAEngine.FangEngineRestApiMessageStruct, processTestInstructionExecutionReveredRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest) (finalTestInstructionExecutionResultMessage *fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage) {
	// Send TestInstruction to FangEngine using RestCall
	var restResponse *http.Response
	var err error
	restResponse, err = restCallsToCAEngine.PostTestInstructionUsingRestCall(fangEngineRestApiMessageValues)

	//**************
	/*
		defer restResponse.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(restResponse.Body)

		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(bodyBytes, &jsonMap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert response body to string
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)

		// Convert response body to Todo struct
		//var todoStruct Todo
		//json.Unmarshal(bodyBytes, &todoStruct)
		//fmt.Printf("%+v\n", todoStruct)


	*/
	//*********************

	// Convert response from restCall into 'Fenix-world-data'
	var testInstructionExecutionStatus fenixExecutionWorkerGrpcApi.TestInstructionExecutionStatusEnum
	if err != nil {
		testInstructionExecutionStatus = fenixExecutionWorkerGrpcApi.TestInstructionExecutionStatusEnum_TIE_UNEXPECTED_INTERRUPTION
	} else {
		switch restResponse.StatusCode {
		case http.StatusOK: // 200
			testInstructionExecutionStatus = fenixExecutionWorkerGrpcApi.TestInstructionExecutionStatusEnum_TIE_FINISHED_OK
		case http.StatusBadRequest: // 400 TODO use correct error
			testInstructionExecutionStatus = fenixExecutionWorkerGrpcApi.TestInstructionExecutionStatusEnum_TIE_FINISHED_NOT_OK

		case http.StatusInternalServerError: // 500
			testInstructionExecutionStatus = fenixExecutionWorkerGrpcApi.TestInstructionExecutionStatusEnum_TIE_UNEXPECTED_INTERRUPTION

		default:
			// Unhandled response code

			common_config.Logger.WithFields(logrus.Fields{
				"ID":                      "f6d86465-9a3c-4277-9730-929537f1b42b",
				"restResponse.StatusCode": restResponse.StatusCode,
			}).Error("Unhandled response from FangEngine")

			testInstructionExecutionStatus = fenixExecutionWorkerGrpcApi.TestInstructionExecutionStatusEnum_TIE_UNEXPECTED_INTERRUPTION
		}

		common_config.Logger.WithFields(logrus.Fields{
			"ID":                             "ae004aab-b900-4126-8093-a2be9238b1d7",
			"restResponse.StatusCode":        restResponse.StatusCode,
			"fangEngineRestApiMessageValues": fangEngineRestApiMessageValues,
		}).Debug("Response from doing restCall")

	}

	// Generate response message to Worker
	finalTestInstructionExecutionResultMessage = &fenixExecutionWorkerGrpcApi.FinalTestInstructionExecutionResultMessage{
		ClientSystemIdentification: &fenixExecutionWorkerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid:                   string(Domains.DomainUUID_SC),
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		},
		TestInstructionExecutionUuid:         processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionExecutionUuid,
		TestInstructionExecutionStatus:       testInstructionExecutionStatus,
		TestInstructionExecutionEndTimeStamp: timestamppb.Now(),
	}

	return finalTestInstructionExecutionResultMessage
}
