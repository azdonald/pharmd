package permissions

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=permissions permission.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=permissions permission.yaml
