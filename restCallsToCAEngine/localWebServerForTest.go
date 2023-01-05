package restCallsToCAEngine

import (
	"FenixSCConnector/common_config"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"errors"
	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
)

// Structs used when converting json messages in RestAPI

// RestUserMessageStruct
type RestUserMessageStruct struct {
	ParameterMap map[string]string `json:"ParameterMap,omitempty"`
}

/*
type RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct struct {
	UserId                                 string                                                                        `protobuf:"bytes,1,opt,name=UserId,proto3" json:"UserId,omitempty"`
	PinnedTestInstructionMessages          []*fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionMessage          `protobuf:"bytes,3,rep,name=PinnedTestInstructionMessages,proto3" json:"PinnedTestInstructionMessages,omitempty"`
	PinnedTestInstructionContainerMessages []*fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionContainerMessage `protobuf:"bytes,4,rep,name=PinnedTestInstructionContainerMessages,proto3" json:"PinnedTestInstructionContainerMessages,omitempty"`
}

*/

func RestAPIServer() {
	log.Println("starting API server for Test")
	//create a new router
	router := mux.NewRouter()
	router.UseEncodedPath()
	log.Println("creating routes")
	//specify endpoints

	router.HandleFunc("/health-check", healthCheck).Methods("GET")
	router.HandleFunc("/ExampleTestStepClass/DoSomething", doSomething).Methods("POST")
	router.HandleFunc("/TestCaseExecution/ExecuteTestActionMethod/GeneralSetupTearDown/Setup", generalSetupTearDownSetup).Methods("POST")
	router.HandleFunc("/TestCaseExecution/ExecuteTestActionMethod/GeneralSetupTearDown/TearDown", generalSetupTearDownTearDown).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)
	//router.HandleFunc("/*", allOtherRoutes).Methods("POST")
	/*
		router.HandleFunc("/ExampleTestStepClass/DoSomething1{expectedToBePassed}", doSomething1).Methods("POST")
		router.HandleFunc("/ExampleTestStepClass/DoSomething2{expectedToBePassed}", doSomething2).Methods("POST")
		router.HandleFunc("/ExampleTestStepClass/DoSomethingWithTestException{expectedToBePassed}", doSomethingWithTestException).Methods("POST")
		router.HandleFunc("/ExampleTestStepClass/DoSomethingWithException{expectedToBePassed}", doSomethingWithException).Methods("POST")
	*/

	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)
}

// RestApi to check if local TestWeb-server is up and running
func notFound(w http.ResponseWriter, r *http.Request) {
	// curl --request GET localhost:8080/health-check

	common_config.Logger.WithFields(logrus.Fields{
		"id": "42c2cdca-4ce1-4802-888d-ccc6eb82996f",
	}).Debug("Incoming 'RestApi - *'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "fab7676d-c303-4b20-8980-397d7a59282e",
	}).Debug("Outgoing 'RestApi - *'")

	// Create base for response body
	var responseBody map[string]string
	responseBody = make(map[string]string)
	responseBody["type"] = "FenixSCConnector - internal Web Server"

	// Create Header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	// Create Response message
	responseBody["title"] = "Error - Not Found"
	responseBody["status"] = "404"
	responseBody["detail"] = "The address used is not valid"
	responseBody["traceId"] = "6d7c074e-2110-49ef-a45a-2a41a5a83b2c"

	responseBodydata, _ := json.Marshal(responseBody)

	fmt.Fprintf(w, string(responseBodydata))

	return

}

// RestApi to check if local TestWeb-server is up and running
func healthCheck(w http.ResponseWriter, r *http.Request) {
	// curl --request GET localhost:8080/health-check

	common_config.Logger.WithFields(logrus.Fields{
		"id": "fb3c1ecb-3da8-4d27-b1c4-16d5120e7125",
	}).Debug("Incoming 'RestApi - /health-check'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "fab7676d-c303-4b20-8980-397d7a59282e",
	}).Debug("Outgoing 'RestApi - /health-check'")

	// Create base for response body
	var responseBody map[string]string
	responseBody = make(map[string]string)
	responseBody["type"] = "FenixSCConnector - internal Web Server"

	// Create Header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Create Response message
	responseBody["title"] = "OK"
	responseBody["status"] = "200"
	responseBody["detail"] = "The Connectors Internal Test API is up and running"
	responseBody["traceId"] = "5a263f5c-8a86-4552-bc61-3b5cce52b208"

	responseBodydata, _ := json.Marshal(responseBody)

	fmt.Fprintf(w, string(responseBodydata))

	return
	// Create Response message
	fmt.Fprintf(w, "API is up and running")
}

