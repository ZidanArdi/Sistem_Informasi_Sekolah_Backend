package config

var allowedOrigins = []string{
	"http://localhost:5173",
	"http://localhost:5174",
	"https://my-fe-ten.vercel.app", // This will be replaced with the actual Vercel URL
}

func GetAllowedOrigins() []string {
	return allowedOrigins
}
