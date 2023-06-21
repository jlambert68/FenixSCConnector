package gcp

import (
	"FenixSCConnector/common_config"
	"FenixSCConnector/notToBeSentToGithub"
	"fmt"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sirupsen/logrus"
	"github.com/toqueteos/webbrowser"
	"golang.org/x/net/context"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"html/template"
	"net/http"
	"time"
)

func (gcp *GcpObjectStruct) GenerateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Chose correct method for authentication
	if common_config.UseServiceAccount == true {
		// Use Service account
		appendedCtx, returnAckNack, returnMessage = gcp.GenerateGCPAccessTokenForServiceAccount(ctx)
	} else {
		// User log into GCP via web
		appendedCtx, returnAckNack, returnMessage = gcp.GenerateGCPAccessTokenForAuthorizedUser(ctx)
	}

	return appendedCtx, returnAckNack, returnMessage

}

// GenerateGCPAccessTokenForServiceAccount Generate Google access token for a service account. Used when running in GCP
func (gcp *GcpObjectStruct) GenerateGCPAccessTokenForServiceAccount(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// If there shouldn't be any call to Worker then just exit
	if common_config.TurnOffCallToWorker == true {
		common_config.Logger.WithFields(logrus.Fields{
			"id": "dbc95a47-3b66-4b14-a364-6df8017662b5",
		}).Debug("Execution Worker shouldn't be called, exit 'GenerateGCPAccessTokenForServiceAccount'")

		return
	}

	// Only create the token if there is none, or it has expired
	if gcp.gcpAccessTokenForServiceAccounts == nil || gcp.gcpAccessTokenForServiceAccounts.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		/*
			tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+common_config.FenixGuiBuilderServerAddress)
			if err != nil {
				fenixGuiBuilderProxyServerObject.logger.WithFields(logrus.Fields{
					"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
					"err": err,
				}).Error("Couldn't generate access token")

				return nil, false, "Couldn't generate access token"
			}

			token, err := tokenSource.Token()
		*/
		/*
			var eMailAndPrivateKey = struct {
				Email      string `json:"client_email"`
				PrivateKey string `json:"private_key"`
			}{}
			json.Unmarshal(serviceAccountKeyJson, &eMailAndPrivateKey)
			config := &jwt.Config{
				Email:      eMailAndPrivateKey.Email,
				PrivateKey: []byte(eMailAndPrivateKey.PrivateKey),
				Scopes: []string{
					gcp_scope,
				},
				TokenURL:   google.JWTTokenURL,
				UseIDToken: false,
			}

		*/

		tokenSource, err := idtoken.NewTokenSource(ctx, notToBeSentToGithub.Gcp_scope, idtoken.WithCredentialsJSON(notToBeSentToGithub.ServiceAccountKeyJson))

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "c11aba6d-a177-4b88-a521-15bec26a3312",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		//token, err := config.TokenSource(oauth2.NoContext).Token()

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "0e0dd854-2088-419d-a0d1-035b6242c585",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":    "1ee9fd6d-c83d-4dbb-bce7-fd4901ff3f87",
				"token": "Nothing to see", //token,
			}).Debug("Got Bearer Token")
		}

		gcp.gcpAccessTokenForServiceAccounts = token

	}

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "fa86e0bd-9c47-4f1e-b9dd-2b1f08880b2b",
		"fenixGuiBuilderProxyServerObject.gcpAccessTokenForServiceAccounts": "Nothing to see", //fenixGuiBuilderProxyServerObject.gcpAccessTokenForServiceAccounts,
	}).Debug("Will use Bearer Token")

	// Add token to gRPC Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.gcpAccessTokenForServiceAccounts.AccessToken)

	return appendedCtx, true, ""

}

// DoneChannel - channel used for to close down local web server
var DoneChannel chan bool

