package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/tools"
	"binadesa2020-backend/lib/variable"
	"context"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

// News Structure
type News struct {
	Base       `bson:",inline"`
	Author     string `bson:"author" json:"author"`
	Title      string `bson:"title" json:"title"`
	ImageCover string `bson:"image_cover" json:"image_cover"`
	Content    string `bson:"content" json:"content"`

	// automatically generate when content created
	Slug string `bson:"slug" json:"slug"`
}

// Collection mongodb for this struct
func (News) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.News)
}

// Create news object
// modify date,, slug, and insert to DB,
// get ID from database to this variable
func (n *News) Create() (*mongo.InsertOneResult, error) {
	n.InitDate()
	n.Slug = url.QueryEscape(n.Title) + tools.RandomString(2)

	result, err := n.Collection().InsertOne(context.Background(), *n)
	if err != nil {
		return nil, err
	}
	n.ID = result.InsertedID.(primitive.ObjectID)
	return result, nil
}

// GetByID news and write into this struct
func (n *News) GetByID(ID string) error {
	return n.Collection().FindOne(context.Background(), bson.M{"_id": ID}).Decode(n)
}

// GetBySlug news and write into this struct
func (n *News) GetBySlug(slug string) error {
	return n.Collection().FindOne(context.Background(), bson.M{"slug": slug}).Decode(n)
}

// Update news and save to database
func (n *News) Update() (*mongo.UpdateResult, error) {
	n.ModifiedAt = variable.DateTimeNowPtr()

	update := bson.M{"$set": *n}
	return n.Collection().UpdateOne(context.Background(), bson.M{"_id": n.ID}, update)
}
