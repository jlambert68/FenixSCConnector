package common_config

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// ExecutionLocationForConnector - Where is the Worker running
var ExecutionLocationForConnector ExecutionLocationTypeType

// ExecutionLocationForFenixExecutionWorkerServer  - Where is Fenix Execution Server running
var ExecutionLocationForFenixExecutionWorkerServer ExecutionLocationTypeType

// ExecutionLocationTypeType - Definitions for where client and Fenix Server is running
type ExecutionLocationTypeType int

// Constants used for where stuff is running
const (
	LocalhostNoDocker ExecutionLocationTypeType = iota
	LocalhostDocker
	GCP
)

// Address to Fenix Execution Worker & Execution Connector, will have their values from Environment variables at startup
var (
	FenixExecutionWorkerAddress                                     string
	FenixExecutionWorkerPort                                        int
	FenixExecutionWorkerAddressToDial                               string
	ExecutionConnectorPort                                          int
	GCPAuthentication                                               bool
	CAEngineAddress                                                 string
	CAEngineAddressPath                                             string
	UseInternalWebServerForTest                                     bool
	UseServiceAccount                                               bool
	ApplicationShouldRunInTray                                      bool
	TurnOffCallToWorker                                             bool
	AuthClientId                                                    string
	AuthClientSecret                                                string
	GcpProject                                                      string
	UsePubSubToReceiveMessagesFromWorker                            bool
	LocalServiceAccountPath                                         string
	TestInstructionExecutionPubSubTopicBase                         string
	ThisDomainsUuid                                                 string
	UseNativeGcpPubSubClientLibrary                                 bool
	ForceNewBaseLineForTestInstructionsAndTestInstructionContainers bool
)
