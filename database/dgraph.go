package database

import (
	"context"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

var Conn *dgo.Dgraph

func Init() {
	NewClient()
}

func NewClient() {
	if Conn == nil {
		fmt.Println("Creating new database connection")
		// d, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")), grpc.WithInsecure())
		// if err != nil {
		// 	log.Fatal(err)
		// }

		conn, err := dgo.DialCloud("https://blue-surf-460029.us-east-1.aws.cloud.dgraph.io/graphql", "NzcxMTkwM2ZkYzhmMzFmNDliZTY0MWE4ODdjMDVlYmQ=")
		if err != nil {
			log.Fatal(err)
		}

		// defer conn.Close()

		Conn = dgo.NewDgraphClient(
			api.NewDgraphClient(conn),
		)

		// if os.Getenv("ENVIRONMENT") == "development" {
		// 	reset(Conn)
		// }

		SetSchema(Conn)
	} else {
		fmt.Printf("Dabase instance already defined")
	}
}

func SetSchema(c *dgo.Dgraph) {
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			name: string .
			data: string .
			nodes: string .

			type Node {
				name
				data
				nodes
			}
		`,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func reset(c *dgo.Dgraph) {
	fmt.Println("Restarting database")
	err := c.Alter(context.Background(), &api.Operation{
		DropOp: api.Operation_ALL,
	})

	if err != nil {
		log.Fatal(err)
	}
}
