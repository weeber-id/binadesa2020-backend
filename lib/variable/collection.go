package variable

type collection struct {
	Admin      string
	User       string
	Complaint  string
	SubMission submission
}

type submission struct {
	KartuKeluarga string
}

// CollectionNames in MongoDB
var CollectionNames collection = collection{
	Admin:     "admin",
	User:      "user",
	Complaint: "complaint",
	SubMission: submission{
		KartuKeluarga: "submissionKartuKeluarga",
	},
}

// ProjectName for vokasi: bina desa
var ProjectName string = "vokasi-binadesa"
