package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base struct for data in MongoDB
type Base struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:",omitempty"`
	CreatedAt  *time.Time         `bson:"created_at,omitempty"`
	ModifiedAt *time.Time         `bson:"modified_at,omitempty"`
	// DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}
