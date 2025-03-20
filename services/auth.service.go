package services

import (
	"diandi-backend/domains"
	"diandi-backend/lib"
)

type AuthService struct {
	logger lib.Logger
	env    lib.Env
}

func NewAuthService(env lib.Env, logger lib.Logger) domains.AuthService {
	return AuthService{
		env:    env,
		logger: logger,
	}
}

func (as AuthService) Authorize(tokenString string) (bool, error) {
	return true, nil
}

func (as AuthService) CreateToken() string {
	return ""
}
