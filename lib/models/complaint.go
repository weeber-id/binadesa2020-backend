package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Complaint struct
// Untuk kolom pengaduan dari warga
type Complaint struct {
	Base      `bson:",inline"`
	Name      string `bson:"name" json:"name"`
	RT        string `bson:"RT" json:"rt"`
	RW        string `bson:"RW" json:"rw"`
	Address   string `bson:"address" json:"address"`
	Complaint string `bson:"complaint" json:"complaint"`
}

// Collection for complaint data
func (Complaint) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.Complaint)
}

// Create a complaint data
func (c *Complaint) Create() (*mongo.InsertOneResult, error) {
	c.CreatedAt = variable.DateTimeNowPtr()
	c.ModifiedAt = variable.DateTimeNowPtr()
	return c.Collection().InsertOne(context.Background(), *c)
}
