package gcp

import (
	"github.com/markbates/goth"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type GcpObjectStruct struct {
	logger                                    *logrus.Logger
	gcpAccessTokenForServiceAccounts          *oauth2.Token
	gcpAccessTokenForAuthorizedAccounts       goth.User
	gcpAccessTokenForServiceAccountsPubSub    *oauth2.Token
	gcpAccessTokenForAuthorizedAccountsPubSub goth.User
}

var Gcp GcpObjectStruct
