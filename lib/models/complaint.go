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
	Name      string `bson:"name"`
	RT        string `bson:"RT"`
	RW        string `bson:"RW"`
	Address   string `bson:"address"`
	Complaint string `bson:"complaint"`
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
