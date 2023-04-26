package users

import "net/http"

type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	Get() func(w http.ResponseWriter, r *http.Request)
	GetMulti() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
	Update() func(w http.ResponseWriter, r *http.Request)
	Me() func(w http.ResponseWriter, r *http.Request)
	UpdateMe() func(w http.ResponseWriter, r *http.Request)
	UpdatePassword() func(w http.ResponseWriter, r *http.Request)
	UpdatePasswordMe() func(w http.ResponseWriter, r *http.Request)
	LogoutAllAdmin() func(w http.ResponseWriter, r *http.Request)
}
