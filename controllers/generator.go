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
		text = fmt.Sprintf("%s = float(%v)", nodeVarName(node.ID), node.Data.Value)
	case "Df_conditional":
		text = condicionalGenerator(node)
	case "Df_for":
		text = forGenerator(node)
	default:
		text = ""
	}
	return text
}

func addGenerator(node Node) string {
	var text string

	if len(node.Inputs.Input1.Connections) <= 0 || len(node.Inputs.Input2.Connections) <= 0 {
		text = ""
		return text
	}

	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	text = fmt.Sprintf("add%v = %s + %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
	return text
}

func divideGenerator(node Node) string {
	var text string

	if len(node.Inputs.Input1.Connections) <= 0 || len(node.Inputs.Input2.Connections) <= 0 {
		text = ""
		return text
	}

	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	text = fmt.Sprintf("divide%v = %s / %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
	return text
}

func multiplyGenerator(node Node) string {
	var text string

	if len(node.Inputs.Input1.Connections) <= 0 || len(node.Inputs.Input2.Connections) <= 0 {
		text = ""
		return text
	}

	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	text = fmt.Sprintf("multiply%v = %s * %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
	return text
}

func substractionGenerator(node Node) string {
	var text string

	if len(node.Inputs.Input1.Connections) <= 0 || len(node.Inputs.Input2.Connections) <= 0 {
		text = ""
		return text
	}

	id1, err := strconv.Atoi(node.Inputs.Input1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	id2, err := strconv.Atoi(node.Inputs.Input2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	text = fmt.Sprintf("substraction%v = %s - %s\n", node.ID, nodeVarName(id1), nodeVarName(id2))
	return text
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
	default:
		text = ""
	}
	return text
}

func getSentenceCode(node Node) string {
	text := ""

	switch nodeType := node.HTML; nodeType {
	case "Df_code":
		text = node.Data.Value
	case "Df_print":
		text = fmt.Sprintf("print(\"%s\")", node.Data.Value)
	default:
		text = ""
	}
	return text
}

func condicionalGenerator(node Node) string {
	var text string

	if len(node.Outputs.Output1.Connections) <= 0 || len(node.Outputs.Output2.Connections) <= 0 {
		text = fmt.Sprintf("if %s:\n    \nelse:\n    ", node.Data.Condition)
		return text
	}

	thenId, err := strconv.Atoi(node.Outputs.Output1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}
	ifNode := getNodeById(thenId)

	elseId, err := strconv.Atoi(node.Outputs.Output2.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}
	elseNode := getNodeById(elseId)

	text = fmt.Sprintf("if %s:\n    %s\nelse:\n    %s", node.Data.Condition, getSentenceCode(ifNode), getSentenceCode(elseNode))
	return text
}

func forGenerator(node Node) string {
	var text string

	if len(node.Outputs.Output1.Connections) <= 0 {
		text = fmt.Sprintf("for i in range(%s,%s):\n", node.Data.From, node.Data.Till)
		return text
	}

	outputId, err := strconv.Atoi(node.Outputs.Output1.Connections[0].Node)
	if err != nil {
		log.Fatal(err)
	}

	nodeOutput := getNodeById(outputId)

	text = fmt.Sprintf("for i in range(%s,%s):\n    %s\n", node.Data.From, node.Data.Till, getSentenceCode(nodeOutput))

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
