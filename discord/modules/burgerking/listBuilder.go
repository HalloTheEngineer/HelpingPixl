package burgerking

import (
	"HelpingPixl/models"
	"HelpingPixl/utils"
	"bytes"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image/jpeg"
	"strconv"
	"time"
)

var qrCodeSize = 512

func BuildCouponMsg(coupon *models.Coupon, ref snowflake.ID) discord.MessageCreate {
	startTime := time.Now().UnixNano()

	b := discord.NewMessageCreateBuilder()
	eBuilder := discord.NewEmbedBuilder()
	eBuilder2 := discord.NewEmbedBuilder()

	if coupon.Warning != "" {
		eBuilder.AddField("Warning", coupon.Warning, false)
	}

	//QR Generator
	var rawCode string
	if coupon.Plu != "" && coupon.Plu != "null" {
		rawCode = coupon.Plu
	} else {
		rawCode = coupon.ConstantPlu
	}

	enc := qrcode.NewQRCodeWriter()
	img, _ := enc.Encode(rawCode, gozxing.BarcodeFormat_QR_CODE, qrCodeSize, qrCodeSize, nil)

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, img, &jpeg.Options{Quality: 80})

	b.AddFiles(&discord.File{
		Name:   coupon.Id + ".jpeg",
		Reader: buf,
	})

	eBuilder2.SetImage(coupon.ImageUrl)
	eBuilder.SetTitle("Coupon 📄")
	eBuilder.SetImage(coupon.ImageUrl)
	eBuilder.SetFooter(fmt.Sprintf("by @hallotheengineer | %dms", (time.Now().UnixNano()-startTime)/1e6), "https://cdn.discordapp.com/avatars/592779824519446538/b3992968a0bce170a4ac0b22e40fa97e.webp?size=40")
	eBuilder.AddField("Product", coupon.Description, true)
	eBuilder.AddField("Price Change", formatPrice(coupon.OfferPrice, coupon.Discount), true)
	eBuilder.AddField("Code", getCode(coupon), false)
	eBuilder.AddField("Validity", formatStartEndDate(coupon.StartDate, coupon.ExpirationDate), true)
	eBuilder.AddField("Useful Links", fmt.Sprintf("> [WebView](%s)", coupon.WebViewUrl), true)
	eBuilder.SetColor(16750848)

	b.AddEmbeds(eBuilder.Build())

	b.SetMessageReferenceByID(ref)
	b.SetEphemeral(true)

	return b.Build()
}

func getCode(coupon *models.Coupon) (str string) {
	str = coupon.Plu
	if coupon.ConstantPlu != "" {
		str += " | " + coupon.ConstantPlu
	}
	return
}
func BuildCouponCompMsg(cache *models.CouponCache) discord.MessageCreate {
	b := discord.NewMessageCreateBuilder()
	eBuilder := discord.NewEmbedBuilder()

	eBuilder.SetTitle("🍔 BurgerKing® Coupons")
	eBuilder.SetColor(16750848)
	eBuilder.SetFooter("by @hallotheengineer", "https://cdn.discordapp.com/avatars/592779824519446538/b3992968a0bce170a4ac0b22e40fa97e.webp?size=40")
	eBuilder.SetDescription("Hello there!\nAll current **[BurgerKing®](https://www.burgerking.de/)** coupons are fetched every day at midnight.\nTo view the coupons, use the dropdown menu below!\nBon appetit!")

	b.AddEmbeds(eBuilder.Build())

	var coupons []discord.StringSelectMenuOption

	for _, coupon := range cache.Coupons {
		coupons = append(coupons, discord.StringSelectMenuOption{
			Label:       formatDropdownString(&coupon),
			Value:       coupon.Id,
			Description: coupon.Description,
		})
	}

	for i, cChunk := range utils.ChunkBy[discord.StringSelectMenuOption](coupons, 25) {
		if i > 4 {
			continue
		}
		dropdown := discord.NewStringSelectMenu("coupon-chooser-"+strconv.Itoa(i), fmt.Sprintf("View coupons (%d) ...", i+1), cChunk...)
		b.AddActionRow(dropdown)
	}

	return b.Build()
}

func formatStartEndDate(start int64, end int64) (str string) {
	if start != 0 {
		str += fmt.Sprintf("> Starting: <t:%d:R>\n", start)
	}
	if end != 0 {
		str += fmt.Sprintf("> Ending: <t:%d:R>", end)
	}
	if str == "" {
		str = "Unknown"
	}
	return
}
func formatDropdownString(coupon *models.Coupon) string {
	priceStr := formatPrice(coupon.OfferPrice, coupon.Discount)

	return fmt.Sprintf("%s (%s)", coupon.Title, priceStr)
}
func formatPrice(price int, discount int) string {
	var priceStr string

	if discount == 0 {
		priceStr = strconv.FormatFloat(float64(price)/100.0, 'f', -1, 64) + "€"
	} else {
		priceStr = "-" + strconv.Itoa(discount) + "%"
	}
	return priceStr
}
