package shiki

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	api := newAPI()

	res, err := api.GetHistory(context.Background(), HistoryCall{
		UserID: 843357,
		Page:   0,
		Limit:  100,
	})
	require.NoError(t, err)
	t.Logf("%#v", res)
}
