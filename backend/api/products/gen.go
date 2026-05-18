package products

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=products products.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=products products.yaml
