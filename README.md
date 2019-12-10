# scilla_checker
## Project Structure
*  `pkg/ast` Scilla AST related code
*  `pkg/ir` IR related code


## Installing the static checker

Requirements:
- Go version `1.12` or higher

Run the following command from the root folder of the project:

`go mod init github.com/ChainSecurity/scilla_static_checker`

## How to run static checker

`go run cmd/scilla_static/main.go examples/inc.json`

## How to plot IR dot graph

There is exists a graphical output of IR.
Example:

`go run cmd/graph_plot/main.go examples/inc.json a.dot`

`dot -Tpng -O a.dot`
