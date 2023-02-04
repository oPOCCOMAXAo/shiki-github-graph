package shiki

type parser byte

type SelfParseFunc func(string) (bool, error)

const Parser = parser(0)

func (parser) ParseDescription(data string) (*WatchDescription, error) {
	if !descWatchRegexp.MatchString(data) {
		return nil, nil //nolint:nilnil
	}

	res := WatchDescription{}

	for _, parse := range []SelfParseFunc{
		res.ParseFinished,
		res.ParseStarted,
		res.ParseSingle,
		res.ParseMulti,
		res.ParseRange,
	} {
		ok, err := parse(data)
		if err != nil {
			return nil, err
		}

		if ok {
			return &res, nil
		}
	}

	res.Failed = true

	return &res, nil
}
