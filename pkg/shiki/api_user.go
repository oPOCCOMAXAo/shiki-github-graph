package shiki

import (
	"context"
	"net/http"
)

func (api *api) GetUserByNick(ctx context.Context, nick string) (*User, error) {
	var res User

	err := api.callAPI(ctx, apiCall{
		Method: http.MethodGet,
		Path:   "/api/users/" + nick,
		ResPtr: &res,
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}
