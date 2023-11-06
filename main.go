package main

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	"FenixSCConnector/restCallsToCAEngine"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

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
		_, returnMessageAckNack, returnMessageString := gcp.Gcp.GenerateGCPAccessToken(ctx, gcp.GetTokenForGrpcAndPubSub)
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
