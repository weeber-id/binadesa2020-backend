package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/tools"
	"binadesa2020-backend/lib/variable"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Admin models in mongoDB
type Admin struct {
	Base     `bson:",inline"`
	Username string `bson:"username"`
	Password string `bson:"password" json:"-"`
	Name     string `bson:"name"`
	Level    int    `bson:"level"`
}

// Collection for admin data
func (Admin) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.Admin)
}

// Create admin from this struct
func (a *Admin) Create() (*mongo.InsertOneResult, error) {
	a.CreatedAt = variable.DateTimeNowPtr()
	a.ModifiedAt = variable.DateTimeNowPtr()
	return a.Collection().InsertOne(context.Background(), *a)
}

// FindByUsername and write to internal variable
func (a *Admin) FindByUsername(username string) bool {
	a.Collection().FindOne(context.Background(), bson.M{"username": username}).Decode(a)
	if (*a == Admin{}) {
		return false
	}
	return true
}

// DeleteByUsername static method
func (a Admin) DeleteByUsername(username string) *mongo.SingleResult {
	return a.Collection().FindOneAndDelete(context.Background(), bson.M{"username": username})
}

// Verify username and password return true is verify and vice versa
func (a *Admin) Verify(username, password string) bool {
	found := a.FindByUsername(username)
	if !found {
		return false
	}

	if a.Password != tools.EncodeMD5(password) {
		return false
	}
	return true
}
