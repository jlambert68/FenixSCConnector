package main

import (
	"FenixSCConnector/common_config"
	"github.com/sirupsen/logrus"
	"strconv"
	//"flag"
	"fmt"
	"log"
	"os"
)

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(environmentVariableName string) string {

	var environmentVariable string

	if useInjectedEnvironmentVariables == "true" {
		// Extract environment variables from parameters feed into program at compilation time

		switch environmentVariableName {
		case "AuthClientId":
			environmentVariable = authClientId

		case "AuthClientSecret":
			environmentVariable = authClientSecret

		case "CAEngineAddress":
			environmentVariable = cAEngineAddress

		case "CAEngineAddressPath":
			environmentVariable = cAEngineAddressPath

		case "ExecutionConnectorPort":
			environmentVariable = executionConnectorPort

		case "ExecutionLocationForConnector":
			environmentVariable = executionLocationForConnector

		case "ExecutionLocationForWorker":
			environmentVariable = executionLocationForWorker

		case "ExecutionWorkerAddress":
			environmentVariable = executionWorkerAddress

		case "ExecutionWorkerPort":
			environmentVariable = executionWorkerPort

		case "GCPAuthentication":
			environmentVariable = gCPAuthentication

		case "GcpProject":
			environmentVariable = gcpProject

		case "LocalServiceAccountPath":
			environmentVariable = localServiceAccountPath

		case "LoggingLevel":
			environmentVariable = loggingLevel

		case "RunInTray":
			environmentVariable = runInTray

		case "TestInstructionExecutionPubSubTopicBase":
			environmentVariable = testInstructionExecutionPubSubTopicBase

		case "ThisDomainsUuid":
			environmentVariable = thisDomainsUuid

		case "TurnOffCallToWorker":
			environmentVariable = turnOffCallToWorker

		case "UseInternalWebServerForTest":
			environmentVariable = useInternalWebServerForTest

		case "UsePubSubToReceiveMessagesFromWorker":
			environmentVariable = usePubSubToReceiveMessagesFromWorker

		case "UseServiceAccount":
			environmentVariable = useServiceAccount

		case "UseNativeGcpPubSubClientLibrary":
			environmentVariable = useNativeGcpPubSubClientLibrary

		default:
			log.Fatalf("Warning: %s environment variable not among injected variables.\n", environmentVariableName)

		}

		if environmentVariable == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", environmentVariableName)
		}

	} else {
		//
		environmentVariable = os.Getenv(environmentVariableName)
		if environmentVariable == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", environmentVariableName)
		}

	}
	return environmentVariable
}

// Variables injected at compilation time
var (
	useInjectedEnvironmentVariables string

	authClientId                            string
	authClientSecret                        string
	cAEngineAddress                         string
	cAEngineAddressPath                     string
	executionConnectorPort                  string
	executionLocationForConnector           string
	executionLocationForWorker              string
	executionWorkerAddress                  string
	executionWorkerPort                     string
	gCPAuthentication                       string
	gcpProject                              string
	localServiceAccountPath                 string
	loggingLevel                            string
	runInTray                               string
	testInstructionExecutionPubSubTopicBase string
	thisDomainsUuid                         string
	turnOffCallToWorker                     string
	useInternalWebServerForTest             string
	usePubSubToReceiveMessagesFromWorker    string
	useServiceAccount                       string
	useNativeGcpPubSubClientLibrary         string
)

func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

