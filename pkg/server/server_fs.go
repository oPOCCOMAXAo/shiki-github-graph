package server

import "github.com/opoccomaxao/shiki-github-graph/pkg/app"

func (s *Server) imagePath(nick string) string {
	return app.DirData + "/" + nick + ".svg"
}
