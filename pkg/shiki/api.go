package shiki

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

const host = "https://shikimori.one"

type api struct {
	client  http.Client
	limiter *rate.Limiter
}

func newAPI() *api {
	return &api{
		client: http.Client{
			Timeout: time.Minute,
		},
		limiter: rate.NewLimiter(90/60, 5),
	}
}

type apiCall struct {
	Method string
	Path   string
	ResPtr interface{}
	Debug  bool
}

func (api *api) callAPI(ctx context.Context, call apiCall) error {
	err := api.limiter.Wait(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	url := host + call.Path

	req, err := http.NewRequestWithContext(ctx, call.Method, url, http.NoBody)
	if err != nil {
		return errors.WithStack(err)
	}

	res, err := api.client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	if call.Debug {
		log.Printf("%s %s\n\n%s\n", call.Method, url, body)
	}

	if call.ResPtr == nil {
		return nil
	}

	err = json.Unmarshal(body, call.ResPtr)
	if err != nil {
		if !call.Debug {
			log.Printf("%s %s\n\n%s\n", call.Method, url, body)
		}

		return errors.WithStack(err)
	}

	return nil
}
