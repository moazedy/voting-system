package repository

import (
	"fmt"
	"time"
	"voting-system/config"

	"github.com/couchbase/gocb/v2"
)

// DBS struct for managing database connections
type dbs struct {
	Couch *gocb.Cluster
}

var DBS dbs

// Init function for initialize databases connection
func Init() {
	couchBaseConnection()
}

// couchBaseConnection function for connection to couchbase server
func couchBaseConnection() {
	cluster, err := gocb.Connect(
		config.Configs.Couchbase.Addresses,
		gocb.ClusterOptions{
			Username: config.Configs.Couchbase.Username,
			Password: config.Configs.Couchbase.Password,
		})
	if err != nil {
		panic(err)
	}
	// Waiting for the cluster to be ready
	err = cluster.WaitUntilReady(time.Second*10, nil)
	if err != nil {
		panic(err)
	}
	// Pinging cluster to ensure query service working
	pingResult, err := cluster.Ping(&gocb.PingOptions{
		ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeQuery},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("couchbase connection ... ")
	fmt.Printf("%#v", pingResult.Services)
	DBS.Couch = cluster
}
