package httpserver

import "../../atella"

// vectorChannelListener listen channel from client at same host as a server.
func (s *Server) vectorChannelListener(vectorChannel chan map[string]atella.HostVector) {
	for {
		v := <-vectorChannel
		s.MasterServerVector.SetElement(s.hostname, v)
	}
}
