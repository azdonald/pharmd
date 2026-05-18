package users

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=users user.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=users user.yaml
