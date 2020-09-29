package variable

type collection struct {
	Admin      string
	User       string
	Complaint  string
	News       string
	SubMission submission
}

type submission struct {
	KartuKeluarga string
	AktaKelahiran string
}

// CollectionNames in MongoDB
var CollectionNames collection = collection{
	Admin:     "admin",
	User:      "user",
	Complaint: "complaint",
	News:      "news",
	SubMission: submission{
		KartuKeluarga: "submissionKartuKeluarga",
		AktaKelahiran: "submissionAktaKelahiran",
	},
}

// ProjectName for vokasi: bina desa
var ProjectName string