func doSomething(w http.ResponseWriter, r *http.Request) {

	// curl -X POST localhost:8080/ExampleTestStepClass/DoSomething?expectedToBePassed=true -H 'Content-Type: application/json' -d '{"UserId":"s41797", "TestInstructionUuid":"myUuid", "TestInstructionName":"myName"}'

	common_config.Logger.WithFields(logrus.Fields{
		"id": "2472dda1-701d-4b23-8326-757e43df4af4",
	}).Debug("Incoming 'RestApi - (POST) /DoSomething")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "db318ff4-ad36-43d4-a8d4-3e0ac4ff08c6",
	}).Debug("Outgoing 'RestApi - (POST) /DoSomething'")

	// Variable where Rest-json-payload will end up in
	//jsonData := &RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct{}

	// Create base for response body
	var responseBody map[string]string
	responseBody = make(map[string]string)
	responseBody["type"] = "FenixSCConnector - internal Web Server"

	// read message body
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Println(error)
		return
	}

	// close message body
	r.Body.Close()

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract 'expectedToBePassedSlice'
	variables := r.URL.Query() //mux.Vars(r)
	expectedToBePassedSlice, existInMap := variables["expectedToBePassed"]

	// Missing parameter 'expectedToBePassedSlice'
	if existInMap == false {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Missing parameter 'expectedToBePassed'"
		responseBody["traceId"] = "15f7f628-c80e-4010-8853-66df1ffa1a59"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// Exact one parameter 'expectedToBePassed' must exist
	if len(expectedToBePassedSlice) != 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		fmt.Fprintf(w, "Parameter 'expectedToBePassed' must contain exactly one value")

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Parameter 'expectedToBePassed' must contain exactly one value"
		responseBody["traceId"] = "dcdfc951-1eb5-4ed9-8c54-5f22bb718ae7"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// Parameter 'expectedToBePassed' should be 'true' or 'false'
	if expectedToBePassedSlice[0] != "true" && expectedToBePassedSlice[0] != "false" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Parameter 'expectedToBePassed' should be 'true' or 'false'"
		responseBody["traceId"] = "2c82ed7f-18f6-4362-8ca7-a4c3602d81ac"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// 'expectedToBePassed' should be 'true'
	if expectedToBePassedSlice[0] == "true" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create Response message
		responseBody["title"] = "OK"
		responseBody["status"] = "200"
		responseBody["detail"] = "OK Test from Test Web server"
		responseBody["traceId"] = "8f374286-d692-4196-83b4-575f66c12684"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return

	}

	// 'expectedToBePassed' is 'false' - Will allways go in here
	if expectedToBePassedSlice[0] == "false" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //TODO Fang must change this

		// Create Response message
		fmt.Fprintf(w, "Not a OK Test from Test Web server")
		// Create Response message
		responseBody["title"] = "Error - Internal Server Error"
		responseBody["status"] = "500"
		responseBody["detail"] = "Not a OK Test from Test Web server"
		responseBody["traceId"] = "7f139cbd-2fb2-4ba2-9f8b-4d42faefc69f"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

	}
}

