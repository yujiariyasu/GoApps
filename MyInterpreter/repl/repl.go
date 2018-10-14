package repl

import (
	"bufio"
	"fmt"
	"github.com/yujiariyasu/GoApps/MyInterpreter/lexer"
	"github.com/yujiariyasu/GoApps/MyInterpreter/parser"
	// "github.com/yujiariyasu/GoApps/MyInterpreter/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		fmt.Printf("%+v", program)
	}
}
