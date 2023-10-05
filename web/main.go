package main

import (
	"encoding/json"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"net/http"
)

func main() {
	http.HandleFunc("/compile", handler)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	 var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, ok := data["code"].(string)
    if !ok {
        http.Error(w, "Missing or invalid 'name' field in JSON", http.StatusBadRequest)
        return
    }
	env := object.NewEnvironment()
	l := lexer.New(code)
	p := parser.New(l)
	
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(w, p.Errors())
		return
	}

	evaluator.Builtins["puts"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				io.WriteString(w, arg.Inspect())
				io.WriteString(w, "\n")
			}
			return evaluator.NULL
		},
	}
	evaluated :=  evaluator.Eval(program, env)
	if evaluated != nil {
		if evaluated.Inspect() != "null" {
			io.WriteString(w, evaluated.Inspect())
			io.WriteString(w, "\n")
		}

	}

}

func printParserErrors(out io.Writer, errors []string) {
	const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
