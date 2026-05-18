package product_categories

//go:generate oapi-codegen -generate=chi-server -o=server_gen.go -package=product_categories categories.yaml
//go:generate oapi-codegen -generate=types -o=types_gen.go -package=product_categories categories.yaml
