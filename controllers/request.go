package node

import (
	"errors"
	"net/http"
	"strings"

	"github.com/zebek95/draflow-api/models"
)

// NodeRequest is the request payload for Node data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type NodeRequest struct {
	*models.Node

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *NodeRequest) Bind(r *http.Request) error {
	// a.Node is nil if no Node fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Node == nil {
		return errors.New("missing required Node fields.")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	a.ProtectedID = ""                         // unset the protected ID
	a.Node.Name = strings.ToLower(a.Node.Name) // as an example, we down-case
	return nil
}
