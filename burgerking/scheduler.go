package burgerking

import (
	"HelpingPixl/config"
	"HelpingPixl/models"
	"fmt"
	"github.com/disgoorg/disgo/webhook"
	"github.com/robfig/cron/v3"
	"log/slog"
	"strconv"
	"time"
)

func ScheduleDailyRefresh() error {

	slog.Info(fmt.Sprintf("(✓) Found %d coupon update hook(s)", len(config.Config.Discord.BKUpdateHookList)))

	c := cron.New()

	id, err := c.AddFunc("0 23 * * *", updateCoupons)
	if err != nil {
		return err
	}
	c.Start()
	c.Run()

	slog.Info(fmt.Sprintf("(✓) BurgerKing coupon scheduler started (next run in %s)", time.Until(c.Entry(id).Next).String()))

	return nil
}
func updateCoupons() {
	startTime := time.Now().UnixNano()
	coupons, _, _, err := Crawl()
	if err != nil {
		slog.Error("Error while fetching coupons: ", err.Error())
		return
	}

	go func() {
		postToHooks(&coupons, int((time.Now().UnixNano()-startTime)/1e6))

		SaveCoupons(coupons)
	}()
}

func postToHooks(coupons *[]models.Coupon, timeElapsed int) {
	messages := GetCouponUpdateEmbeds(coupons, &CachedCoupons.Coupons, timeElapsed)

	for i, hook := range config.Config.Discord.BKUpdateHookList {
		client, err := webhook.NewWithURL(hook)
		if err != nil {
			continue
		}

		for _, msg := range messages {
			_, err := client.CreateMessage(msg)
			if err != nil {
				slog.Error("Error while posting coupons to hook "+strconv.Itoa(i)+": ", err.Error())
				break
			}
			time.Sleep(time.Second)
		}

		slog.Info("Posted coupons to hook (" + strconv.Itoa(i) + ")")
	}
}
