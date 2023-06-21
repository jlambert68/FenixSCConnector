package restCallsToCAEngine

import (
	"FenixSCConnector/common_config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/SubCustody/FangEngineClassesAndMethods"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/SubCustody/TestInstructions"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/TypeAndStructs"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func ConvertTestInstructionIntoFangEngineRestCallMessage(processTestInstructionExecutionReveredRequest *fenixExecutionWorkerGrpcApi.ProcessTestInstructionExecutionReveredRequest) (fangEngineRestApiMessageValues *FangEngineRestApiMessageStruct, err error) {

	// Extract TestInstructionUUID from 'TestInstructionExecutionRequest'
	testInstructionUuid := processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionUuid

	// Extract TestInstructionAttributes from 'TestInstructionExecutionRequest'
	testInstructionAttributes := processTestInstructionExecutionReveredRequest.TestInstruction.TestInstructionAttributes

	// Extract relevant FangEngineData to be used in mapping
	fangEngineData, existsInMap := fangEngineClassesMethodsAttributesMap[TypeAndStructs.OriginalElementUUIDType(testInstructionUuid)]
	if existsInMap != true {
		// Must exist in map
		common_config.Logger.WithFields(logrus.Fields{
			"id":                  "acbfdd00-ed23-4882-893c-0e6b4e61338f",
			"testInstructionUuid": testInstructionUuid,
		}).Error("Couldn't find correct FangEngineData in 'fangEngineClassesMethodsAttributesMap'")

		errorID := "4faf3e89-f647-494c-8cd2-3f0623db68c6"
		err = errors.New(fmt.Sprintf("couldn't find correct FangEngineData in 'fangEngineClassesMethodsAttributesMap' for TestInstructionUuid:'%s', [ErrorID='%s']", testInstructionUuid, errorID))

		return nil, err
	}

	// Values to be used in RestCall
	fangEngineRestApiMessageValues = &FangEngineRestApiMessageStruct{
		FangEngineClassNameNAME:  fangEngineData.FangEngineClassNameNAME,
		FangEngineMethodNameNAME: fangEngineData.FangEngineMethodNameNAME,
		FangAttributes:           make(map[TypeAndStructs.TestInstructionAttributeUUIDType]*FangEngineClassesAndMethods.FangEngineAttributesStruct),
	}

	// Loop all Attributes and populate message to be used for RestCall
	for _, testInstructionAttribute := range testInstructionAttributes {

		// Separate Attribute 'ExpectedToBePassed', which is used in url instead as a parameter in the body of the rest call
		if testInstructionAttribute.TestInstructionAttributeName != string(TestInstructions.TestInstructionAttributeName_SC_ExpectedToBePassed) {

			// Create and add Attribute with value
			var tempTestInstructionAttributesUuidAndValue TestInstructionAttributesUuidAndValueStruct
			tempTestInstructionAttributesUuidAndValue = TestInstructionAttributesUuidAndValueStruct{
				TestInstructionAttributeUUID:          TypeAndStructs.TestInstructionAttributeUUIDType(testInstructionAttribute.TestInstructionAttributeUuid),
				TestInstructionAttributeName:          TypeAndStructs.TestInstructionAttributeNameType(testInstructionAttribute.TestInstructionAttributeName),
				TestInstructionAttributeValueAsString: TypeAndStructs.AttributeValueAsStringType(testInstructionAttribute.AttributeValueAsString),
			}
			fangEngineRestApiMessageValues.TestInstructionAttribute = append(
				fangEngineRestApiMessageValues.TestInstructionAttribute,
				tempTestInstructionAttributesUuidAndValue)

			// Extract FangEngineAttribute-data
			fangEngineDataAttribute, existsInMap := fangEngineData.Attributes[TypeAndStructs.TestInstructionAttributeUUIDType(testInstructionAttribute.TestInstructionAttributeUuid)]
			if existsInMap != true {
				// Must exist in map
				common_config.Logger.WithFields(logrus.Fields{
					"id":                           "9230e70e-6054-43ec-a184-940d9519fc7d",
					"TestInstructionAttributeUuid": TypeAndStructs.TestInstructionAttributeUUIDType(testInstructionAttribute.TestInstructionAttributeUuid),
				}).Error("Couldn't find correct FangEngineData in 'fangEngineData.Attributes'")

				errorID := "0ce6aea7-ce07-4ae9-8c9f-227efb621e92"
				err = errors.New(fmt.Sprintf("couldn't find correct FangEngineData in 'fangEngineData.Attributes' for TestInstructionAttributeUuid:'%s', [ErrorID='%s']", testInstructionUuid, errorID))

				return nil, err
			}

			// Create and add reference between Attribute FangAttribute-name to be used in RestRequest
			var tempFangAttributes *FangEngineClassesAndMethods.FangEngineAttributesStruct
			tempFangAttributes = &FangEngineClassesAndMethods.FangEngineAttributesStruct{
				TestInstructionAttributeUUID:     TypeAndStructs.TestInstructionAttributeUUIDType(testInstructionAttribute.TestInstructionAttributeUuid),
				TestInstructionAttributeName:     TypeAndStructs.TestInstructionAttributeNameType(testInstructionAttribute.TestInstructionAttributeName),
				TestInstructionAttributeTypeUUID: TypeAndStructs.TestInstructionAttributeTypeUUIDType(testInstructionAttribute.TestInstructionAttributeTypeUuid),
				FangEngineAttributeNameUUID:      fangEngineDataAttribute.FangEngineAttributeNameUUID,
				FangEngineAttributeNameName:      fangEngineDataAttribute.FangEngineAttributeNameName,
			}
			// Add Attribute
			fangEngineRestApiMessageValues.FangAttributes[TypeAndStructs.TestInstructionAttributeUUIDType(testInstructionAttribute.TestInstructionAttributeUuid)] = tempFangAttributes

		} else {
			// Attribute is 'ExpectedToBePassedValue'
			fangEngineRestApiMessageValues.FangEngineExpectedToBePassedValue = TypeAndStructs.AttributeValueAsStringType(testInstructionAttribute.AttributeValueAsString)
		}

	}

	return fangEngineRestApiMessageValues, err
}

func PostTestInstructionUsingRestCall(fangEngineRestApiMessageValues *FangEngineRestApiMessageStruct) (restResponse *http.Response, err error) {

	switch fangEngineRestApiMessageValues.FangEngineMethodNameNAME {

	default:
		rand.Seed(time.Now().UnixNano())
		min := 5
		max := 10
		myRandomNumber := rand.Intn(max-min+1) + min
		fmt.Println("....sleeping for " + strconv.Itoa(myRandomNumber) + " seconds...")
		time.Sleep(time.Second * time.Duration(myRandomNumber))

	}

	fmt.Println("2. Performing Http Post...")

	common_config.Logger.WithFields(logrus.Fields{
		"id":                             "38c5fd40-0aee-4cd0-9107-0974331db0cc",
		"fangEngineRestApiMessageValues": fangEngineRestApiMessageValues,
	}).Debug("Posting TestInstruction to Sub Custodys execution Engine")

	// Generate json-body for RestCall, need to do it manually, because of strange json-structure with parameters just added instead of using an array
	attributesMap := make(map[string]string)

	for _, testInstructionAttribute := range fangEngineRestApiMessageValues.TestInstructionAttribute {

		// Add attribute to Map
		attributesMap[string(testInstructionAttribute.TestInstructionAttributeName)] = string(testInstructionAttribute.TestInstructionAttributeValueAsString)
	}

	attributesAsJson, err := json.Marshal(attributesMap)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":            "e1e74131-5040-43fa-abfc-1023f09d4388",
			"attributesMap": attributesMap,
		}).Error("Couldn't Marshal Attributes-map into json request")

		return nil, err
	}

	// Create request-url
	/*
		https://igc-SubCustodytestauto-cax-test.sebshift.dev.sebank.se/TestCaseExecution/ExecuteTestActionMethod/a/b?expectedToBePassed=true

		curl -X 'POST' \
		  'https://igc-SubCustodytestauto-cax-test.sebshift.dev.sebank.se/TestCaseExecution/ExecuteTestActionMethod/a/b?expectedToBePassed=true' \
		  -H 'accept: text/plain' \
		  -H 'Content-Type: application/json' \
		  -d '{
		  "additionalProp1": "string",
		  "additionalProp2": "string",
		  "additionalProp3": "string"
		}'
	*/
	var fangEngineUrl string
	fangEngineUrl = common_config.CAEngineAddressPath + "/" + string(fangEngineRestApiMessageValues.FangEngineClassNameNAME) +
		"/" + string(fangEngineRestApiMessageValues.FangEngineMethodNameNAME) +
		"?" + "expectedToBePassed=" + string(fangEngineRestApiMessageValues.FangEngineExpectedToBePassedValue)

	// Use Local web server for test or FangEngine
	if common_config.UseInternalWebServerForTest == true {
		// Use Local web server for testing
		fangEngineUrl = "http://" + common_config.LocalWebServerAddressAndPort + fangEngineUrl

	} else {
		// Use FangEngine
		fangEngineUrl = common_config.CAEngineAddress + fangEngineUrl
	}
	// Create Http-client
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Do RestCall to FangEngine
	restResponse, err = httpClient.Post(fangEngineUrl, "application/json; charset=utf-8", bytes.NewBuffer(attributesAsJson))
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":            "b98c2fb4-e717-4fc4-8d2c-6c791c523175",
			"fangEngineUrl": fangEngineUrl,
		}).Error("Couldn't do call to Rest-execution-server")

		return restResponse, err
	}

	return restResponse, err
}