func generalSetupTearDownSetup(w http.ResponseWriter, r *http.Request) {

	// curl -X POST localhost:8080/GeneralSetupTearDown/SetupexpectedToBePassed=true -H 'Content-Type: application/json' -d '{"UserId":"s41797", "TestInstructionUuid":"myUuid", "TestInstructionName":"myName"}'

	common_config.Logger.WithFields(logrus.Fields{
		"id": "5c68e681-73c1-438f-aed5-cf7f6f4f9072",
	}).Debug("Incoming 'RestApi - (POST) /GeneralSetupTearDown/Setup'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "9195d621-eb4a-477f-8f68-1109c4aa69c1",
	}).Debug("Outgoing 'RestApi - (POST) /GeneralSetupTearDown/Setup'")

	// Variable where Rest-json-payload will end up in
	//jsonData := &RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct{}

	// Create base for response body
	var responseBody map[string]string
	responseBody = make(map[string]string)
	responseBody["type"] = "FenixSCConnector - internal Web Server"

	// read message body
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Println(error)
		return
	}

	// close message body
	r.Body.Close()

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		// Just print Incoming parameter
		common_config.Logger.WithFields(logrus.Fields{
			"id":      "12789a24-f02c-494d-b8d8-18642d7588ef",
			"jsonMap": jsonMap,
		}).Debug("Incoming Parameters")
	}

	// Extract 'expectedToBePassedSlice'
	variables := r.URL.Query() //mux.Vars(r)
	expectedToBePassedSlice, existInMap := variables["expectedToBePassed"]

	// Missing parameter 'expectedToBePassedSlice'
	if existInMap == false {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Missing parameter 'expectedToBePassed'"
		responseBody["traceId"] = "88da16af-1ddf-49ca-945b-aeb2a1da6470"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// Exact one parameter 'expectedToBePassed' must exist
	if len(expectedToBePassedSlice) != 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		fmt.Fprintf(w, "Parameter 'expectedToBePassed' must contain exactly one value")

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Parameter 'expectedToBePassed' must contain exactly one value"
		responseBody["traceId"] = "71202ccb-8dc4-45d6-a64a-fe1a775d1a73"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// Parameter 'expectedToBePassed' should be 'true' or 'false'
	if expectedToBePassedSlice[0] != "true" && expectedToBePassedSlice[0] != "false" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Parameter 'expectedToBePassed' should be 'true' or 'false'"
		responseBody["traceId"] = "955c80e1-a745-42ae-800a-1ec972fcf255"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// 'expectedToBePassed' should be 'true'
	if expectedToBePassedSlice[0] == "true" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create Response message
		responseBody["title"] = "OK"
		responseBody["status"] = "200"
		responseBody["detail"] = "OK Test from Test Web server"
		responseBody["traceId"] = "30418b44-a015-4d94-823f-49faa7622ca6"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return

	}

	// 'expectedToBePassed' is 'false' - Will always go in here
	if expectedToBePassedSlice[0] == "false" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //TODO Fang must change this

		// Create Response message
		fmt.Fprintf(w, "Not a OK Test from Test Web server")
		// Create Response message
		responseBody["title"] = "Error - Internal Server Error"
		responseBody["status"] = "500"
		responseBody["detail"] = "Not a OK Test from Test Web server"
		responseBody["traceId"] = "dc2adc95-b656-4894-a296-f0a73409aedb"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

	}
}

