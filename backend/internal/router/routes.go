package router

type Route struct {
	Method string
	Path   string
	Name   string
}

var Routes = []Route{
	// Auth
	{"POST", "/api/auth/register", "Register User"},
	{"POST", "/api/auth/login", "Login User"},

	// PKI (future)
	{"POST", "/api/key/generate", "Generate RSA Key"},
	{"POST", "/api/csr/create", "Create CSR"},
	{"POST", "/api/cert/sign", "Sign Certificate"},
}