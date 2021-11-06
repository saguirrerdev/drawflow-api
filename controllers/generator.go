package node

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/zebek95/draflow-api/error_handler"
	"github.com/zebek95/draflow-api/models"
)

type Df struct {
	Drawflow []DfModule `json:"drawflow"`
}

type DfModule struct {
	Home []DfNode `json:"Home"`
}

type DfNode struct {
	Id       int         `json:"id"`
	Name     string      `json:"name"`
	Data     DfData      `json:"data"`
	Class    string      `json:"class"`
	Typenode string      `json:"typenode"`
	Html     string      `json:"html"`
	PosX     int         `json:"pos_x"`
	PosY     int         `json:"pos_y"`
	Inputs   []DfInputs  `json:"inputs"`
	Outputs  []DfOutputs `json:"outputs"`
}

type DfData struct {
	Value string `json:"value"`
}

type DfOutputs struct {
	Output1 []DfConnectionOutput `json:"output_1"`
	Output2 []DfConnectionOutput `json:"output_2"`
	Output3 []DfConnectionOutput `json:"output_3"`
}

type DfConnectionOutput struct {
	Node   string `json:"node"`
	Output string `json:"output"`
}

type DfInputs struct {
	Input1 []DfConnectionInput `json:"input_1"`
	Input2 []DfConnectionInput `json:"input_2"`
	Input3 []DfConnectionInput `json:"input_3"`
}

type DfConnectionInput struct {
	Node  string `json:"node"`
	Input string `json:"input"`
}

func GetCode(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the node
	// context because this handler is a child of the NodeCtx
	// middleware. The worst case, the recoverer middleware will save us.
	node := r.Context().Value("node").(*models.Node)

	var data Df

	json.Unmarshal([]byte(node.Data), &data)

	fmt.Println(node.Data)

	if err := render.Render(w, r, NewNodeResponse(node)); err != nil {
		render.Render(w, r, error_handler.ErrRender(err))
		return
	}
}

func generate(data map[string]interface{}) string {
	fmt.Println(data)
	module := data["drawflow"]
	fmt.Println(module)
	return ""
}
