package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"
	"context"

	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

// KartuKeluarga submission structur
type KartuKeluarga struct {
	Base               `bson:",inline"`
	NamaKepalaKeluarga string             `bson:"nama_kepala_keluarga"`
	NoTelp             string             `bson:"no_telp"`
	File               KartuKeluargaFiles `bson:"file"`
}

// KartuKeluargaFiles for file submission
type KartuKeluargaFiles struct {
	KTPSuami          string `bson:"ktp_suami"`
	KTPIstri          string `bson:"ktp_istri"`
	SuratNikah        string `bson:"surat_nikah"`
	AktaKelahiranAnak string `bson:"akta_kelahiran_anak"`
}

// Collection for submission kartu keluarga
func (KartuKeluarga) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.SubMission.KartuKeluarga)
}

// Create submission from this struct
// Modify date, insert to DB, and get ID
func (k *KartuKeluarga) Create() (*mongo.InsertOneResult, error) {
	now := variable.DateTimeNowPtr()
	k.CreatedAt = now
	k.ModifiedAt = now

	result, err := k.Collection().InsertOne(context.Background(), *k)
	if err != nil {
		return nil, err
	}

	k.ID = result.InsertedID.(primitive.ObjectID)
	return result, nil
}

// ChangeAllFiles submission to this struct
// Change files in this struct, update to DB based on this ID
func (k *KartuKeluarga) ChangeAllFiles(files *KartuKeluargaFiles) (*mongo.UpdateResult, error) {
	k.File = *files

	update := bson.M{"$set": *k}
	result, err := k.Collection().UpdateOne(context.Background(), bson.M{"_id": k.ID}, update)
	if err != nil {
		return nil, err
	}
	k.ModifiedAt = variable.DateTimeNowPtr()
	return result, nil
}
