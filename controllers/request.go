package node

import (
	"errors"
	"net/http"
	"strings"

	"github.com/zebek95/draflow-api/models"
)

type NodeRequest struct {
	*models.Node
}

func (a *NodeRequest) Bind(r *http.Request) error {
	if a.Node == nil {
		return errors.New("missing required Node fields.")
	}

	a.Node.Name = strings.ToLower(a.Node.Name)
	return nil
}
