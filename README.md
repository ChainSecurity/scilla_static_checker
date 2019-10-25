# scilla_checker
## Project Structure
*  `pkg/ast` Scilla AST related code
*  `pkg/ir` IR related code

## How to run

There is exists a graphical output of IR.
Example:

`go run cmd/graph_plot/main.go examples/inc.json a.dot`

`dot -Tpng -O a.dot`

To get more dot graphs, code modifications in `pkg/ir/plot.go` are needed.
