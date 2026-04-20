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
	{"POST", "/api/auth/refresh", "Refresh JWT Token"},
	{"POST", "/api/auth/logout", "Logout User (Protected)"},
	{"GET", "/api/auth/me", "Get Current User (Protected)"},

	// PKI (future)
	{"POST", "/api/key/generate", "Generate RSA Key"},
	{"POST", "/api/csr/create", "Create CSR"},
	{"POST", "/api/cert/sign", "Sign Certificate"},
}