package main

import (
	"net/http"

	server "carrier/Explorer/domestic/explorerserver"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := httprouter.New()
	router.POST("/addServer", server.AddServer)
	log.Info("starting REST server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
