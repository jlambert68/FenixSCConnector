package gcp

import (
	"github.com/markbates/goth"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GcpObjectStruct struct {
	logger                                 *logrus.Logger
	gcpAccessTokenForServiceAccounts       *oauth2.Token
	gcpAccessTokenForAuthorizedAccounts    goth.User
	gcpAccessTokenForServiceAccountsPubSub *oauth2.Token

	gcpAccessTokenForAuthorizedAccountsPubSub goth.User

	refreshTokenResponse *RefreshTokenResponse

	// The following token is received from Worker, needs to be this due to the setup at SEB
	GcpAccessTokenFromWorkerToBeUsedWithPubSub string
}

var Gcp GcpObjectStruct
