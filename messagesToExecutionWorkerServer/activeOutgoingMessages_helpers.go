package messagesToExecutionWorkerServer

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/grpcurl"
	"context"
	"crypto/tls"
	"encoding/base64"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
	"time"
)

// ********************************************************************************************************************

// SetConnectionToFenixExecutionWorkerServer - Set upp connection and Dial to FenixExecutionServer
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) SetConnectionToFenixExecutionWorkerServer(ctx context.Context) (_ context.Context, err error) {
	/* THis is done in the new
	var opts []grpc.DialOption

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		panic(fmt.Sprintf("cannot load root CA certs, err: %s", err))
	}

	//When running on GCP then use credential otherwise not

	if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            systemRoots,
		})

		opts = []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
		}
	}
	*/

	// slice with sleep time, in milliseconds, between each attempt to Dial to Worker
	var sleepTimeBetweenDialAttempts []int
	sleepTimeBetweenDialAttempts = []int{100, 100, 200, 200, 300, 300, 500, 500, 600, 1000} // Total: 3.6 seconds

	// Do multiple attempts to do connection to Execution Worker
	var numberOfDialAttempts int
	var dialAttemptCounter int
	numberOfDialAttempts = len(sleepTimeBetweenDialAttempts)
	dialAttemptCounter = 0

	for {

		// Set up connection to Fenix Execution Worker
		// When run on GCP, use credentials
		var newGrpcClientConnection *grpc.ClientConn
		if common_config.ExecutionLocationForFenixExecutionWorkerServer == common_config.GCP {
			// Worker runs on GCP

			if common_config.ExecutionLocationForConnector == common_config.LocalhostNoDocker {
				// Connector runs Locally
				ctx, newGrpcClientConnection = dialFromGrpcurl(ctx)
				remoteFenixExecutionWorkerServerConnection = newGrpcClientConnection
			} else {
				// Connector runs on GCP
				creds := credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				})

				var opts []grpc.DialOption
				opts = []grpc.DialOption{
					grpc.WithTransportCredentials(creds),
				}
				remoteFenixExecutionWorkerServerConnection, err = grpc.Dial(common_config.FenixExecutionWorkerAddressToDial, opts...)

			}

		} else {
			// Worker runs Local
			remoteFenixExecutionWorkerServerConnection, err = grpc.Dial(common_config.FenixExecutionWorkerAddressToDial, grpc.WithInsecure())
		}
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID": "50b59b1b-57ce-4c27-aa84-617f0cde3100",
				"common_config.FenixExecutionWorkerAddressToDial": common_config.FenixExecutionWorkerAddressToDial,
				"error message":      err,
				"dialAttemptCounter": dialAttemptCounter,
			}).Error("Did not connect to FenixExecutionServer via gRPC")

			// Add to counter for how many Dial attempts that have been done
			dialAttemptCounter = dialAttemptCounter + 1

			// Only return the error after last attempt
			if dialAttemptCounter >= numberOfDialAttempts {
				return nil, err
			}

		} else {

			// Creates a new Client
			fenixExecutionWorkerGrpcClient = fenixExecutionWorkerGrpcApi.NewFenixExecutionWorkerConnectorGrpcServicesClient(remoteFenixExecutionWorkerServerConnection)

			break
		}

		// Sleep for some time before retrying to connect
		time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenDialAttempts[dialAttemptCounter-1]))

	}

	return ctx, err
}

var (
	isUnixSocket func() bool
)

func dialFromGrpcurl(ctx context.Context) (context.Context, *grpc.ClientConn) {

	target := common_config.FenixExecutionWorkerAddressToDial

	dialTime := 10 * time.Second

	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var opts []grpc.DialOption

	var creds credentials.TransportCredentials

	var tlsConf *tls.Config

	creds = credentials.NewTLS(tlsConf)

	grpcurlUA := "FenixSCConnector"
	//if grpcurl.version == grpcurl.no_version {
	//	grpcurlUA = "grpcurl/dev-build (no version set)"
	//}

	opts = append(opts, grpc.WithUserAgent(grpcurlUA))
	//opts = append(opts, grpc.WithNoProxy())

	network := "tcp"
	if isUnixSocket != nil && isUnixSocket() {
		network = "unix"
	}

	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		log.Panicln("Failed to Dial, ", target, err.Error())
	}
	return ctx, cc

}

// MetadataFromHeaders converts a list of header strings (each string in
// "Header-Name: Header-Value" form) into metadata. If a string has a header
// name without a value (e.g. does not contain a colon), the value is assumed
// to be blank. Binary headers (those whose names end in "-bin") should be
// base64-encoded. But if they cannot be base64-decoded, they will be assumed to
// be in raw form and used as is.
func MetadataFromHeaders(headers []string) metadata.MD {
	md := make(metadata.MD)
	for _, part := range headers {
		if part != "" {
			pieces := strings.SplitN(part, ":", 2)
			if len(pieces) == 1 {
				pieces = append(pieces, "") // if no value was specified, just make it "" (maybe the header value doesn't matter)
			}
			headerName := strings.ToLower(strings.TrimSpace(pieces[0]))
			val := strings.TrimSpace(pieces[1])
			if strings.HasSuffix(headerName, "-bin") {
				if v, err := decode(val); err == nil {
					val = v
				}
			}
			md[headerName] = append(md[headerName], val)
		}
	}
	return md
}

var base64Codecs = []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding}

func decode(val string) (string, error) {
	var firstErr error
	var b []byte
	// we are lenient and can accept any of the flavors of base64 encoding
	for _, d := range base64Codecs {
		var err error
		b, err = d.DecodeString(val)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return string(b), nil
	}
	return "", firstErr
}

/*
// Generate Google access token. Used when running in GCP
func (toExecutionWorkerObject *MessagesToExecutionWorkerObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if toExecutionWorkerObject.GcpAccessToken == nil || toExecutionWorkerObject.GcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+common_config.FenixExecutionWorkerAddress)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "0cf31da5-9e6b-41bc-96f1-6b78fb446194",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":    "8b1ca089-0797-4ee6-bf9d-f9b06f606ae9",
				"token": token,
			}).Debug("Got Bearer Token")
		}

		toExecutionWorkerObject.GcpAccessToken = token

	}

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "cd124ca3-87bb-431b-9e7f-e044c52b4960",
		"FenixExecutionWorkerObject.gcpAccessToken": toExecutionWorkerObject.GcpAccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+toExecutionWorkerObject.GcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}

*/
