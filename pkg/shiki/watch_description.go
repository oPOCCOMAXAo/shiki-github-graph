package shiki

import (
	"regexp"
	"strconv"

	"github.com/opoccomaxao-go/generic-collection/gmath"
)

// desc regexps.
var (
	descWatchRegexp    = regexp.MustCompile(`(?i)просмотр`)
	descFinishedRegexp = regexp.MustCompile(`(?i)^просмотрено( и оценено на <b>\d+</b>)?$`)
	descStartedRegexp  = regexp.MustCompile(`(?i)просмотрен(о|ы) (\d+) эпизод(ов|а)`)
	descMultiRegexp    = regexp.MustCompile(`(?i)просмотрены (\d+)[^ ]*?(?:, (\d+)[^ ]*?)? и (\d+)[^ ]*? эпизоды`)
	descSingleRegexp   = regexp.MustCompile(`(?i)просмотрен (\d+)[^ ]*? эпизод`)
	descRangeRegexp    = regexp.MustCompile(`(?i)просмотрены с (\d+).*? по (\d+).*? эпизоды`)
)

type WatchDescription struct {
	EpisodeStart int
	EpisodeEnd   int
	Finished     bool
	Failed       bool
}

func (desc *WatchDescription) ParseFinished(data string) (bool, error) {
	desc.Finished = descFinishedRegexp.MatchString(data)

	return desc.Finished, nil
}

func (desc *WatchDescription) ParseStarted(data string) (bool, error) {
	matches := descStartedRegexp.FindStringSubmatch(data)

	if len(matches) == 0 {
		return false, nil
	}

	desc.EpisodeStart = 1

	var err error

	desc.EpisodeEnd, err = strconv.Atoi(matches[2])
	if err != nil {
		return true, err
	}

	return true, nil
}

func (desc *WatchDescription) ParseSingle(data string) (bool, error) {
	matches := descSingleRegexp.FindStringSubmatch(data)

	if len(matches) == 0 {
		return false, nil
	}

	return true, desc.parseMinMaxNoZero(matches[1:2])
}

func (desc *WatchDescription) parseMinMaxNoZero(matches []string) error {
	episodes := make([]int, 0, len(matches))

	for _, match := range matches {
		if match == "" {
			continue
		}

		episode, err := strconv.Atoi(match)
		if err != nil {
			return err
		}

		episodes = append(episodes, episode)
	}

	desc.EpisodeStart = gmath.Min(episodes[0], episodes...)
	desc.EpisodeEnd = gmath.Max(episodes[0], episodes...)

	return nil
}

func (desc *WatchDescription) ParseMulti(data string) (bool, error) {
	matches := descMultiRegexp.FindStringSubmatch(data)

	if len(matches) == 0 {
		return false, nil
	}

	return true, desc.parseMinMaxNoZero(matches[1:4])
}

func (desc *WatchDescription) ParseRange(data string) (bool, error) {
	matches := descRangeRegexp.FindStringSubmatch(data)

	if len(matches) == 0 {
		return false, nil
	}

	return true, desc.parseMinMaxNoZero(matches[1:3])
}
