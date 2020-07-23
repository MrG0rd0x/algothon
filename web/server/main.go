package server

import (
	"context"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/double-nibble/algothon/userdb"
	"github.com/double-nibble/algothon/web/routes"
)

var log = logrus.WithFields(logrus.Fields{"service": "web"})

// ListenAndServe starts a new server that should only return on fatal errors
func ListenAndServe() {
	log.Info("starting TLS listener for registration server")
	db := userdb.NewConnection("redis:6379", os.Getenv("REDIS_PASSWORD"))
	if !db.Verify(context.Background()) {
		log.Fatal("Could not create connection to userDB")
	}
	log.Debug("DB connection successful")
	r := routes.NewRouter(&routes.Config{
		DB:         db,
		SessionKey: []byte(os.Getenv("SESSION_KEY")),
		Salt:       os.Getenv("SALT"),
	})
	http.ListenAndServe(":8080", r)
}
