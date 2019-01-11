package main

import (
	"github.com/op/go-logging"
	"os"
	"launchomega.com/MYTH/server"
)

var log = logging.MustGetLogger("")

func init() {
	format := logging.MustStringFormatter(`[%{level}] %{time} %{shortfunc}: %{message}`)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}

func main() {
	log.Info("Starting up MYTH...")
	server.StartServer()
}
