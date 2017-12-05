package Cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

// Session Handle
var Session *gocql.Session

// This function will actually run before ../main.go's main function
func init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "lumberyard"
	cluster.ProtoVersion = 4
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}
