set -e
for j in ./examples/*.json ;
do
    echo $j
    go run ./cmd/scilla_static/main.go $j
done
