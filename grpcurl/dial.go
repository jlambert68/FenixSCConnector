package grpcurl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"net"
)

const no_version = "dev build <no version set>"

var version = no_version

/*
var (

	exit = os.Exit

	isUnixSocket func() bool // nil when run on non-unix platform

	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	help = flags.Bool("help", false, prettify(`
		Print usage instructions and exit.`))
	printVersion = flags.Bool("version", false, prettify(`
		Print version.`))
	plaintext = flags.Bool("plaintext", false, prettify(`
		Use plain-text HTTP/2 when connecting to server (no TLS).`))
	insecure = flags.Bool("insecure", false, prettify(`
		Skip server certificate and domain verification. (NOT SECURE!) Not
		valid with -plaintext option.`))
	cacert = flags.String("cacert", "", prettify(`
		File containing trusted root certificates for verifying the server.
		Ignored if -insecure is specified.`))
	cert = flags.String("cert", "", prettify(`
		File containing client certificate (public key), to present to the
		server. Not valid with -plaintext option. Must also provide -key option.`))
	key = flags.String("key", "", prettify(`
		File containing client private key, to present to the server. Not valid
		with -plaintext option. Must also provide -cert option.`))
	protoset      multiString
	protoFiles    multiString
	importPaths   multiString
	addlHeaders   multiString
	rpcHeaders    multiString
	reflHeaders   multiString
	expandHeaders = flags.Bool("expand-headers", false, prettify(`
		If set, headers may use '${NAME}' syntax to reference environment
		variables. These will be expanded to the actual environment variable
		value before sending to the server. For example, if there is an
		environment variable defined like FOO=bar, then a header of
		'key: ${FOO}' would expand to 'key: bar'. This applies to -H,
		-rpc-header, and -reflect-header options. No other expansion/escaping is
		performed. This can be used to supply credentials/secrets without having
		to put them in command-line arguments.`))
	authority = flags.String("authority", "", prettify(`
		The authoritative name of the remote server. This value is passed as the
		value of the ":authority" pseudo-header in the HTTP/2 protocol. When TLS
		is used, this will also be used as the server name when verifying the
		server's certificate. It defaults to the address that is provided in the
		positional arguments.`))
	userAgent = flags.String("user-agent", "", prettify(`
		If set, the specified value will be added to the User-Agent header set
		by the grpc-go library.
		`))
	data = flags.String("d", "", prettify(`
		Data for request contents. If the value is '@' then the request contents
		are read from stdin. For calls that accept a stream of requests, the
		contents should include all such request messages concatenated together
		(possibly delimited; see -format).`))
	format = flags.String("format", "json", prettify(`
		The format of request data. The allowed values are 'json' or 'text'. For
		'json', the input data must be in JSON format. Multiple request values
		may be concatenated (messages with a JSON representation other than
		object must be separated by whitespace, such as a newline). For 'text',
		the input data must be in the protobuf text format, in which case
		multiple request values must be separated by the "record separator"
		ASCII character: 0x1E. The stream should not end in a record separator.
		If it does, it will be interpreted as a final, blank message after the
		separator.`))
	allowUnknownFields = flags.Bool("allow-unknown-fields", false, prettify(`
		When true, the request contents, if 'json' format is used, allows
		unknown fields to be present. They will be ignored when parsing
		the request.`))
	connectTimeout = flags.Float64("connect-timeout", 0, prettify(`
		The maximum time, in seconds, to wait for connection to be established.
		Defaults to 10 seconds.`))
	formatError = flags.Bool("format-error", false, prettify(`
		When a non-zero status is returned, format the response using the
		value set by the -format flag .`))
	keepaliveTime = flags.Float64("keepalive-time", 0, prettify(`
		If present, the maximum idle time in seconds, after which a keepalive
		probe is sent. If the connection remains idle and no keepalive response
		is received for this same period then the connection is closed and the
		operation fails.`))
	maxTime = flags.Float64("max-time", 0, prettify(`
		The maximum total time the operation can take, in seconds. This is
		useful for preventing batch jobs that use grpcurl from hanging due to
		slow or bad network links or due to incorrect stream method usage.`))
	maxMsgSz = flags.Int("max-msg-sz", 0, prettify(`
		The maximum encoded size of a response message, in bytes, that grpcurl
		will accept. If not specified, defaults to 4,194,304 (4 megabytes).`))
	emitDefaults = flags.Bool("emit-defaults", false, prettify(`
		Emit default values for JSON-encoded responses.`))
	protosetOut = flags.String("protoset-out", "", prettify(`
		The name of a file to be written that will contain a FileDescriptorSet
		proto. With the list and describe verbs, the listed or described
		elements and their transitive dependencies will be written to the named
		file if this option is given. When invoking an RPC and this option is
		given, the method being invoked and its transitive dependencies will be
		included in the output file.`))
	msgTemplate = flags.Bool("msg-template", false, prettify(`
		When describing messages, show a template of input data.`))
	verbose = flags.Bool("v", false, prettify(`
		Enable verbose output.`))
	veryVerbose = flags.Bool("vv", false, prettify(`
		Enable very verbose output.`))
	serverName = flags.String("servername", "", prettify(`
		Override server name when validating TLS certificate. This flag is
		ignored if -plaintext or -insecure is used.
		NOTE: Prefer -authority. This flag may be removed in the future. It is
		an error to use both -authority and -servername (though this will be
		permitted if they are both set to the same value, to increase backwards
		compatibility with earlier releases that allowed both to be set).`))
	reflection = optionalBoolFlag{val: true}

)
*/
type optionalBoolFlag struct {
	set, val bool
}

