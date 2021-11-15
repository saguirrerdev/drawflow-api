package node

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/zebek95/draflow-api/error_handler"
	"github.com/zebek95/draflow-api/models"
)

func ListNodes(w http.ResponseWriter, r *http.Request) {
	nodes := models.DbGetNodes()

	if err := render.RenderList(w, r, NewNodeListResponse(nodes.Node)); err != nil {
		render.Render(w, r, error_handler.ErrRender(err))
		return
	}
}

func CreateNode(w http.ResponseWriter, r *http.Request) {
	data := &NodeRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, error_handler.ErrInvalidRequest(err))
		return
	}

	node := models.DbNewNode(data.Node)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewNodeResponse(node))
}

func GetNode(w http.ResponseWriter, r *http.Request) {
	node := r.Context().Value("node").(*models.Node)

	if err := render.Render(w, r, NewNodeResponse(node)); err != nil {
		render.Render(w, r, error_handler.ErrRender(err))
		return
	}
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {

	node := r.Context().Value("node").(*models.Node)

	models.DbRemoveNode(node.Uid)

	render.Render(w, r, DeletedNodeResponse())
}

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
