package node

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/zebek95/draflow-api/models"
)

type NodeResponse struct {
	models.Node

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

type DeleteNodeResponse struct {
	Message string `json:"message"`
}

func NewNodeResponse(node *models.Node) *NodeResponse {
	resp := &NodeResponse{Node: *node}
	return resp
}

func (rd *NodeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func DeletedNodeResponse() *DeleteNodeResponse {
	resp := &DeleteNodeResponse{
		Message: "Node has been deleted",
	}
	return resp
}

func (rd *DeleteNodeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewNodeListResponse(nodes []models.Node) []render.Renderer {

	list := []render.Renderer{}
	for _, node := range nodes {
		list = append(list, NewNodeResponse(&node))
	}
	return list
}
