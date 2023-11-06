include MakeFileSecretVariables

RunGrpcGui:
	cd ~/egen_kod/go/go_workspace/src/jlambert/grpcui/standalone && grpcui -plaintext localhost:6672

filename :=
filenamePartFirst := FenixSCConnectorCrossBuild_
filenamePartLast := .exe
datetime := `date +'%y%m%d_%H%M%S'`

GenerateDateTime:
	$(eval fileName := $(filenamePartFirst)$(datetime)$(filenamePartLast))

	echo $(fileName)

GenerateTrayIcons:
	./bundleIcons.sh

BuildExeForWindows:
#	fyne-cross windows -arch=amd64 --ldflags="-X 'main.useInjectedEnvironmentVariables=true' -X 'main.runInTray=truex' -X 'main.loggingLevel=DebugLevel' -X 'main.executionConnectorPort=6672' -X 'main.executionLocationForConnector=LOCALHOST_NODOCKER' -X 'main.executionLocationForWorker=GCP' -X 'main.executionWorkerAddress=fenixexecutionworker-ca-nwxrrpoxea-lz.a.run.app' -X 'main.executionWorkerPort=443' -X 'main.gcpAuthentication=false'"
#	GOOD=windows GOARCH=amd64 go build -o FenixSCConnectorWindow.exe -ldflags="-X 'main.useInjectedEnvironmentVariables=true' -X 'main.runInTray=truex' -X 'main.loggingLevel=DebugLevel' -X 'main.executionConnectorPort=6672' -X 'main.executionLocationForConnector=LOCALHOST_NODOCKER' -X 'main.executionLocationForWorker=GCP' -X 'main.executionWorkerAddress=fenixexecutionworker-ca-nwxrrpoxea-lz.a.run.app' -X 'main.executionWorkerPort=443' -X  'main.gcpAuthentication=true' -X 'main.caEngineAddress=127.0.0.1' -X 'main.caEngineAddressPath=/"
	env GOOD=windows GOARCH=amd64 go build  -o FenixSCConnector.WindowsExe -ldflags="-X 'main.useInjectedEnvironmentVariables=true' -X 'main.runInTray=truex' -X 'main.loggingLevel=DebugLevel' -X 'main.executionConnectorPort=6672' -X 'main.executionLocationForConnector=LOCALHOST_NODOCKER' -X 'main.executionLocationForWorker=GCP' -X 'main.executionWorkerAddress=fenixexecutionworker-ca-must-be-logged-in-nwxrrpoxea-lz.a.run.app' -X 'main.executionWorkerPort=443' -X 'main.gcpAuthentication=true' -X 'main.caEngineAddress=127.0.0.1' -X 'main.caEngineAddressPath=x' -X 'main.useInternalWebServerForTest=true' -X 'main.useServiceAccount=true'" /home/jlambert/egen_kod/go/go_workspace/src/jlambert/FenixSCConnector
BuildExeForLinux:
	GOOD=linux GOARCH=amd64 go build  -o FenixSCConnector.LinuxExe -ldflags="-X 'main.useInjectedEnvironmentVariables=true' -X 'main.authClientId=$(authClientId)' -X  'main.authClientSecret=$(authClientSecret)' -X  'main.cAEngineAddress=http://127.0.0.1:5000' -X  'main.cAEngineAddressPath=/TestCaseExecution/ExecuteTestActionMethod' -X  'main.executionConnectorPort=6672' -X  'main.executionLocationForConnector=LOCALHOST_NODOCKER' -X  'main.executionLocationForWorker=GCP' -X  'main.executionWorkerAddress=$(executionWorkerAddress)' -X  'main.executionWorkerPort=443' -X  'main.gCPAuthentication=true' -X  'main.gcpProject=$(gcpProject)' -X  'main.localServiceAccountPath=#' -X  'main.loggingLevel=DebugLevel' -X  'main.runInTray=truex' -X  'main.testInstructionExecutionPubSubTopicBase=ProcessTestInstructionExecutionRequest' -X  'main.thisDomainsUuid=7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae' -X  'main.turnOffCallToWorker=false' -X  'main.useInternalWebServerForTest=true' -X  'main.usePubSubToReceiveMessagesFromWorker=true' -X  'main.useServiceAccount=false'"

CrossBuildForWindows_SEB_test:
	$(eval fileName := $(filenamePartFirst)$(datetime)$(filenamePartLast))
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o $(fileName) -ldflags="-X 'main.useInjectedEnvironmentVariables=true' -X 'main.authClientId=$(authClientId)' -X  'main.authClientSecret=$(authClientSecret)' -X  'main.cAEngineAddress=http://127.0.0.1:5000' -X  'main.cAEngineAddressPath=/TestCaseExecution/ExecuteTestActionMethod' -X  'main.executionConnectorPort=6672' -X  'main.executionLocationForConnector=LOCALHOST_NODOCKER' -X  'main.executionLocationForWorker=GCP' -X  'main.executionWorkerAddress=$(executionWorkerAddress)' -X  'main.executionWorkerPort=443' -X  'main.gCPAuthentication=true' -X  'main.gcpProject=$(gcpProject)' -X  'main.localServiceAccountPath=#' -X  'main.loggingLevel=DebugLevel' -X  'main.runInTray=truex' -X  'main.testInstructionExecutionPubSubTopicBase=ProcessTestInstructionExecutionRequest' -X  'main.thisDomainsUuid=7edf2269-a8d3-472c-aed6-8cdcc4a8b6ae' -X  'main.turnOffCallToWorker=false' -X  'main.useInternalWebServerForTest=true' -X  'main.usePubSubToReceiveMessagesFromWorker=true' -X  'main.useServiceAccount=false'" .