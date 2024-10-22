package burgerking

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"log/slog"
	"time"
)

func ScheduleDailyRefresh() error {

	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}

	j, err := s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(
			func() {
				coupons, _, _, err := Crawl()
				if err != nil {
					slog.Error("Error while fetching coupons: ", err.Error())
					return
				}
				SaveCoupons(coupons)
			},
		),
	)
	if err != nil {
		return err
	}
	s.Start()

	nextRun, err := j.NextRun()
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("(✓) BurgerKing coupon scheduler started (next run at %s)", nextRun.Format(time.DateTime)))

	return nil
}
