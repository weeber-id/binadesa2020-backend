package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/tools"
	"binadesa2020-backend/lib/variable"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// Admin models in mongoDB
type Admin struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:",omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password" json:"-"`
	Name     string             `bson:"name"`
	Level    int                `bson:"level"`
}

// Collection for admin data
func (Admin) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.Admin)
}

// FindByUsername and write to internal variable
func (a *Admin) FindByUsername(username string) bool {
	a.Collection().FindOne(context.Background(), bson.M{"username": username}).Decode(a)
	if (*a == Admin{}) {
		return false
	}
	return true
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
