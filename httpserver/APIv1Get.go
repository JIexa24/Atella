package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// getHostAPIv1 return json, contains hostname and agent version.
func (s *Server) getHostAPIv1(w http.ResponseWriter, r *http.Request) {
	var m map[string]string = make(map[string]string, 0)
	m["hostname"] = s.hostname
	m["version"] = s.Version
	response := newUniversalResponse()
	response.ResponseBody = m
	responseJSON, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "%s", responseJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getVectorAPIv1 return agent master vector.
func (s *Server) getVectorAPIv1(w http.ResponseWriter, r *http.Request) {
	response := newUniversalResponse()
	response.ResponseBody = s.MasterServerVector.GetVectorCopy()
	responseJSON, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "%s", responseJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// getVersionAPIv1 return agent version.
func (s *Server) getVersionAPIv1(w http.ResponseWriter, r *http.Request) {
	response := newUniversalResponse()
	response.ResponseBody = s.Version
	responseJSON, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "%s", responseJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) getHostnameAPIv1(w http.ResponseWriter, r *http.Request) {
	response := newUniversalResponse()
	response.ResponseBody = s.hostname
	responseJSON, err := json.Marshal(response)
	if err != nil {
		s.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "%s", responseJSON)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
