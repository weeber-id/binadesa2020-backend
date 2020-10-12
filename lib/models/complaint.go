package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

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
	IsRead    bool   `bson:"is_read" json:"is_read"`
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

// GetByID complaint and write into this struct
func (c *Complaint) GetByID(ID string) (bool, error) {
	var empty Complaint

	objectID, _ := primitive.ObjectIDFromHex(ID)
	err := c.Collection().FindOne(context.Background(), bson.M{"_id": objectID}).Decode(c)
	if err != nil {
		return false, err
	}

	if *c == empty {
		return false, nil
	}

	if c.IsRead == false {
		c.IsRead = true
		c.Update()
	}
	return true, nil
}

// Update this struct to database
func (c *Complaint) Update() (*mongo.UpdateResult, error) {
	c.ModifiedAt = variable.DateTimeNowPtr()

	update := bson.M{"$set": *c}
	return c.Collection().UpdateOne(context.Background(), bson.M{"_id": c.ID}, update)
}
