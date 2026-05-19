package pos

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=pos pos.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=pos pos.yaml
