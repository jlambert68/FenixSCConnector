package main

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"FenixSCConnector/restCallsToCAEngine"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"strconv"
	"time"

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
		case "RunInTray":
			environmentVariable = runInTray
		case "LoggingLevel":
			environmentVariable = loggingLevel

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
			environmentVariable = gcpAuthentication

		case "CAEngineAddress":
			environmentVariable = caEngineAddress

		case "CAEngineAddressPath":
			environmentVariable = caEngineAddressPath

		case "UseInternalWebServerForTest":
			environmentVariable = useInternalWebServerForTest

		case "UseServiceAccount":
			environmentVariable = useServiceAccount

		case "TurnOffCallToWorker":
			environmentVariable = turnOffCallToWorker

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
	runInTray                       string
	loggingLevel                    string
	executionConnectorPort          string
	executionLocationForConnector   string
	executionLocationForWorker      string
	executionWorkerAddress          string
	executionWorkerPort             string
	gcpAuthentication               string
	caEngineAddress                 string
	caEngineAddressPath             string
	useInternalWebServerForTest     string
	useServiceAccount               string
	turnOffCallToWorker             string
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

func main() {

	// Initiate logger in common_config
	InitLogger("")

	// When Execution Worker runs on GCP, then set up access
	if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP &&
		common_config.GCPAuthentication == true &&
		common_config.TurnOffCallToWorker == false {
		gcp.Gcp = gcp.GcpObjectStruct{}

		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

		// Generate first time Access token
		_, returnMessageAckNack, returnMessageString := gcp.Gcp.GenerateGCPAccessToken(ctx)
		if returnMessageAckNack == false {

			// If there was any problem then exit program
			common_config.Logger.WithFields(logrus.Fields{
				"id": "20c90d94-eef7-4819-ba8c-b7a56a39f995",
			}).Fatalf("Couldn't generate access token for GCP, return message: '%s'", returnMessageString)

		}
	}

	// InitiateRestCallsToCAEngine()
	restCallsToCAEngine.InitiateRestCallsToCAEngine()

	// If local web server, used for testing, should be used instead of FangEngine
	if common_config.UseInternalWebServerForTest == true {

		common_config.Logger.WithFields(logrus.Fields{
			"id": "353930b1-5c6f-4826-955c-19f543e2ab85",
		}).Info("Using internal web server instead of FangEngine, for RestCall")

		// Run local test web server in a go-routine
		go func() {
			restCallsToCAEngine.RestAPIServer()
		}()

	}

	// Start Connector Engine
	fenixExecutionConnectorMain()

	/*

		// Run as console program and exit as on standard exiting signals
		sig := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sig
			fmt.Println()
			fmt.Println(sig)
			done <- true

			fmt.Println("ctrl+c")
		}()

		fmt.Println("awaiting signal")
		<-done
		fmt.Println("exiting")

	*/

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

}
