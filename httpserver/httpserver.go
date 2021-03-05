package httpserver

import (
	// "net/http"
	"fmt"
	"net/http"
	"sync"

	"../../atella"
)

// Server is server config structure.
type Server struct {
	Hosts              []atella.Host
	SelfIndexes        []int64
	hostname           string
	Version            string
	code               string
	IAmIsMaster        bool
	logger             atella.Logger
	MasterServerVector *atella.MasterVector
}

// NewHTTPServer return a new HTTP server descriptor.
func NewHTTPServer(Hosts []atella.Host, parsedHosts *atella.ParsedHosts,
	version, hostname, code string,
	logger atella.Logger) *Server {
	s := &Server{
		Hosts:       Hosts,
		SelfIndexes: parsedHosts.SelfIndexes,
		hostname:    hostname,
		Version:     version,
		code:        code,
		IAmIsMaster: atella.SubsetInt64(parsedHosts.MasterIndexes,
			parsedHosts.SelfIndexes),
		logger: logger,
		MasterServerVector : atella.NewMasterVector()}
	s.APIv1()
	return s
}

// Start server on multiple ports, when listed in Hosts by SelfIndexes.
func (s *Server) Start(vectorChannel chan map[string]atella.HostVector) error {
	for _, index := range s.SelfIndexes {
		i := index
		wait := &sync.WaitGroup{}
		wait.Add(1)
		go func() {
			addr := fmt.Sprintf(":%s", s.Hosts[i].Port)
			s.logger.Infof("Server started on %s", addr)
			wait.Done()
			if err := http.ListenAndServe(addr, nil); err != nil {
				s.logger.Fatalf("%s", err.Error())
			}
		}()
		wait.Wait()
	}

	if s.IAmIsMaster && vectorChannel != nil {
		s.logger.Info("Start master server")
		go s.vectorChannelListener(vectorChannel)
	}

	return nil
}