func generalSetupTearDownTearDown(w http.ResponseWriter, r *http.Request) {

	// curl -X POST localhost:8080/GeneralSetupTearDown/SetupexpectedToBePassed=true -H 'Content-Type: application/json' -d '{"UserId":"s41797", "TestInstructionUuid":"myUuid", "TestInstructionName":"myName"}'

	common_config.Logger.WithFields(logrus.Fields{
		"id": "38d357b9-6b5d-44f1-ae8a-02586e15fd85",
	}).Debug("Incoming 'RestApi - (POST) /GeneralSetupTearDown/TearDown'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "f5593e1f-849d-4522-8cd6-b48938190e64",
	}).Debug("Outgoing 'RestApi - (POST) /GeneralSetupTearDown/TearDown'")

	// Variable where Rest-json-payload will end up in
	//jsonData := &RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct{}

	// Create base for response body
	var responseBody map[string]string
	responseBody = make(map[string]string)
	responseBody["type"] = "FenixSCConnector - internal Web Server"

	// read message body
	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Println(error)
		return
	}

	// close message body
	r.Body.Close()

	jsonMap := make(map[string]interface{})
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		// Just print Incoming parameter
		common_config.Logger.WithFields(logrus.Fields{
			"id":      "3aefede2-ac76-4a20-b3fe-43f4ec5d1d4a",
			"jsonMap": jsonMap,
		}).Debug("Incoming Parameters")
	}

	// Extract 'expectedToBePassedSlice'
	variables := r.URL.Query() //mux.Vars(r)
	expectedToBePassedSlice, existInMap := variables["expectedToBePassed"]

	// Missing parameter 'expectedToBePassedSlice'
	if existInMap == false {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Missing parameter 'expectedToBePassed'"
		responseBody["traceId"] = "40966b3b-1b09-4405-a887-5a67d3667af1"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// Exact one parameter 'expectedToBePassed' must exist
	if len(expectedToBePassedSlice) != 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		fmt.Fprintf(w, "Parameter 'expectedToBePassed' must contain exactly one value")

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Parameter 'expectedToBePassed' must contain exactly one value"
		responseBody["traceId"] = "6f0bacce-c83f-4eaa-8817-3d23681b3997"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// Parameter 'expectedToBePassed' should be 'true' or 'false'
	if expectedToBePassedSlice[0] != "true" && expectedToBePassedSlice[0] != "false" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		// Create Response message
		responseBody["title"] = "Error - Bad Request"
		responseBody["status"] = "400"
		responseBody["detail"] = "Parameter 'expectedToBePassed' should be 'true' or 'false'"
		responseBody["traceId"] = "3a30ef44-11f0-466a-97c2-8b5611d9c8ac"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return
	}

	// 'expectedToBePassed' should be 'true'
	if expectedToBePassedSlice[0] == "true" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Create Response message
		responseBody["title"] = "OK"
		responseBody["status"] = "200"
		responseBody["detail"] = "OK Test from Test Web server"
		responseBody["traceId"] = "252b6600-0219-4b26-ba12-584ddba61e66"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

		return

	}

	// 'expectedToBePassed' is 'false' - Will always go in here
	if expectedToBePassedSlice[0] == "false" {

		// Create Header
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) //TODO Fang must change this

		// Create Response message
		fmt.Fprintf(w, "Not a OK Test from Test Web server")
		// Create Response message
		responseBody["title"] = "Error - Internal Server Error"
		responseBody["status"] = "500"
		responseBody["detail"] = "Not a OK Test from Test Web server"
		responseBody["traceId"] = "d185c85c-7a46-4caf-9b29-aff95399044f"

		responseBodydata, _ := json.Marshal(responseBody)

		fmt.Fprintf(w, string(responseBodydata))

	}
}

func extractAndValidateJsonBody(responseWriterPointer *http.ResponseWriter, httpRequest *http.Request, myInputTypeVariable interface{}) (err error) {
	// If the Content-Type header is present, check that it has the value
	// application/json. Note that we are using the gddo/httputil/header
	// package to parse and extract the value here, so the check works
	// even if the client includes additional charset or boundary
	// information in the header.
	responseWriter := *responseWriterPointer
	if httpRequest.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(httpRequest.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(responseWriter, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// Use http.MaxBytesReader to enforce a maximum read of 1MB from the
	// response body. A request body larger than that will now result in
	// Decode() returning a "http: request body too large" error.
	httpRequest.Body = http.MaxBytesReader(responseWriter, httpRequest.Body, 1048576)

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause Decode() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(httpRequest.Body)
	dec.DisallowUnknownFields()

	var p = myInputTypeVariable //RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct
	err = dec.Decode(&p)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field in our RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct struct. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(responseWriter, msg, http.StatusBadRequest)

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(responseWriter, msg, http.StatusRequestEntityTooLarge)

		// Otherwise default to logging the error and sending a 500 Internal
		// Server Error response.
		default:
			log.Println(err.Error())
			http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return err
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(responseWriter, msg, http.StatusBadRequest)
		return
	}

	//fmt.Fprintf(responseWriter, "RestSavePinnedInstructionsAndTestInstructionContainersToFenixGuiBuilderServerStruct: %+v", p)

	return nil
}
