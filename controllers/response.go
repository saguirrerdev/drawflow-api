package node

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/zebek95/draflow-api/models"
)

// NodeResponse is the response payload for the Node data model.
// See NOTE above in NodeRequest as well.
//
// In the NodeResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type NodeResponse struct {
	models.Node

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func NewNodeResponse(node *models.Node) *NodeResponse {
	resp := &NodeResponse{Node: *node}
	return resp
}

func (rd *NodeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func NewNodeListResponse(nodes []models.Node) []render.Renderer {

	list := []render.Renderer{}
	for _, node := range nodes {
		list = append(list, NewNodeResponse(&node))
	}
	return list
}
