package node

import (
	"errors"
	"net/http"
	"strings"

	"github.com/zebek95/draflow-api/models"
)

type NodeRequest struct {
	*models.Node

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *NodeRequest) Bind(r *http.Request) error {
	if a.Node == nil {
		return errors.New("missing required Node fields.")
	}

	// just a post-process after a decode..
	a.ProtectedID = ""
	a.Node.Name = strings.ToLower(a.Node.Name)
	return nil
}
