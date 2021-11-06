package node

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/zebek95/draflow-api/error_handler"
	"github.com/zebek95/draflow-api/models"
)

func ListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := models.DbGetNodes()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(nodes)

	if err := render.RenderList(w, r, NewNodeListResponse(nodes.Node)); err != nil {
		render.Render(w, r, error_handler.ErrRender(err))
		return
	}
}

// NodeCtx middleware is used to load an Node object from
// the URL parameters passed through as the request. In case
// the Node could not be found, we stop here and return a 404.
func NodeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var node *models.Node
		var err error

		if nodeID := chi.URLParam(r, "nodeID"); nodeID != "" {
			node, err = models.DbGetNode(nodeID)
		} else {
			render.Render(w, r, error_handler.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, error_handler.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "node", node)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateNode persists the posted Node and returns it
// back to the client as an acknowledgement.
func CreateNode(w http.ResponseWriter, r *http.Request) {
	data := &NodeRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, error_handler.ErrInvalidRequest(err))
		return
	}

	node, err := models.DbNewNode(data.Node)
	if err != nil {
		log.Fatal(err)
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewNodeResponse(node))
}

// GetNode returns the specific Node. You'll notice it just
// fetches the Node right off the context, as its understood that
// if we made it this far, the Node must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetNode(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the node
	// context because this handler is a child of the NodeCtx
	// middleware. The worst case, the recoverer middleware will save us.
	node := r.Context().Value("node").(*models.Node)

	if err := render.Render(w, r, NewNodeResponse(node)); err != nil {
		render.Render(w, r, error_handler.ErrRender(err))
		return
	}
}

// UpdateNode updates an existing Node in our persistent store.
func UpdateNode(w http.ResponseWriter, r *http.Request) {
	node := r.Context().Value("node").(*models.Node)

	data := &NodeRequest{Node: node}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, error_handler.ErrInvalidRequest(err))
		return
	}
	node = data.Node
	// models.DbUpdateNode(node.ID, node)

	render.Render(w, r, NewNodeResponse(node))
}

// DeleteNode removes an existing Node from our persistent store.
func DeleteNode(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the node
	// context because this handler is a child of the NodeCtx
	// middleware. The worst case, the recoverer middleware will save us.
	node := r.Context().Value("node").(*models.Node)

	// node, err = models.DbRemoveNode(node.ID)
	if err != nil {
		render.Render(w, r, error_handler.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewNodeResponse(node))
}

// This is entirely optional, but I wanted to demonstrate how you could easily
// add your own logic to the render.Respond method.
func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

// NOTE: as a thought, the request and response payloads for an Node could be the
// same payload type, perhaps will do an example with it as well.
// type NodePayload struct {
//   *Node
// }
