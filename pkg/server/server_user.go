package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opoccomaxao-go/generic-collection/gmath"
	"github.com/opoccomaxao-go/generic-collection/slice"

	"github.com/opoccomaxao/shiki-github-graph/pkg/image"
	"github.com/opoccomaxao/shiki-github-graph/pkg/storage"
)

func (s *Server) GetUser(ctx *gin.Context) {
	nick := ctx.Param("nick")

	err := s.prepareDummy(nick)
	if err != nil {
		log.Printf("%+v\n", err)
	}

	ctx.Redirect(http.StatusFound, "/image/"+nick+".svg")

	go s.taskScheduleUserUpdate(nick)
}

func (s *Server) taskScheduleUserUpdate(nick string) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	user, err := s.Storage.GetUserByName(ctx, nick)
	if err != nil {
		log.Printf("%+v\n", err)

		return
	}

	if user == nil {
		shikiUser, err := s.Shiki.GetUserByNick(ctx, nick)
		if err != nil {
			log.Printf("%+v\n", err)

			return
		}

		user = &storage.User{
			ID:   shikiUser.ID,
			Name: shikiUser.Nickname,
		}
	}

	user.RequestedAt = time.Now().Unix()

	err = s.Storage.SaveUser(ctx, user)
	if err != nil {
		log.Printf("%+v\n", err)

		return
	}

	s.usersProcessNotify <- struct{}{}
}

func (s *Server) processUsersUpdate(ctx context.Context) {
	done := ctx.Done()

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := s.runSingleUsersUpdate(ctx)
			if err != nil {
				log.Printf("%+v\n", err)
			}
		case <-s.usersProcessNotify:
			err := s.runSingleUsersUpdate(ctx)
			if err != nil {
				log.Printf("%+v\n", err)
			}
		}
	}
}

func (s *Server) runSingleUsersUpdate(ctx context.Context) error {
	for {
		user, err := s.Storage.GetUserToUpdate(ctx)
		if err != nil {
			return err
		}

		if user == nil {
			return nil
		}

		err = s.taskUpdateUser(ctx, user)
		if err != nil {
			return err
		}

		s.animeProcessNotify <- struct{}{}
	}
}

func (s *Server) taskUpdateUser(ctx context.Context, user *storage.User) error {
	log.Printf("update user #%d %s\n", user.ID, user.Name)

	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	err := s.taskUpdateUserLoadHistory(ctx, user.ID)
	if err != nil {
		return err
	}

	err = s.taskUpdateUserFixHistory(ctx, user.ID)
	if err != nil {
		return err
	}

	log.Printf("build user image #%d %s\n", user.ID, user.Name)
	full, err := s.taskUpdateUserBuildCalendar(ctx, user)
	if err != nil {
		return err
	}

	err = s.Storage.SaveUser(ctx, &storage.User{
		ID:            user.ID,
		LastUpdatedAt: time.Now().Unix(),
		Full:          &full,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) taskUpdateUserLoadHistory(ctx context.Context, userID int64) error {
	for result := range s.Shiki.LoadHistory(ctx, userID) {
		if result.Error != nil {
			return result.Error
		}

		for _, msg := range result.Warning {
			s.warningLogger.Println(msg)
		}

		history := make([]*storage.History, len(result.Entries))

		for i, entry := range result.Entries {
			history[i] = &storage.History{
				ID:           entry.ID,
				CreatedAt:    entry.CreatedAt.Unix(),
				UserID:       userID,
				AnimeID:      int32(entry.Target.ID),
				EpisodeStart: int16(entry.Watch.EpisodeStart),
				EpisodeEnd:   int16(entry.Watch.EpisodeEnd),
			}
		}

		err := s.Storage.SaveHistory(ctx, history)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) taskUpdateUserFixHistory(
	ctx context.Context,
	userID int64,
) error {
	history, err := s.Storage.GetHistoryBroken(ctx, userID)
	if err != nil {
		return err
	}

	animesID := make([]int32, len(history))
	for i, hist := range history {
		animesID[i] = hist.AnimeID
	}

	history, err = s.Storage.GetHistoryAnimes(ctx, userID, animesID)
	if err != nil {
		return err
	}

	toFix := []*storage.History{}
	maxEpisodeValid := map[int32]int16{}

	for _, hist := range history {
		maxValid := maxEpisodeValid[hist.AnimeID]

		if hist.EpisodeStart == 0 {
			hist.EpisodeStart = maxValid + 1
			toFix = append(toFix, hist)
		} else {
			maxEpisodeValid[hist.AnimeID] = gmath.Max(maxValid, hist.EpisodeEnd)
		}
	}

	err = s.Storage.SaveHistory(ctx, history)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) taskUpdateUserBuildCalendar(
	ctx context.Context,
	user *storage.User,
) (bool, error) {
	finished, err := s.Storage.IsCalendarUnfinished(ctx, user.ID)
	if err != nil {
		return finished, err
	}

	calendar, err := s.Storage.GetUserCalendar(ctx, user.ID)
	if err != nil {
		return finished, err
	}

	params := image.BuildParams{
		Data: slice.Map(calendar, func(cp *storage.CalendarPoint) image.CalendarData {
			return image.CalendarData{
				Day:     time.Unix(cp.Time, user.ID),
				Seconds: cp.Seconds,
			}
		}),
		End: time.Now(),
	}
	params.Start = params.End.AddDate(-1, 0, 0)

	svgBody, err := image.BuildSVG(params)
	if err != nil {
		return finished, err
	}

	file, err := os.Create(s.imagePath(user.Name))
	if err != nil {
		return finished, err
	}

	_, err = file.WriteString(svgBody)

	return finished, err
}
