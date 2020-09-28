package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// AktaKelahiran submission structur
type AktaKelahiran struct {
	BaseSubmission     `bson:",inline"`
	NamaKepalaKeluarga string             `bson:"nama_kepala_keluarga" json:"nama_kepala_keluarga"`
	Email              string             `bson:"email" json:"email"`
	File               AktaKelahiranFiles `bson:"file" json:"file"`
}

// AktaKelahiranFiles for file submission
type AktaKelahiranFiles struct {
	SuratKelahiran string `bson:"surat_kelahiran" json:"surat_kelahiran"`
	KTPSuami       string `bson:"ktp_suami" json:"ktp_suami"`
	KTPIstri       string `bson:"ktp_istri" json:"ktp_istri"`
	KTPSaksi1      string `bson:"ktp_saksi_1" json:"ktp_saksi_1"`
	KTPSaksi2      string `bson:"ktp_saksi_2" json:"ktp_saksi_2"`
	SuratNikah     string `bson:"surat_nikah" json:"surat_nikah"`
}

// Collection for submission akta kelahiran
func (AktaKelahiran) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.SubMission.AktaKelahiran)
}

// Create submission from this struct
// generate unique code
// modify date, insert to DB,
// get ID from database to this variable
func (a *AktaKelahiran) Create() (*mongo.InsertOneResult, error) {
	// generate unique code and date created & modified
	a.InitCreate()

	result, err := a.Collection().InsertOne(context.Background(), *a)
	if err != nil {
		return nil, err
	}
	a.ID = result.InsertedID.(primitive.ObjectID)
	return result, nil
}

// Update this struct to database
func (a *AktaKelahiran) Update() (*mongo.UpdateResult, error) {
	a.ModifiedAt = variable.DateTimeNowPtr()

	update := bson.M{"$set": *a}
	return a.Collection().UpdateOne(context.Background(), bson.M{"_id": a.ID}, update)
}
