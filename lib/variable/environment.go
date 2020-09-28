package variable

import (
	"binadesa2020-backend/lib/clog"
	"errors"
	"io/ioutil"
	"log"
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

// MinioConfig for storage
var MinioConfig struct {
	URIEndpoint string // https
	EndPoint    string // IP internal (production) or https (local)
	AccessKey   string
	SecretKey   string
}

// Version this service
var Version string

// Initialization read from variable environment
func Initialization() {
	godotenv.Load("devel.env")

	// Reading version
	ver, err := ioutil.ReadFile("./VERSION")
	if err != nil {
		log.Fatalf("read version file %v \n", err)
	}
	Version = string(ver)

	MongoConfig.Host = os.Getenv("DB_MONGO_HOST")
	MongoConfig.Database = os.Getenv("DB_MONGO_NAME")
	MongoConfig.User = os.Getenv("DB_MONGO_USER")
	MongoConfig.Password = os.Getenv("DB_MONGO_PASS")

	MinioConfig.URIEndpoint = os.Getenv("MINIO_URIENDPOINT")
	MinioConfig.EndPoint = os.Getenv("MINIO_ENDPOINT")
	MinioConfig.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	MinioConfig.SecretKey = os.Getenv("MINIO_SECRET_KEY")

	JWTConfig.Key = os.Getenv("JWY_SECRET_KEY")

	ProjectName = os.Getenv("PROJECT_NAME")
	if ProjectName == "" {
		clog.Fatal(errors.New("PROJECT_NAME variable environment is null"), "init variable environment")
	}
}
