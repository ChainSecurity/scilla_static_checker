# scilla_checker
## Project Structure
*  `pkg/ast` Scilla AST related code
*  `pkg/ir` IR related code

## How to run static checker

`go run cmd/scilla_static/main.go examples/inc.json`


## How to plot IR dot graph

There is exists a graphical output of IR.
Example:

`go run cmd/graph_plot/main.go examples/inc.json a.dot`

`dot -Tpng -O a.dot`
