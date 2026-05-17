package utils

var ValidationExemptRoutes = map[string]bool{
	"/v1/register":   true,
	"/v1/login":      true,
	"/v1/refresh":    true,
	"/health":        true,
}
