package variable

type collection struct {
	Admin string
	User  string
}

// CollectionNames in MongoDB
var CollectionNames collection = collection{
	Admin: "admin",
	User:  "user",
}