func init() {
	//executionLocationForConnector := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Get Environment variable to tell how/were this worker is  running
	var executionLocationForConnector = mustGetenv("ExecutionLocationForConnector")

	switch executionLocationForConnector {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForConnector = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForConnector = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForConnector = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for Connector: " + executionLocationForConnector + ". Expected one of the following: 'LOCALHOST_NODOCKER', 'LOCALHOST_DOCKER', 'GCP'")
		os.Exit(0)

	}

	// Get Environment variable to tell were Fenix Execution Server is running
	var executionLocationForExecutionWorker = mustGetenv("ExecutionLocationForWorker")

	switch executionLocationForExecutionWorker {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixExecutionWorkerServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixExecutionWorkerServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixExecutionWorkerServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for Fenix Execution Worker Server: " + executionLocationForExecutionWorker + ". Expected one of the following: 'LOCALHOST_NODOCKER', 'LOCALHOST_DOCKER', 'GCP'")
		os.Exit(0)

	}

	// Address to Fenix Execution Worker Server
	common_config.FenixExecutionWorkerAddress = mustGetenv("ExecutionWorkerAddress")

	// Port for Fenix Execution Worker Server
	common_config.FenixExecutionWorkerPort, err = strconv.Atoi(mustGetenv("ExecutionWorkerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'ExecutionWorkerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Port for Fenix Execution Connector Server
	common_config.ExecutionConnectorPort, err = strconv.Atoi(mustGetenv("ExecutionConnectorPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'executionConnectorPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Build the Dial-address for gPRC-call
	common_config.FenixExecutionWorkerAddressToDial = common_config.FenixExecutionWorkerAddress + ":" + strconv.Itoa(common_config.FenixExecutionWorkerPort)

	// Extract Debug level
	var loggingLevel = mustGetenv("LoggingLevel")

	switch loggingLevel {

	case "DebugLevel":
		common_config.LoggingLevel = logrus.DebugLevel

	case "InfoLevel":
		common_config.LoggingLevel = logrus.InfoLevel

	default:
		fmt.Println("Unknown loggingLevel '" + loggingLevel + "'. Expected one of the following: 'DebugLevel', 'InfoLevel'")
		os.Exit(0)

	}

	// Extract if there is a need for authentication when going toward GCP
	boolValue, err := strconv.ParseBool(mustGetenv("GCPAuthentication"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'GCPAuthentication:' to an boolean, error: ", err)
		os.Exit(0)
	}
	common_config.GCPAuthentication = boolValue

	// Extract if local web server for test should be used instead of FangEngine
	boolValue, err = strconv.ParseBool(mustGetenv("UseInternalWebServerForTest"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'UseInternalWebServerForTest:' to an boolean, error: ", err)
		os.Exit(0)
	}
	common_config.UseInternalWebServerForTest = boolValue

	// Extract Address, Port and url-path for Sub Custody Rest-Engine
	common_config.CAEngineAddress = mustGetenv("CAEngineAddress")
	common_config.CAEngineAddressPath = mustGetenv("CAEngineAddressPath")

	// Extract if Service Account should be used towards GCP or should the user log in via web
	boolValue, err = strconv.ParseBool(mustGetenv("UseServiceAccount"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'UseServiceAccount:' to an boolean, error: ", err)
		os.Exit(0)
	}
	common_config.UseServiceAccount = boolValue

	// Extract if there should be calls to ExecutionWorker or not. Used when testing Connector
	boolValue, err = strconv.ParseBool(mustGetenv("TurnOffCallToWorker"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'TurnOffCallToWorker:' to an boolean, error: ", err)
		os.Exit(0)
	}
	common_config.TurnOffCallToWorker = boolValue

	// Extract OAuth 2.0 Client ID
	common_config.AuthClientId = mustGetenv("AuthClientId")

	// Extract OAuth 2.0 Client Secret
	common_config.AuthClientSecret = mustGetenv("AuthClientSecret")

	// Extract the GCP-project
	common_config.GcpProject = mustGetenv("GcpProject")

	// Extract if PubSub should be used to receive messages from Worker
	common_config.UsePubSubToReceiveMessagesFromWorker, err = strconv.ParseBool(mustGetenv("UsePubSubToReceiveMessagesFromWorker"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'UsePubSubToReceiveMessagesFromWorker:' to an boolean, error: ", err)
		os.Exit(0)
	}

	// Extract the LocalServiceAccountPath
	common_config.LocalServiceAccountPath = mustGetenv("LocalServiceAccountPath")
	// The only way have an OK space is to replace an existing character
	if common_config.LocalServiceAccountPath == "#" {
		common_config.LocalServiceAccountPath = ""
	}

	// Set the environment varaible that Google-client-libraries look for
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", common_config.LocalServiceAccountPath)

	// Extract environment variable for 'TestInstructionExecutionPubSubTopicBase'
	common_config.TestInstructionExecutionPubSubTopicBase = mustGetenv("TestInstructionExecutionPubSubTopicBase")

	// Extract environment variable for 'ThisDomainsUuid'
	common_config.ThisDomainsUuid = mustGetenv("ThisDomainsUuid")

	// Extract if native pubsub client library should be used or not
	common_config.UseNativeGcpPubSubClientLibrary, err = strconv.ParseBool(mustGetenv("UseNativeGcpPubSubClientLibrary"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'UseNativeGcpPubSubClientLibrary:' to an boolean, error: ", err)
		os.Exit(0)
	}

	// Extract if a New Baseline for TestInstructions, TestInstructionContainers and Users should be saved in database
	common_config.ForceNewBaseLineForTestInstructionsAndTestInstructionContainers, err = strconv.ParseBool(
		mustGetenv("ForceNewBaseLineForTestInstructionsAndTestInstructionContainers"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable "+
			"'ForceNewBaseLineForTestInstructionsAndTestInstructionContainers:' to an boolean, error: ", err)
		os.Exit(0)
	}

}