/*
func dial() *grpc.ClientConn {
	dialTime := 10 * time.Second
	if *connectTimeout > 0 {
		dialTime = time.Duration(*connectTimeout * float64(time.Second))
	}
	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var opts []grpc.DialOption
	if *keepaliveTime > 0 {
		timeout := time.Duration(*keepaliveTime * float64(time.Second))
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    timeout,
			Timeout: timeout,
		}))
	}
	if *maxMsgSz > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(*maxMsgSz)))
	}
	var creds credentials.TransportCredentials
	if !*plaintext {
		tlsConf, err := grpcurl.ClientTLSConfig(*insecure, *cacert, *cert, *key)
		if err != nil {
			fail(err, "Failed to create TLS config")
		}

		sslKeylogFile := os.Getenv("SSLKEYLOGFILE")
		if sslKeylogFile != "" {
			w, err := os.OpenFile(sslKeylogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
			if err != nil {
				fail(err, "Could not open SSLKEYLOGFILE %s", sslKeylogFile)
			}
			tlsConf.KeyLogWriter = w
		}

		creds = credentials.NewTLS(tlsConf)

		// can use either -servername or -authority; but not both
		if *serverName != "" && *authority != "" {
			if *serverName == *authority {
				warn("Both -servername and -authority are present; prefer only -authority.")
			} else {
				fail(nil, "Cannot specify different values for -servername and -authority.")
			}
		}
		overrideName := *serverName
		if overrideName == "" {
			overrideName = *authority
		}

		if overrideName != "" {
			opts = append(opts, grpc.WithAuthority(overrideName))
		}
	} else if *authority != "" {
		opts = append(opts, grpc.WithAuthority(*authority))
	}

	grpcurlUA := "grpcurl/" + version
	if version == no_version {
		grpcurlUA = "grpcurl/dev-build (no version set)"
	}
	if *userAgent != "" {
		grpcurlUA = *userAgent + " " + grpcurlUA
	}
	opts = append(opts, grpc.WithUserAgent(grpcurlUA))

	network := "tcp"
	if isUnixSocket != nil && isUnixSocket() {
		network = "unix"
	}
	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		fail(err, "Failed to dial target host %q", target)
	}
	return cc
}

*/

// ClientTransportCredentials is a helper function that constructs a TLS config with
// the given properties (see ClientTLSConfig) and then constructs and returns gRPC
// transport credentials using that config.
//
// Deprecated: Use grpcurl.ClientTLSConfig and credentials.NewTLS instead.
func ClientTransportCredentials(insecureSkipVerify bool, cacertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	tlsConf, err := ClientTLSConfig(insecureSkipVerify, cacertFile, clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(tlsConf), nil
}

// ClientTLSConfig builds transport-layer config for a gRPC client using the
// given properties. If cacertFile is blank, only standard trusted certs are used to
// verify the server certs. If clientCertFile is blank, the client will not use a client
// certificate. If clientCertFile is not blank then clientKeyFile must not be blank.
func ClientTLSConfig(insecureSkipVerify bool, cacertFile, clientCertFile, clientKeyFile string) (*tls.Config, error) {
	var tlsConf tls.Config

	if clientCertFile != "" {
		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %v", err)
		}
		tlsConf.Certificates = []tls.Certificate{certificate}
	}

	if insecureSkipVerify {
		tlsConf.InsecureSkipVerify = true
	} else if cacertFile != "" {
		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(cacertFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %v", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca certs")
		}

		tlsConf.RootCAs = certPool
	}

	return &tlsConf, nil
}

