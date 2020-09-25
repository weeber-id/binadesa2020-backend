package models

import (
	"binadesa2020-backend/lib/tools"
	"binadesa2020-backend/lib/variable"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Base struct for data in MongoDB
type Base struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt  *time.Time         `bson:"created_at,omitempty" json:"created_at"`
	ModifiedAt *time.Time         `bson:"modified_at,omitempty" json:"modified_at"`
	// DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

// BaseSubmission for submission
type BaseSubmission struct {
	Base       `bson:",inline"`
	UniqueCode string `bson:"unique_code" json:"unique_code"`

	// StatusCode
	// 0 waiting
	// 1 process
	// 2 rejected
	// 3 accepted
	StatusCode int `bson:"status_code" json:"status_code"`
}

// InitCreate from this struct
// generate unique code
// modify date
func (b *BaseSubmission) InitCreate() error {
	b.UniqueCode = tools.RandomString(6)
	b.StatusCode = 0

	now := variable.DateTimeNowPtr()
	b.CreatedAt = now
	b.ModifiedAt = now

	return nil
}
