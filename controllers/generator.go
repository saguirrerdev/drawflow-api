package node

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/zebek95/draflow-api/models"
)

type Node struct {
	ID       int     `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Data     Data    `json:"data,omitempty"`
	Class    string  `json:"class,omitempty"`
	HTML     string  `json:"html,omitempty"`
	Typenode string  `json:"typenode,omitempty"`
	Inputs   Inputs  `json:"inputs,omitempty"`
	Outputs  Outputs `json:"outputs,omitempty"`
	PosX     int     `json:"pos_x,omitempty"`
	PosY     int     `json:"pos_y,omitempty"`
}

type Nodes []struct {
	ID       int     `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Data     Data    `json:"data,omitempty"`
	Class    string  `json:"class,omitempty"`
	HTML     string  `json:"html,omitempty"`
	Typenode string  `json:"typenode,omitempty"`
	Inputs   Inputs  `json:"inputs,omitempty"`
	Outputs  Outputs `json:"outputs,omitempty"`
	PosX     int     `json:"pos_x,omitempty"`
	PosY     int     `json:"pos_y,omitempty"`
}
type Data struct {
	Value     string `json:"value,omitempty"`
	Condition string `json:"condition,omitempty"`
	Then      string `json:"then,omitempty"`
	Else      string `json:"else,omitempty"`
}
type ConnectionsInput struct {
	Node  string `json:"node,omitempty"`
	Input string `json:"input,omitempty"`
}
type Input1 struct {
	Connections []ConnectionsInput `json:"connections,omitempty"`
}
type Input2 struct {
	Connections []ConnectionsInput `json:"connections,omitempty"`
}
type Input3 struct {
	Connections []ConnectionsInput `json:"connections,omitempty"`
}
type Inputs struct {
	Input1 Input1 `json:"input_1,omitempty"`
	Input2 Input2 `json:"input_2,omitempty"`
	Input3 Input3 `json:"input_3,omitempty"`
}
type ConnectionsOutput struct {
	Node   string `json:"node,omitempty"`
	Output string `json:"output,omitempty"`
}
type Output1 struct {
	Connections []ConnectionsOutput `json:"connections,omitempty"`
}
type Output2 struct {
	Connections []ConnectionsOutput `json:"connections,omitempty"`
}
type Output3 struct {
	Connections []ConnectionsOutput `json:"connections,omitempty"`
}
type Outputs struct {
	Output1 Output1 `json:"output_1,omitempty"`
	Output2 Output2 `json:"output_2,omitempty"`
	Output3 Output3 `json:"output_3,omitempty"`
}

type Code struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var nodes Nodes

func GetCode(w http.ResponseWriter, r *http.Request) {
	node := r.Context().Value("node").(*models.Node)

	json.Unmarshal([]byte(node.Nodes), &nodes)

	code := Code{Code: generate("", 0), Name: node.Name}

	render.Render(w, r, CodeResponse(code))
}

func CodeResponse(code Code) *Code {
	resp := code
	return &resp
}

func (rd Code) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func generate(code string, pos int) string {

	if pos > len(nodes)-1 {
		return code
	}

	el := nodes[pos]
	text := getNodeCode(el)

	code += text + "\n"

	return generate(code, pos+1)
}

func addGenerator(node Node) string {
	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("add%v = %s + %s", node.ID, nodeVarName(id1), nodeVarName(id2))
}

func getNodeCode(node Node) string {
	text := ""

	switch nodeType := node.HTML; nodeType {
	case "Df_add":
		text = addGenerator(node)
	case "Df_number":
		text = fmt.Sprintf("%s = int(%s)", nodeVarName(node.ID), node.Data.Value)
	case "Df_conditional":
		text = condicionalGenerator(node)
	case "Df_print":
		text = fmt.Sprintf("msg%v = \"%s\"", node.ID, node.Data.Value)
	}
	return text
}

func condicionalGenerator(node Node) string {
	thenId, err := strconv.Atoi(node.Outputs.Output1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	elseId, err := strconv.Atoi(node.Outputs.Output2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("if %s:\n    %s\nelse:\n    %s", node.Data.Condition, nodeVarName(thenId), nodeVarName(elseId))
}

func nodeVarName(id int) string {
	el := getNodeById(id)
	text := ""
	switch nodeType := el.HTML; nodeType {
	case "Df_add":
		text = fmt.Sprintf("add%v", el.ID)
	case "Df_number":
		text = fmt.Sprintf("number%v", el.ID)
	case "Df_print":
		text = fmt.Sprintf("print(msg%v)", el.ID)
	}
	return text
}

func getNodeById(id int) Node {

	var node Node

	for _, v := range nodes {
		if v.ID == id {
			node = v
		}
	}

	return node
}