// ServerTransportCredentials builds transport credentials for a gRPC server using the
// given properties. If cacertFile is blank, the server will not request client certs
// unless requireClientCerts is true. When requireClientCerts is false and cacertFile is
// not blank, the server will verify client certs when presented, but will not require
// client certs. The serverCertFile and serverKeyFile must both not be blank.
func ServerTransportCredentials(cacertFile, serverCertFile, serverKeyFile string, requireClientCerts bool) (credentials.TransportCredentials, error) {
	var tlsConf tls.Config
	// TODO(jh): Remove this line once https://github.com/golang/go/issues/28779 is fixed
	// in Go tip. Until then, the recently merged TLS 1.3 support breaks the TLS tests.
	tlsConf.MaxVersion = tls.VersionTLS12

	// Load the server certificates from disk
	certificate, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load key pair: %v", err)
	}
	tlsConf.Certificates = []tls.Certificate{certificate}

	if cacertFile != "" {
		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(cacertFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %v", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("failed to append ca certs")
		}

		tlsConf.ClientCAs = certPool
	}

	if requireClientCerts {
		tlsConf.ClientAuth = tls.RequireAndVerifyClientCert
	} else if cacertFile != "" {
		tlsConf.ClientAuth = tls.VerifyClientCertIfGiven
	} else {
		tlsConf.ClientAuth = tls.NoClientCert
	}

	return credentials.NewTLS(&tlsConf), nil
}

// BlockingDial is a helper method to dial the given address, using optional TLS credentials,
// and blocking until the returned connection is ready. If the given credentials are nil, the
// connection will be insecure (plain-text).
func BlockingDial(ctx context.Context, network, address string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// grpc.Dial doesn't provide any information on permanent connection errors (like
	// TLS handshake failures). So in order to provide good error messages, we need a
	// custom dialer that can provide that info. That means we manage the TLS handshake.
	result := make(chan interface{}, 1)

	writeResult := func(res interface{}) {
		// non-blocking write: we only need the first result
		select {
		case result <- res:
		default:
		}
	}

	// custom credentials and dialer will notify on error via the
	// writeResult function
	if creds != nil {
		creds = &errSignalingCreds{
			TransportCredentials: creds,
			writeResult:          writeResult,
		}
	}
	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		// NB: We *could* handle the TLS handshake ourselves, in the custom
		// dialer (instead of customizing both the dialer and the credentials).
		// But that requires using insecure.NewCredentials() dial transport
		// option (so that the gRPC library doesn't *also* try to do a
		// handshake). And that would mean that the library would send the
		// wrong ":scheme" metaheader to servers: it would send "http" instead
		// of "https" because it is unaware that TLS is actually in use.
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err != nil {
			writeResult(err)
		}
		return conn, err
	}

	// Even with grpc.FailOnNonTempDialError, this call will usually timeout in
	// the face of TLS handshake errors. So we can't rely on grpc.WithBlock() to
	// know when we're done. So we run it in a goroutine and then use result
	// channel to either get the connection or fail-fast.
	go func() {
		// We put grpc.FailOnNonTempDialError *before* the explicitly provided
		// options so that it could be overridden.
		opts = append([]grpc.DialOption{grpc.FailOnNonTempDialError(true)}, opts...)
		// But we don't want caller to be able to override these two, so we put
		// them *after* the explicitly provided options.
		opts = append(opts, grpc.WithBlock(), grpc.WithContextDialer(dialer))

		if creds == nil {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
		conn, err := grpc.DialContext(ctx, address, opts...)
		var res interface{}
		if err != nil {
			res = err
		} else {
			res = conn
		}
		writeResult(res)
	}()

	select {
	case res := <-result:
		if conn, ok := res.(*grpc.ClientConn); ok {
			return conn, nil
		}
		return nil, res.(error)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// errSignalingCreds is a wrapper around a TransportCredentials value, but
// it will use the writeResult function to notify on error.
type errSignalingCreds struct {
	credentials.TransportCredentials
	writeResult func(res interface{})
}

func (c *errSignalingCreds) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	conn, auth, err := c.TransportCredentials.ClientHandshake(ctx, addr, rawConn)
	if err != nil {
		c.writeResult(err)
	}
	return conn, auth, err
}
