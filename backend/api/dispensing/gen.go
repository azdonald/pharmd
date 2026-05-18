package dispensing

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=dispensing dispensing.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=dispensing dispensing.yaml
