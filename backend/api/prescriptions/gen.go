package prescriptions

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=prescriptions prescriptions.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=prescriptions prescriptions.yaml
