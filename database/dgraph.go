package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

var Conn *dgo.Dgraph

func Init() {
	NewClient()
}

func NewClient() {
	if Conn == nil {
		fmt.Println("Creating new database connection")
		d, err := grpc.Dial(fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		Conn = dgo.NewDgraphClient(
			api.NewDgraphClient(d),
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
	// Set database schema
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			name: string .
			data: string .

			type Node {
				name
				data
			}
		`,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func reset(c *dgo.Dgraph) {
	fmt.Println("Restaring database")
	// Clean database
	err := c.Alter(context.Background(), &api.Operation{
		DropOp: api.Operation_ALL,
	})

	if err != nil {
		log.Fatal(err)
	}
}
