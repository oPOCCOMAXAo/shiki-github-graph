package server

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/opoccomaxao/shiki-github-graph/pkg/image"
)

func (s *Server) mwImage(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Security-Policy", "img-src *")
}

func (s *Server) initImages() error {
	err := os.MkdirAll(s.config.ImageDir, 0o777)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) prepareDummy(nick string) error {
	path := s.imagePath(nick)

	_, err := os.Stat(path)

	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	svgData, err := image.BuildSVG(image.BuildParams{
		Start: time.Now().AddDate(-1, 0, 0),
		End:   time.Now(),
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(path, []byte(svgData), 0o600)
	if err != nil {
		return err
	}

	return nil
}
