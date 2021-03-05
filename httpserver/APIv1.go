package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// universalResponse describe tmplate for api responses in json format
type universalResponse struct {
	ResponseBody interface{} `json:"response"`
}

func newUniversalResponse() universalResponse {
	return universalResponse{
		ResponseBody: nil,
	}
}

// APIv1 mapping methods with locations.
func (s *Server) APIv1() {
	http.HandleFunc("/api/v1/test/code/200", s.testCode200APIv1)
	http.HandleFunc("/api/v1/test/code/401", s.testCode401APIv1)
	http.HandleFunc("/api/v1/test/code/404", s.testCode404APIv1)
	http.HandleFunc("/api/v1/test/code/405", s.testCode405APIv1)
	http.HandleFunc("/api/v1/test/code/500", s.testCode500APIv1)
	http.HandleFunc("/api/v1/test/code/501", s.testCode501APIv1)
	http.HandleFunc("/api/v1/ping", s.pingAPIv1)
	http.HandleFunc("/api/v1/echo/method", s.echoMethodAPIv1)
	http.HandleFunc("/api/v1/get/host", s.getHostAPIv1)
	http.HandleFunc("/api/v1/get/vector", s.getVectorAPIv1)
	http.HandleFunc("/api/v1/get/version", s.getVersionAPIv1)
	http.HandleFunc("/api/v1/get/hostname", s.getHostnameAPIv1)
	http.HandleFunc("/api/v1/set/vector", s.setVectorAPIv1)
}

func (s *Server) testCode200APIv1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) testCode401APIv1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (s *Server) testCode404APIv1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) testCode405APIv1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Server) testCode500APIv1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) testCode501APIv1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) pingAPIv1(w http.ResponseWriter, r *http.Request) {
	response := newUniversalResponse()
	response.ResponseBody = "pong"
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

func (s *Server) echoMethodAPIv1(w http.ResponseWriter, r *http.Request) {
	response := newUniversalResponse()
	response.ResponseBody = r.Method
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
