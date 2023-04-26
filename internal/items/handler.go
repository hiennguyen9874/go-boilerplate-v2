package items

import "net/http"

type Handlers interface {
	Create() func(w http.ResponseWriter, r *http.Request)
	Get() func(w http.ResponseWriter, r *http.Request)
	GetMulti() func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
	Update() func(w http.ResponseWriter, r *http.Request)
}
