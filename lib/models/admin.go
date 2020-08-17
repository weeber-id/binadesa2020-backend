package models

import "github.com/go-bongo/bongo"

// Admin models in mongoDB
type Admin struct {
	bongo.DocumentBase `bson:",inline"`
	Username           string
	Password           string
	Name               string
	Level              int8
}
