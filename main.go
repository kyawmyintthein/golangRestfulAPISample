package main

import (
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
	flag.StringVar(&configFilePath, "config", "config.yml", "absolute path to the configuration file")
	flag.Parse()

	application, err := app.NewApp(configFilePath)
	if err != nil{
		panic(err)
	}

	// init necessary module before start
	err = application.Init()
	if err != nil{
		panic(err)
	}

	// start http server
	application.StartHttpServer()
}