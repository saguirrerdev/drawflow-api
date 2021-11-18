package models

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/zebek95/draflow-api/database"
)

type Node struct {
	Uid   string   `json:"uid"`
	Name  string   `json:"name"`
	Data  string   `json:"data"`
	Nodes string   `json:"nodes"`
	DType []string `json:"dgraph.type"`
}

type NodeResponse struct {
	Node []Node `json:"node"`
}

func DbNewNode(node *Node) *Node {
	node.DType = append(node.DType, "Node")

	nodejson, err := json.Marshal(node)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   nodejson,
	}

	response, err := database.Conn.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		log.Fatal(err)
	}

	keys := reflect.ValueOf(response.Uids).MapKeys()
	uid := response.GetUids()
	key := keys[len(keys)-1]

	nod := DbGetNode(uid[key.String()])

	return nod
}

func DbGetNode(id string) *Node {
	params := map[string]string{"$id": id}

	q := `
		query Node($id: string) {
			node(func: uid($id)) {
				uid
				name
				data
				nodes
				dgraph.type
			}
		}
	`
	response, err := database.Conn.NewTxn().QueryWithVars(context.Background(), q, params)
	if err != nil {
		log.Fatal(err)
	}

	var obj NodeResponse

	err = json.Unmarshal(response.Json, &obj)
	if err != nil {
		log.Fatal(err)
	}

	return &obj.Node[0]
}

func DbGetNodes() *NodeResponse {
	const q = `
		{
			node(func: type(Node)){
				uid
				name
				data
				nodes
				dgraph.type
			}
		}
	`

	response, err := database.Conn.NewTxn().Query(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}

	var obj NodeResponse

	err = json.Unmarshal(response.Json, &obj)
	if err != nil {
		log.Fatal(err)
	}

	return &obj
}

func DbRemoveNode(id string) {
	params := map[string]string{"uid": id}

	pb, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: pb,
	}

	_, err = database.Conn.NewTxn().Mutate(context.Background(), mu)
	if err != nil {
		log.Fatal(err)
	}
}
