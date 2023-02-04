package server

import (
	"context"
	"log"
	"time"

	"github.com/opoccomaxao/shiki-github-graph/pkg/storage"
)

func (s *Server) processAnimeUpdate(ctx context.Context) {
	done := ctx.Done()

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := s.runSingleAnimeUpdate(ctx)
			if err != nil {
				log.Printf("%+v\n", err)
			}
		case <-s.animeProcessNotify:
			err := s.runSingleAnimeUpdate(ctx)
			if err != nil {
				log.Printf("%+v\n", err)
			}
		}
	}
}

func (s *Server) runSingleAnimeUpdate(ctx context.Context) error {
	for {
		ids, err := s.Storage.GetAnimesUnstored(ctx, 10)
		if err != nil {
			return err
		}

		if len(ids) == 0 {
			return nil
		}

		log.Printf("update animes %d\n", ids)

		res := make([]*storage.Anime, 0, len(ids))

		for _, id := range ids {
			anime, err := s.Shiki.GetAnime(ctx, int(id))
			if err != nil {
				log.Printf("%+v\n", err)
			}

			res = append(res, &storage.Anime{
				ID:                     int64(anime.ID),
				Name:                   anime.Name,
				Episodes:               anime.Episodes,
				EpisodeDurationSeconds: anime.Duration * 60,
				UpdatedAt:              time.Now().Unix(),
			})
		}

		err = s.Storage.SaveAnime(ctx, res)
		if err != nil {
			return err
		}
	}
}
