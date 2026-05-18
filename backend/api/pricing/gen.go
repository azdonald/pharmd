package pricing

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=pricing pricing.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=pricing pricing.yaml
