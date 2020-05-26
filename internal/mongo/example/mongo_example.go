package main

import (
	"bitbucket.org/libertywireless/golang_artifacts/go-database/clmongo"
	"context"
	"fmt"
)

func main(){
	mongoCfg := clmongo.MongodbConfig{
		DatabaseName:  "test",
		DatabaseHosts: "localhost:27017",
		TimeOut:       10,
		DialTimeOut:   10,
		PoolSize:      10,
		Username:      "",
		Password:      "",
	}

	mongoConnector, err := clmongo.NewMongodbConnector(&mongoCfg)
	if err != nil{
		panic(err)
	}

	db := mongoConnector.DB(context.Background())
	collection := db.Collection("users")
	var result interface{}
	err = collection.FindOne(context.Background(), map[string]interface{}{"id": "1"}).Decode(&result)
	if result != nil{
		panic(result)
	}
	fmt.Printf("Successfully connected to %s \n", db.Name())
}

