package models

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	jsoniter "github.com/json-iterator/go"
	"log/slog"
	"regexp"
	"strings"
	"time"
)

var dateRegex = regexp.MustCompile(`(?m)(?i)Abgabe bis (\d{1,2}\.\d{1,2}\.\d{4})`)

type (
	Coupon struct {
		Title          string `json:"title"`
		Id             string `json:"id"`
		Description    string `json:"description"`
		Plu            string `json:"plu"`
		ConstantPlu    string `json:"constant_plu"`
		ImageUrl       string `json:"image_url"`
		WebViewUrl     string `json:"web_view_url"`
		OfferPrice     int    `json:"offer_price"`
		Discount       int    `json:"discount"`
		StartDate      int64  `json:"start_date"`
		ExpirationDate int64  `json:"expiration_date"`
		Warning        string `json:"warning"`
		IsAdditional   bool   `json:"is_additional"`
	}
	CouponCache struct {
		Date    int64    `json:"date"`
		Count   int      `json:"count"`
		Coupons []Coupon `json:"coupons"`
	}
)

func (c CouponCache) ToJsonString() string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	str, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return "{}"
	}
	return string(str)
}

func (c CouponCache) GetById(id string) *Coupon {
	for _, coupon := range c.Coupons {
		if coupon.Id == id {
			return &coupon
		}
	}
	return nil
}

func (c *Coupon) ToString() string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	str, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		return "Error"
	}
	return string(str)
}
func (c *Coupon) AddImageUrl(assetId string) {
	c.ImageUrl = fmt.Sprintf(BKCouponImageUrl, assetId+".png")
}
func (c *Coupon) AddBrowserViewUrl(engineId string) {
	c.WebViewUrl = fmt.Sprintf(BKCouponWebViewUrl, engineId)
}

func (c *Coupon) AddSecondaryExpirationDate(infoText string) bool {
	matches := dateRegex.FindAllString(infoText, -1)
	if matches == nil {
		return false
	}
	datetime, err := time.Parse("02.01.2006 15:04:05", strings.Split(matches[0], " ")[2]+" 23:59:59")
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	c.ExpirationDate = datetime.Unix()
	return true
}
func (c *Coupon) AddPrimaryExpirationDate(rules []*gabs.Container) bool {
	typeName := "_type"
	var startDatetime = time.Unix(0, 0)
	var endDatetime = time.Unix(0, 0)
	for _, rule := range rules {
		if rule.Exists(typeName) {
			if trimPS(rule.S(typeName).String()) == "loyalty-between-dates" {
				sd := rule.Path("startDate")
				ed := rule.Path("endDate")
				var err error
				if sd != nil {
					startDatetime, err = time.Parse("2006-01-02 15:04:05", trimPS(sd.String()+" 23:59:59"))
				}
				if ed != nil {
					endDatetime, err = time.Parse("2006-01-02 15:04:05", trimPS(ed.String()+" 23:59:59"))
				}
				if err != nil {
					slog.Error(err.Error())
				}
				break
			}
		}
	}

	if startDatetime.IsZero() && endDatetime.IsZero() {
		return false
	}

	c.StartDate = startDatetime.Unix()
	c.ExpirationDate = endDatetime.Unix()
	return true
}
func trimPS(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}
