package main

import (
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!, I am monkey interpreter, send me some code to compile and run on /compile endpoint")
	})

	app.Post("/compile", handler)

	app.Listen(":8080")

}

func handler(c *fiber.Ctx) error {

	c.Set("Content-Type", "application/json")
	c.Set("Access-Control-Allow-Origin", "*")

	output := make([]string, 0)

	code := c.FormValue("code", "")

	env := object.NewEnvironment()
	l := lexer.New(code)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		output = printParserErrors( p.Errors())
		c.JSON(
			fiber.Map{
				"output": output,
				"eval":   "error",
			},
		)
		return nil
	}

	

	evaluator.Builtins["puts"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				output = append(output, arg.Inspect())
			}
			return evaluator.NULL
		},
	}

	evaluated := evaluator.Eval(program, env)


	c.JSON(fiber.Map{
		"output": output,
		"eval":   evaluated.Inspect(),
	})


	return nil

}

func printParserErrors( errors []string) []string {
	out := make([]string, 0)
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

	out = append(out, MONKEY_FACE)
	out = append(out, "Woops! We ran into some monkey business here!")
	out = append(out, " parser errors:")
	for _, msg := range errors {
		out = append(out, "\t"+msg)
	}
	return out
}
