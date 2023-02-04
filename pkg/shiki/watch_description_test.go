package shiki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWatchDescription(t *testing.T) {
	testCases := []struct {
		input  string
		output WatchDescription
	}{
		{
			input: "Просмотрено",
			output: WatchDescription{
				Finished: true,
			},
		},
		{
			input: "Просмотрено и оценено на <b>5</b>",
			output: WatchDescription{
				Finished: true,
			},
		},
		{
			input: "Просмотрен 2-й эпизод",
			output: WatchDescription{
				EpisodeStart: 2,
				EpisodeEnd:   2,
			},
		},
		{
			input: "Просмотрены 8-й и 9-й эпизоды",
			output: WatchDescription{
				EpisodeStart: 8,
				EpisodeEnd:   9,
			},
		},
		{
			input: "Просмотрены 7-й, 8-й и 9-й эпизоды",
			output: WatchDescription{
				EpisodeStart: 7,
				EpisodeEnd:   9,
			},
		},
		{
			input: "Просмотрено 5 эпизодов",
			output: WatchDescription{
				EpisodeStart: 1,
				EpisodeEnd:   5,
			},
		},
		{
			input: "Просмотрены 4 эпизода",
			output: WatchDescription{
				EpisodeStart: 1,
				EpisodeEnd:   4,
			},
		},
		{
			input: "Просмотрены с 14-го по 18-й эпизоды",
			output: WatchDescription{
				EpisodeStart: 14,
				EpisodeEnd:   18,
			},
		},
	}

	for _, tC := range testCases {
		tC := tC

		t.Run(tC.input, func(t *testing.T) {
			res, err := Parser.ParseDescription(tC.input)
			require.NoError(t, err)
			require.Equal(t, &tC.output, res)
		})
	}
}
