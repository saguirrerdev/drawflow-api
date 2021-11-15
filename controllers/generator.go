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
	From      string `json:"from,omitempty"`
	Till      string `json:"till,omitempty"`
}
type ConnectionsInput struct {
	Node  string `json:"node,omitempty"`
	Input string `json:"input,omitempty"`
}
type Input struct {
	Connections []ConnectionsInput `json:"connections,omitempty"`
}
type Inputs struct {
	Input1 Input `json:"input_1,omitempty"`
	Input2 Input `json:"input_2,omitempty"`
	Input3 Input `json:"input_3,omitempty"`
}
type ConnectionsOutput struct {
	Node   string `json:"node,omitempty"`
	Output string `json:"output,omitempty"`
}
type Output struct {
	Connections []ConnectionsOutput `json:"connections,omitempty"`
}
type Outputs struct {
	Output1 Output `json:"output_1,omitempty"`
	Output2 Output `json:"output_2,omitempty"`
	Output3 Output `json:"output_3,omitempty"`
}

type Code struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var nodes Nodes

func GetCode(w http.ResponseWriter, r *http.Request) {
	node := r.Context().Value("node").(*models.Node)

	err := json.Unmarshal([]byte(node.Nodes), &nodes)
	if err != nil {
		log.Fatal(err)
	}

	code := Code{Code: generate("", 0), Name: node.Name}

	err1 := render.Render(w, r, CodeResponse(code))
	if err1 != nil {
		log.Fatal(err1)
	}
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

func getNodeCode(node Node) string {
	text := ""

	switch nodeType := node.HTML; nodeType {
	case "Df_add":
		text = addGenerator(node)
	case "Df_substraction":
		text = substractionGenerator(node)
	case "Df_divide":
		text = divideGenerator(node)
	case "Df_multiply":
		text = multiplyGenerator(node)
	case "Df_number":
		text = fmt.Sprintf("%s = int(%v)", nodeVarName(node.ID), node.Data.Value)
	case "Df_conditional":
		text = condicionalGenerator(node)
	case "Df_print":
		text = fmt.Sprintf("msg%v = \"%v\"", node.ID, node.Data.Value)
	case "Df_for":
		text = forGenerator(node)
	case "Df_code":
		text = fmt.Sprintf("%s", node.Data.Value)
	default:
		text = "Node code no available"
	}
	return text
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

	return fmt.Sprintf("add%v = %s + %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
}

func divideGenerator(node Node) string {
	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("divide%v = %s / %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
}
func multiplyGenerator(node Node) string {
	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("multiply%v = %s * %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
}

func substractionGenerator(node Node) string {
	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("substraction%v = %s - %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
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

func forGenerator(node Node) string {
	return fmt.Sprintf("for i in range(%s,%s):\n    %s\n", node.Data.From, node.Data.Till, "print(\"Hello world\")")
}

func nodeVarName(id int) string {
	el := getNodeById(id)
	text := ""
	switch nodeType := el.HTML; nodeType {
	case "Df_add":
		text = fmt.Sprintf("add%v", el.ID)
	case "Df_substractionadd":
		text = fmt.Sprintf("substraction%v", el.ID)
	case "Df_divide":
		text = fmt.Sprintf("divide%v", el.ID)
	case "Df_multiply":
		text = fmt.Sprintf("multiply%v", el.ID)
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
			break
		}
	}

	return node
}
