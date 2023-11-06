package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/gcp"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsDataAdmin/Domains"
	"github.com/sirupsen/logrus"
	"time"
)

func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) ConnectorIsReadyToReceiveWork(
	stopSending *chan common_config.StopAliveToWorkerTickerChannelStruct,
	accessTokenReceivedChannelPtr *chan bool) {

	accessTokenReceivedChannel := *accessTokenReceivedChannelPtr

	common_config.Logger.WithFields(logrus.Fields{
		"id": "0c8aa2a2-6f3a-478e-95c8-352f28dfe488",
	}).Debug("Incoming 'ConnectorIsReadyToReceiveWork'")

	defer common_config.Logger.WithFields(logrus.Fields{
		"id": "37013017-fb87-4040-90ab-ba5990451b0f",
	}).Debug("Outgoing 'ConnectorIsReadyToReceiveWork'")

	// Create the message informing Worker that Connector is ready for Work
	var connectorIsReadyMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage
	connectorIsReadyMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage{
		ClientSystemIdentification: &fenixExecutionWorkerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: string(Domains.DomainUUID_CA),
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.
				CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		},
		ConnectorIsReady: fenixExecutionWorkerGrpcApi.ConnectorIsReadyEnum_CONNECTOR_IS_READY_TO_RECEIVE_WORK,
	}

	// Set up ticker to send I'm alive to Worker
	ticker := time.NewTicker(15 * time.Second)

	// Send a Message to Worker that Connector is Ready, every 15 seconds until program exists
	var exitForLoop bool
	var incomingChannelMessage common_config.StopAliveToWorkerTickerChannelStruct

	// Call Worker to start with
	toExecutionWorkerObject.SendConnectorInformsItIsAlive(connectorIsReadyMessage)
	// If an access token was return then start PubSub subscription receiver
	if len(gcp.Gcp.GcpAccessTokenFromWorkerToBeUsedWithPubSub) > 0 {
		accessTokenReceivedChannel <- true
	} else {
		accessTokenReceivedChannel <- false
	}

	for {
		select {
		case <-*stopSending:
			exitForLoop = true
		case <-ticker.C:
			// Call Worker
			toExecutionWorkerObject.SendConnectorInformsItIsAlive(connectorIsReadyMessage)
		}
		if exitForLoop == true {
			break
		}
	}

	ticker.Stop()
	common_config.Logger.WithFields(logrus.Fields{
		"id": "2488a817-1789-4117-936e-56e5eb3b32fe",
	}).Debug("Ticker is stopped within 'ConnectorIsReadyToReceiveWork'")

	// Inform that Ticker has been stopped
	incomingChannelMessage = <-*stopSending
	*incomingChannelMessage.ReturnChannel <- true

}

func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) ConnectorIsShuttingDown() {

	common_config.Logger.WithFields(logrus.Fields{
		"id": "58b062b7-ec76-45cb-aea8-6fdc237e6a0a",
	}).Debug("Incoming 'ConnectorIsShuttingDown'")

	common_config.Logger.WithFields(logrus.Fields{
		"id": "8ac98724-f207-45a5-83e3-a33e7c08adc3",
	}).Debug("Outgoing 'ConnectorIsShuttingDown'")

	// Create the message informing Worker that Connector is shutting down
	var connectorIsReadyMessage *fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage
	connectorIsReadyMessage = &fenixExecutionWorkerGrpcApi.ConnectorIsReadyMessage{
		ClientSystemIdentification: &fenixExecutionWorkerGrpcApi.ClientSystemIdentificationMessage{
			DomainUuid: "",
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.
				CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		},
		ConnectorIsReady: fenixExecutionWorkerGrpcApi.ConnectorIsReadyEnum_CONNECTOR_IS_SHUTTING_DOWN,
	}

	// Call Worker
	toExecutionWorkerObject.SendConnectorInformsItIsAlive(connectorIsReadyMessage)

}
