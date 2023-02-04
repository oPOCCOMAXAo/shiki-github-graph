package shiki

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type HistoryCall struct {
	UserID          int64
	Page            int
	Limit           int
	TargetTypeAnime bool
}

func (api *api) GetHistory(ctx context.Context, call HistoryCall) ([]*HistoryEntry, error) {
	var res []*HistoryEntry

	target := ""
	if call.TargetTypeAnime {
		target = "&target_type=Anime"
	}

	err := api.callAPI(ctx, apiCall{
		Method: http.MethodGet,
		Path: fmt.Sprintf("/api/users/%d/history?limit=%d&page=%d%s",
			call.UserID,
			call.Limit,
			call.Page,
			target,
		),
		ResPtr: &res,
	})
	if err != nil {
		return nil, err
	}

	for _, he := range res {
		he.Watch, err = Parser.ParseDescription(he.Description)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return res, nil
}
