package shiki

import (
	"context"
	"fmt"

	"github.com/opoccomaxao-go/generic-collection/slice"
)

type HistoryResult struct {
	Entries []*HistoryEntry
	Warning []string
	Error   error
}

func (s *Service) LoadHistory(ctx context.Context, userID int64) <-chan HistoryResult {
	res := make(chan HistoryResult, 10)

	go s.fetchAllHistory(ctx, userID, res)

	return res
}

func (s *Service) fetchAllHistory(ctx context.Context, userID int64, resChan chan<- HistoryResult) {
	defer close(resChan)

	done := ctx.Done()

	call := HistoryCall{
		UserID:          userID,
		Page:            0,
		Limit:           100,
		TargetTypeAnime: true,
	}

	for {
		res, err := s.api.GetHistory(ctx, call)
		if err != nil || len(res) == 0 {
			resChan <- HistoryResult{
				Error: err,
			}

			return
		}

		sendResult := HistoryResult{
			Entries: slice.Filter(res, func(he *HistoryEntry) bool {
				return he.Watch != nil
			}),
		}

		for _, entry := range sendResult.Entries {
			if entry.Watch.Failed {
				sendResult.Warning = append(sendResult.Warning, fmt.Sprintf(`Failed on "%s"`, entry.Description))
			}

			if entry.Watch.Finished {
				entry.Watch.EpisodeEnd = entry.Target.Episodes

				if entry.Target.Episodes == 1 {
					entry.Watch.EpisodeStart = 1
				}
			}
		}

		select {
		case <-done:
			return
		case resChan <- sendResult:
		}

		call.Page++
	}
}
