package config

var (
	// API Settings
	APP_PORT = 5000
	MODE     = "DEV"
	// Auth
	JWT_SECRET  = `Enter Your Secret`
	JWT_EXPIRES = int64(84600) // One Day
	// DB Settings
	DB_USERNAME = "test"
	DB_PASSWORD = "123"
	DB_HOST     = "db"
	DB_DATABASE = "test"
	DB_PORT     = "3306"
)
