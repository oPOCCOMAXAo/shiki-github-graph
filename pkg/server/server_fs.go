package server

func (s *Server) imagePath(nick string) string {
	return s.config.ImageDir + "/" + nick + ".svg"
}
