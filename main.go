package main

import (
	_ "github.com/kyawmyintthein/golangRestfulAPISample/docs"
	"flag"
	"github.com/kyawmyintthein/golangRestfulAPISample/app"
)

// @title App Name
// @version 1.0
// @description  App name documentation.

// @contact.name API Support
// @contact.email email address

// @host localhost:3030
// @BasePath /
func main(){
	var configFilePath string
	var serverPort string
	flag.StringVar(&configFilePath, "config", "config.yml", "absolute path to the configuration file")
	flag.StringVar(&serverPort, "server_port", "4000", "port on which server runs")
	flag.Parse()

	application := app.New(configFilePath)

	// init necessary module before start
	application.Init()

	// start http server
	application.Start(serverPort)
}
