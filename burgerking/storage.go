package burgerking

import (
	"HelpingPixl/models"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

var CachedCoupons models.CouponCache

func Load() {
	now := time.Now()

	todayStr := now.Format("2006-01-02")

	_, err := os.Lstat("coupons")
	if err != nil {
		_ = os.Mkdir("coupons", 0755)
		slog.Info("(✓) Coupons folder created")
	}

	dir, _ := os.ReadDir("coupons")

	if dir != nil {
		for _, info := range dir {
			if strings.HasPrefix(info.Name(), todayStr) && !info.IsDir() {

				file, _ := os.OpenFile("coupons/"+info.Name(), os.O_RDONLY, 0755)

				json := jsoniter.ConfigCompatibleWithStandardLibrary

				bytes, err := io.ReadAll(file)
				if err != nil {
					_ = os.Remove("coupons/" + info.Name())
					return
				}

				_ = json.Unmarshal(bytes, &CachedCoupons)

				slog.Info(fmt.Sprintf("(✓) Loaded %d BurgerKing coupons from file \"%s\"", len(CachedCoupons.Coupons), info.Name()))
				return
			}
		}
		fetchedCoupons, _, ms, err := Crawl()
		if err != nil {
			slog.Error("Fetching coupons failed: ", err)
			return
		}
		SaveCoupons(fetchedCoupons)
		slog.Info(fmt.Sprintf("(✓) Fetched BurgerKing coupons successfully (%dms)", ms))
	}
}

func SaveCoupons(coupons []models.Coupon) {
	cache := models.CouponCache{
		Date:    time.Now().Unix(),
		Count:   len(coupons),
		Coupons: coupons,
	}
	CachedCoupons = cache

	_ = os.Mkdir("coupons", 0755)
	now := time.Now()

	file, _ := os.OpenFile(fmt.Sprintf("coupons/%s.json", now.Format("2006-01-02")), os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	bytes, _ := json.MarshalIndent(&cache, "", "   ")
	_, _ = file.Write(bytes)
}
