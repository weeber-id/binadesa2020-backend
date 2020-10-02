package models

import (
	"binadesa2020-backend/lib/services/mongodb"
	"binadesa2020-backend/lib/variable"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

// SuratKeterangan submission structur
type SuratKeterangan struct {
	BaseSubmission `bson:",inline"`
	Tipe           string               `bson:"tipe" json:"tipe"`
	Nama           string               `bson:"nama" json:"nama"`
	Email          string               `bson:"email" json:"email"`
	IsPaid         bool                 `bson:"is_paid" json:"is_paid"`
	File           SuratKeteranganFiles `bson:"file" json:"file"`
}

// SuratKeteranganFiles for file submission
type SuratKeteranganFiles struct {
	SuratPernyataan   string `bson:"surat_pernyataan" json:"surat_pernyataan"`
	KTP               string `bson:"ktp" json:"ktp"`
	LampiranPendukung string `bson:"lampiran_pendukung" json:"lampiran_pendukung"`
}

// Collection for submission surat keterangan
func (SuratKeterangan) Collection() *mongo.Collection {
	return mongodb.DB.Collection(variable.CollectionNames.SubMission.SuratKeterangan)
}

// Create submission from this struct
// check 'tipe' variable
// generate unique code
// write to database
// get ID from database to this variable
func (s *SuratKeterangan) Create() (*mongo.InsertOneResult, error) {
	// Check tipe variable
	switch s.Tipe {
	case "Tanah Bebas Sengketa",
		"Kepemilikan Tanah",
		"Pindah",
		"Hak Pergi Haji",
		"Keringanan Pembiayaan Kendaraan Bermotor",
		"Salah Nama":
		s.IsPaid = true

	case "Meninggal Dunia", "Tidak Mampu":
		s.IsPaid = false

	default:
		return nil, errors.New("'tipe' variable must be: Tanah Bebas Sengketa, Kepemilikan Tanah, Pindah, Hak Pergi Haji, Keringanan Pembiayaan Kendaraan Bermotor, Salah Nama")
	}

	// generate unique code and date created & modified
	s.InitCreate()

	// write to database
	result, err := s.Collection().InsertOne(context.Background(), *s)
	if err != nil {
		return nil, err
	}
	s.ID = result.InsertedID.(primitive.ObjectID)
	return result, nil
}

// GetByUniqueCode surat keterangan and write into this struct
func (s *SuratKeterangan) GetByUniqueCode(code string) (bool, error) {
	var empty SuratKeterangan

	err := s.Collection().FindOne(context.Background(), bson.M{"unique_code": code}).Decode(s)
	if err != nil {
		return false, err
	}

	if *s == empty {
		return false, nil
	}
	return true, nil
}

// Update this struct to database
func (s *SuratKeterangan) Update() (*mongo.UpdateResult, error) {
	s.ModifiedAt = variable.DateTimeNowPtr()

	update := bson.M{"$set": *s}
	return s.Collection().UpdateOne(context.Background(), bson.M{"_id": s.ID}, update)
}
