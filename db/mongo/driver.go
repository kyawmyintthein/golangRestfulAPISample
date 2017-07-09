package mongo

import "golangRestfulAPISample/bootstrap"

var maxPool int

// init mongodb 
func Init() {
	var adapter string
	adapter = bootstrap.App.DBConfig.String("adapter")
	if adapter == "mongodb" {
		maxPool = bootstrap.App.DBConfig.Int("mongodb.max_pool")
		checkAndInitServiceConnection()
	}
}

// checkAndInitServiceConnection
func checkAndInitServiceConnection() {
	var err error
	if service.baseSession == nil {
		service.URL = bootstrap.App.DBConfig.String("mongodb.path")
		if err = service.New(); err != nil {
			panic(err)
		}
	}
}
