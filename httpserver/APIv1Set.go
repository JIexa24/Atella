package httpserver

import (
	"encoding/json"
	"net/http"

	"../../atella"
)

// setVectorAPIv1 handle requset from agents and save received vectors.
func (s *Server) setVectorAPIv1(w http.ResponseWriter, r *http.Request) {
	if !s.IAmIsMaster {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	switch r.Method {
	case "POST":
		r.ParseForm()
		code := r.Header.Get("X-Atella-Auth")
		if code != s.code {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		currentClientHostname := r.Form.Get("hostname")
		vector := r.Form.Get("vector")
		if vector == "" || currentClientHostname == "" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		v := atella.NewVector()
		err := json.Unmarshal([]byte(vector), v)
		if err != nil {
			s.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.MasterServerVector.SetElement(currentClientHostname, v.List)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
