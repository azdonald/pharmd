package purchases

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=purchases purchases.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=purchases purchases.yaml
