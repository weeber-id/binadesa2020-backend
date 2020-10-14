package models

import (
	"binadesa2020-backend/lib/tools"
	"binadesa2020-backend/lib/variable"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type statusCode struct {
	Waiting  int
	Process  int
	Rejected int
	Accepted int
}

// StatusCode for submission
var StatusCode statusCode = statusCode{
	Waiting:  0,
	Process:  1,
	Rejected: 2,
	Accepted: 3,
}

// Base struct for data in MongoDB
type Base struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt  *time.Time         `bson:"created_at,omitempty" json:"created_at"`
	ModifiedAt *time.Time         `bson:"modified_at,omitempty" json:"modified_at"`
}

// InitDate in createdAt and modifiedAt
func (b *Base) InitDate() error {
	now := variable.DateTimeNowPtr()
	b.CreatedAt = now
	b.ModifiedAt = now
	return nil
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
// modify date
// generate unique code
func (b *BaseSubmission) InitCreate() error {
	b.InitDate()

	b.UniqueCode = tools.RandomString(6)
	b.StatusCode = StatusCode.Waiting
	return nil
}
