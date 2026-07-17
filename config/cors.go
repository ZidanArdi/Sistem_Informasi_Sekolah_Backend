package config

var allowedOrigins = []string{
	"http://localhost:5173",
	"http://localhost:5174",
	"https://sisteminformasisekolahbackend-production.up.railway.app",
}

func GetAllowedOrigins() []string {
	return allowedOrigins
}
