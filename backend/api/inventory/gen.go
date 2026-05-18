package inventory

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=inventory inventory.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=inventory inventory.yaml
