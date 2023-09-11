package common_config

import (
	"github.com/sirupsen/logrus"
)

// Used for keeping track of the proto file versions for ExecutionServer and this Worker
var highestConnectorProtoFileVersion int32 = -1
var highestExecutionWorkerProtoFileVersion int32 = -1

// Logger that all part of the system can use
var Logger *logrus.Logger

const LocalWebServerAddressAndPort = "127.0.0.1:8080"

// Unique 'Uuid' for this running instance. Created at start up. Used as identification
var ApplicationRunTimeUuid string

// Used when Stopping Ticker for sending info to Worker that Connecotr is Ready, used when doing shut down
type StopAliveToWorkerTickerChannelStruct struct {
	ReturnChannel *chan bool
}
