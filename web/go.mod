module github.com/double-nibble/algothon/web

go 1.13

require (
	github.com/double-nibble/algothon/userdb v0.0.0
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/sessions v1.2.0
	github.com/sirupsen/logrus v1.6.0
)

replace github.com/double-nibble/algothon/userdb => ../userdb