func (gcp *GcpObjectStruct) GenerateGCPAccessTokenForAuthorizedUser(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Secure that User is initiated
	gcp.initiateUserObject()

	// Only create the token if there is none, or it has expired (or 5 minutes before expiration
	timeToCompareTo := time.Now().Add(-time.Minute * 5)
	if !(gcp.gcpAccessTokenForAuthorizedAccounts.IDToken == "" || gcp.gcpAccessTokenForAuthorizedAccounts.ExpiresAt.Before(timeToCompareTo)) {
		// We already have a ID-token that can be used, so return that
		appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.gcpAccessTokenForAuthorizedAccounts.IDToken)

		return appendedCtx, true, ""
	}

	// Need to create a new ID-token

	key := common_config.ApplicationRunTimeUuid // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30                        // 30 days
	isProd := false                             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		// Use 'Fenix End User Authentication'
		google.New(
			common_config.AuthClientId,
			common_config.AuthClientSecret,
			"http://localhost:3000/auth/google/callback",
			"email", "profile"),
	)

	router := pat.New()

	router.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {

		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {

			fmt.Fprintln(res, err)

			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)

		// Save ID-token
		gcp.gcpAccessTokenForAuthorizedAccounts = user

		// Trigger Close of Web Server, and 'true' means that a ID-to
		DoneChannel <- true

	})

	router.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})

	// Initiate channel used to stop server
	DoneChannel = make(chan bool, 1)

	// Initiate http server
	localWebServer := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	// Start Local Web Server as go routine
	go gcp.startLocalWebServer(localWebServer)

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "689d42de-3cc0-4237-b1e9-3a6c769f65ea",
	}).Debug("Local webServer Started")

	// Wait for message in channel to stop local web server
	gotIdTokenResult := <-DoneChannel

	// Shutdown local web server
	gcp.stopLocalWebServer(context.Background(), localWebServer)

	// Depending on the outcome of getting a token return different results
	if gotIdTokenResult == true {
		// Success in getting an ID-token
		appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+gcp.gcpAccessTokenForAuthorizedAccounts.IDToken)

		return appendedCtx, true, ""
	} else {
		// Didn't get any ID-token
		return nil, false, "Couldn't generate access token"
	}

}

// Start and run Local Web Server
func (gcp *GcpObjectStruct) startLocalWebServer(webServer *http.Server) {

	go func() {
		time.Sleep(1 * time.Second)
		err := webbrowser.Open("http://localhost:3000")

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "17bc0305-4594-48e1-bb8d-c642579e5e56",
				"err": err,
			}).Fatalf("Couldn't open the web browser")
		}
	}()
	err := webServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		common_config.Logger.WithFields(logrus.Fields{
			"ID": "8226cf74-0cdc-4e29-a441-116504b4b333",
		}).Fatalf("Local Web Server failed to listen: %s\n", err)

	}
}

// Close down Local Web Server
func (gcp *GcpObjectStruct) stopLocalWebServer(ctx context.Context, webServer *http.Server) {

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "1f4e0354-2a09-4a1d-be61-67ecda781142",
	}).Debug("Trying to stop local web server")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := webServer.Shutdown(ctx)
	if err != nil {
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID": "ea06dfab-39b9-4df6-b3ca-7f5f56b3cb91",
			}).Fatalf("Local Web Server Shutdown Failed:%+v", err)

		} else {
			common_config.Logger.WithFields(logrus.Fields{
				"ID": "ea06dfab-39b9-4df6-b3ca-7f5f56b3cb91",
			}).Debug("Local Web Server Exited Properly")
		}

	}

}

// SetLogger
// Set to use the same Logger reference as is used by central part of system
func (gcp *GcpObjectStruct) SetLogger(logger *logrus.Logger) {

	//grpcOutVaraible = GRPCOutStruct{}

	//gcp.logger = logger

	return

}

// initiateUserObject
// Set to use the same Logger reference as is used by central part of system
func (gcp *GcpObjectStruct) initiateUserObject() {

	// Only do initiation if it's not done before

	if gcp.gcpAccessTokenForAuthorizedAccounts.UserID == "" {
		gcp.gcpAccessTokenForAuthorizedAccounts = goth.User{}
	}

	return

}
