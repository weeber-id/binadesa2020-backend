package models

// JwtClaims structure
type JwtClaims struct {
	Role     string
	Username string
	Name     string
	Level    int
}
