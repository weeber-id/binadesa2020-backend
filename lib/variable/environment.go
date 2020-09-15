package variable

import (
	"os"

	"github.com/joho/godotenv"
)

// MongoConfig data type
var MongoConfig struct {
	Host     string
	Database string
	User     string
	Password string
}

// JWTConfig datatype
var JWTConfig struct {
	Key string
}

// Initialization read from variable environment
func Initialization() {
	godotenv.Load("devel.env")

	MongoConfig.Host = os.Getenv("DB_MONGO_HOST")
	MongoConfig.Database = os.Getenv("DB_MONGO_NAME")
	MongoConfig.User = os.Getenv("DB_MONGO_USER")
	MongoConfig.Password = os.Getenv("DB_MONGO_PASS")

	JWTConfig.Key = os.Getenv("JWY_SECRET_KEY")
}
