package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"

	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// KartuKeluarga submission structur
type KartuKeluarga struct {
	BaseSubmission     `bson:",inline"`
	Nama               string             `bson:"nama" json:"nama"`
	NamaKepalaKeluarga string             `bson:"nama_kepala_keluarga" json:"nama_kepala_keluarga"`
	Email              string             `bson:"email" json:"email"`
	File               KartuKeluargaFiles `bson:"file" json:"file"`
}

// KartuKeluargaFiles for file submission
type KartuKeluargaFiles struct {
	KTPSuami          string `bson:"ktp_suami" json:"ktp_suami"`
	KTPIstri          string `bson:"ktp_istri" json:"ktp_istri"`
	SuratNikah        string `bson:"surat_nikah" json:"surat_nikah"`
	AktaKelahiranAnak string `bson:"akta_kelahiran_anak" json:"akta_kelahiran_anak"`
}

// Collection for submission kartu keluarga
func (KartuKeluarga) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.SubMission.KartuKeluarga)
}

// Create submission from this struct
// generate unique code
// modify date, insert to DB,
// get ID from database
func (k *KartuKeluarga) Create() (*mongo.InsertOneResult, error) {
	k.InitCreate()

	result, err := k.Collection().InsertOne(context.Background(), *k)
	if err != nil {
		return nil, err
	}
	k.ID = result.InsertedID.(primitive.ObjectID)
	return result, nil
}

// GetByUniqueCode and write in this variable
// return isfound and error
func (k *KartuKeluarga) GetByUniqueCode(code string) (bool, error) {
	var empty KartuKeluarga

	err := k.Collection().FindOne(context.Background(), bson.M{"unique_code": code}).Decode(k)
	if err != nil {
		return false, err
	}

	if *k == empty {
		return false, nil
	}
	return true, nil
}

// Update this struct to database
func (k *KartuKeluarga) Update() (*mongo.UpdateResult, error) {
	k.ModifiedAt = variable.DateTimeNowPtr()

	update := bson.M{"$set": *k}
	return k.Collection().UpdateOne(context.Background(), bson.M{"_id": k.ID}, update)
}
