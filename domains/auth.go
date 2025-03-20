package domains

type AuthService interface {
	Authorize(tokenString string) (bool, error)
	CreateToken() string
}
