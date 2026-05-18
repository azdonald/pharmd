package patients

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=patients patients.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=patients patients.yaml
