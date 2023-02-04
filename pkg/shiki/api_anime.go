package shiki

import (
	"context"
	"fmt"
	"net/http"
)

func (api *api) GetAnime(ctx context.Context, id int) (*Anime, error) {
	var res Anime

	err := api.callAPI(ctx, apiCall{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/animes/%d", id),
		ResPtr: &res,
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}
