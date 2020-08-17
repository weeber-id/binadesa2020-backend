package main

import (
	"binadesa2020-backend/lib/controllers"
	"binadesa2020-backend/lib/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// type Admin struct {
// 	bongo.DocumentBase `bson:",inline"`
// 	Username           string
// 	Password           string
// 	Name               string
// 	Level              int8
// }

func main() {
	godotenv.Load("devel.env")

	services.MongodbConnection()
	defer services.MDB.Session.Close()

	// admin1 := &Admin{
	// 	Username: "bayu3490",
	// 	Password: tools.EncodeMD5("testing"),
	// 	Name:     "Bayu Aditya",
	// 	Level:    0,
	// }

	// err := conn.Collection("admin").Save(admin1)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// row := &models.Admin{}
	// results := conn.Collection("admin").Find(bson.M{})
	// for results.Next(row) {
	// 	fmt.Println(row.Id)
	// }

	r := gin.Default()

	r.GET("/admins", controllers.GetAdmin)
	r.Run()
}
