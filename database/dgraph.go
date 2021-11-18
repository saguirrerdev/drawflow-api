package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
)

var Conn *dgo.Dgraph

func Init() {
	NewClient()
}

func NewClient() {
	if Conn == nil {
		conn, err := dgo.DialCloud(os.Getenv("DB_HOST"), os.Getenv("DB_ACCESS_KEY"))
		if err != nil {
			log.Fatal(err)
		}

		Conn = dgo.NewDgraphClient(
			api.NewDgraphClient(conn),
		)

		SetSchema(Conn)
	} else {
		fmt.Printf("Database instance already defined")
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
