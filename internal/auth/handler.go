package auth

import "net/http"

type Handlers interface {
	SignIn() func(w http.ResponseWriter, r *http.Request)
	RefreshToken() func(w http.ResponseWriter, r *http.Request)
	GetPublicKey() func(w http.ResponseWriter, r *http.Request)
	Logout() func(w http.ResponseWriter, r *http.Request)
	LogoutAllToken() func(w http.ResponseWriter, r *http.Request)
	VerifyEmail() func(w http.ResponseWriter, r *http.Request)
	ForgotPassword() func(w http.ResponseWriter, r *http.Request)
	ResetPassword() func(w http.ResponseWriter, r *http.Request)
}
