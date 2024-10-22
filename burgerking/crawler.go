package burgerking

import (
	"HelpingPixl/config"
	"HelpingPixl/models"
	"github.com/Jeffail/gabs/v2"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Crawl() (coupons []models.Coupon, additionalCount int, timeElapsedMs int64, err error) {

	slog.Info("Fetching BurgerKing coupons...")

	startTime := time.Now().UnixNano()

	coupons, err = fetchGraphQlCoupons(false)

	additionalCoupons, err := fetchGraphQlCoupons(true)

	if additionalCoupons != nil {
		coupons = slices.Concat(coupons, additionalCoupons)
	}

	sort.Slice(coupons, func(i, j int) bool {
		return coupons[i].OfferPrice < coupons[j].OfferPrice
	})

	return coupons, len(additionalCoupons), (time.Now().UnixNano() - startTime) / 1e6, nil
}

func fetchGraphQlCoupons(useAdditionalBindings bool) ([]models.Coupon, error) {
	client := resty.New()
	client.SetHeaders(Headers)

	var requestString = map[string]string{}

	if useAdditionalBindings {
		requestString["query"] = models.BKConfigOffersQuery
	} else {
		requestString["query"] = models.BKSystemwideOffersQuery
	}

	response, err := client.R().ForceContentType("application/json").SetBody(requestString).Post(models.BKGraphQlUrl)
	if err != nil || response.IsError() {
		return nil, err
	}

	responseJson, err := gabs.ParseJSON(response.Body())
	if err != nil {
		return nil, err
	}

	coupons := make([]models.Coupon, 0)

	var base []string

	if useAdditionalBindings {
		base = []string{"data", "allConfigOffers", "*"}
	} else {
		base = []string{"data", "allSystemwideOffers", "*"}
	}
	for _, child := range responseJson.S(base...).Children() {
		if useAdditionalBindings {
			newCoupon := buildCoupon(child)

			if newCoupon.ExpirationDate < time.Now().Unix() {
				continue
			}

			coupons = append(coupons, newCoupon)
		} else {
			newCoupon := buildCoupon(child)

			if newCoupon.ExpirationDate < time.Now().Unix() {
				continue
			}

			coupons = append(coupons, newCoupon)
			/*
				path := []string{"sortedSystemwideOffers", "*"}
				if child.Exists(path...) {
					for _, coupon := range child.S(path...).Children() {

					}
				}
			*/
		}
	}

	tmp := coupons[:0]
	for _, coupon := range coupons {
		if !((coupon.Plu == "" || coupon.Plu == "null") && (coupon.ConstantPlu == "" || coupon.ConstantPlu == "null")) {
			tmp = append(tmp, coupon)
		}
	}
	coupons = tmp

	return coupons, nil
}
func buildCoupon(coupon *gabs.Container) (newCoupon models.Coupon) {
	namePath := []string{"name", "deRaw", "0", "children", "0", "text"}
	idPath := []string{"_id"}
	imgDescPath := []string{"localizedImage", "de", "imageDescription"}
	descsPath := []string{"description", "deRaw"}
	engineIdPath := []string{"loyaltyEngineId"}
	offerPricePath := []string{"offerPrice"}
	shortCodePath := []string{"shortCode"}
	constantPluPath := []string{"vendorConfigs", "partner", "constantPlu"}
	moreInfoPath := []string{"moreInfo", "deRaw", "0", "children", "0", "text"}
	rulePath := []string{"rules"}
	couponTypePath := []string{"_type"}
	assetUrlPath := []string{"localizedImage", "de", "app", "asset", "url"}

	var dateSuccess bool

	//Title
	if coupon.Exists(namePath...) {
		newCoupon.Title = trimPS(coupon.S(namePath...).String())
	}
	//Id
	if coupon.Exists(idPath...) {
		newCoupon.Id = trimPS(coupon.S(idPath...).String())
	}
	//Image Description
	if coupon.Exists(imgDescPath...) && coupon.S(imgDescPath...).String() != "null" {
		str := trimPS(coupon.S(imgDescPath...).String())
		newCoupon.Description = str
	} else if coupon.Exists(descsPath...) {
		for _, child := range coupon.S(descsPath...).Children() {
			path := []string{"children", "0", "text"}
			if child.Exists(path...) {
				str := trimPS(child.S(path...).String())
				if strings.ContainsRune(str, '+') {
					newCoupon.Description = str
					break
				}
			}
		}

		if newCoupon.Description == "" {
			newCoupon.Description = newCoupon.Title
		}
	}
	//Asset URL (Image)
	if coupon.Exists(assetUrlPath...) {
		newCoupon.ImageUrl = trimPS(coupon.S(assetUrlPath...).String())
	}
	//Engine ID (TID for website)
	if coupon.Exists(engineIdPath...) {
		newCoupon.AddBrowserViewUrl(trimPS(coupon.S(engineIdPath...).String()))
	}
	//Price
	if coupon.Exists(offerPricePath...) {
		i, err := strconv.Atoi(coupon.S(offerPricePath...).String())
		if err != nil {
			i = -1
		}
		if i == 0 {
			if strings.HasPrefix(newCoupon.Title, "2") {
				newCoupon.Discount = 50
			}
		}
		newCoupon.OfferPrice = i
	}
	//Code (PLU)
	if coupon.Exists(shortCodePath...) {
		newCoupon.Plu = trimPS(coupon.S(shortCodePath...).String())
	}
	//Constant PLU
	if coupon.Exists(constantPluPath...) {
		newCoupon.ConstantPlu = trimPS(coupon.S(constantPluPath...).String())
	}
	//ExpirationDate
	if coupon.Exists(rulePath...) {
		if newCoupon.AddPrimaryExpirationDate(coupon.S(rulePath...).Children()) {
			dateSuccess = true
		}
	}
	if coupon.Exists(couponTypePath...) {
		if trimPS(coupon.S(couponTypePath...).String()) == "configOffer" {
			newCoupon.IsAdditional = true
		}
	}
	//ExpirationDate Backup
	if !dateSuccess {
		if coupon.Exists(moreInfoPath...) {
			newCoupon.AddSecondaryExpirationDate(trimPS(coupon.S(moreInfoPath...).String()))
		} else {
			newCoupon.Warning = config.Config.BurgerKing.NoExpirationDate
		}
	}

	return
}

var Headers = map[string]string{
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Origin":             "https://www.burgerking.de",
	"Content-Type":       "application/json",
	"sec-ch-ua":          "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"99\", \"Google Chrome\";v=\"99\"",
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": "\"Windows\"",
	"sec-fetch-dest":     "empty",
	"sec-fetch-mode":     "cors",
	"sec-fetch-site":     "cross-site",
	"x-ui-language":      "de",
	"x-ui-platform":      "web",
	"x-ui-region":        "DE",
}

func trimPS(str string) string {
	return strings.ReplaceAll(str, "\"", "")
}
