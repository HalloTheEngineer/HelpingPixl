package burgerking

import (
	"HelpingPixl/config"
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
	"strings"
	"time"
)

const (
	maxChars   = 6000
	qrCodeSize = 512
)

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
	eBuilder.AddField("Product", formatDescription(coupon.Title, coupon.Description), true)
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

func formatDescription(title, description string) string {
	description = strings.TrimPrefix(description, " ")
	if strings.HasPrefix(description, "+") {
		return title + " " + description
	}
	return description
}

func BuildCouponCompMsg(cache *models.CouponCache) discord.MessageCreate {
	b := discord.NewMessageCreateBuilder()
	eBuilder := discord.NewEmbedBuilder()

	eBuilder.SetTitle("🍔 BurgerKing® Coupons")
	eBuilder.SetColor(16750848)
	eBuilder.SetFooter("by @hallotheengineer", "https://cdn.discordapp.com/avatars/592779824519446538/b3992968a0bce170a4ac0b22e40fa97e.webp?size=40")
	eBuilder.SetDescription(config.Config.Formatting.BKCouponInfoDesc)

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
func GetCouponUpdateEmbeds(coupons *[]models.Coupon, oldCoupons *[]models.Coupon, timeElapsed int) (messages []discord.WebhookMessageCreate) {
	b := discord.NewWebhookMessageCreateBuilder()

	eBu := discord.NewEmbedBuilder()
	eBu.SetTitle("🍔 BurgerKing Coupons")
	eBu.SetDescription(fmt.Sprintf(config.Config.Formatting.BKUpdateDesc, len(*coupons)))
	eBu.SetFooter(fmt.Sprintf("by @hallotheengineer | %dms", timeElapsed), "https://cdn.discordapp.com/avatars/592779824519446538/b3992968a0bce170a4ac0b22e40fa97e.webp?size=40")
	eBu.AddField("Useful links", "> [SparKings](https://burgerking.de/sparkings)\n> [WebView](https://www.burgerking.de/rewards/offers)\n> [KingFinder](https://burgerking.de/store-locator)", false)
	eBu.SetColor(16750848)

	var embedList []discord.Embed

	embedList = append(embedList, eBu.Build())

	//Comparison / New Coupons
	var newCoupons []models.Coupon
	for _, newCoupon := range *coupons {
		var existing bool
		for _, oldCoupon := range *oldCoupons {
			if oldCoupon.Id == newCoupon.Id {
				existing = true
			}
		}

		if !existing {
			newCoupons = append(newCoupons, newCoupon)
		}
	}
	newCEmbed := discord.NewEmbedBuilder()
	newCEmbed.SetTitle("New Coupons (**" + strconv.Itoa(len(newCoupons)) + "**)")
	for _, coupon := range newCoupons {
		newCEmbed.Description += "> Coupon #" + getCode(&coupon) + ": " + coupon.Title + "\n"
	}
	embedList = append(embedList, newCEmbed.Build())

	//Chunking
	chunks := utils.ChunkBy[models.Coupon](*coupons, 18)
	for i, chunk := range chunks {
		eBuilder := discord.NewEmbedBuilder()
		eBuilder.SetTitle(fmt.Sprintf("Current Coupons (%d/%d)", i+1, len(chunks)))
		for _, coupon := range chunk {
			eBuilder.AddField(coupon.Title, fmt.Sprintf("> **Description**: %s\n> **Price**: %s\n> **Code**: %s", coupon.Description, formatPrice(coupon.OfferPrice, coupon.Discount), getCode(&coupon)), false)
		}
		eBuilder.SetFooter("by @hallotheengineer", "https://cdn.discordapp.com/avatars/592779824519446538/b3992968a0bce170a4ac0b22e40fa97e.webp?size=40")
		embedList = append(embedList, eBuilder.Build())
	}

	var currentBatch []discord.Embed
	currentChars := 0

	for _, embed := range embedList {
		embedChars := calculateCharacterCount(embed)
		if currentChars+embedChars > maxChars {
			b.AddEmbeds(currentBatch...)
			messages = append(messages, b.Build())
			currentBatch = []discord.Embed{}
			currentChars = 0
		}

		currentBatch = append(currentBatch, embed)
		currentChars += embedChars
	}

	return
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
func getCode(coupon *models.Coupon) (str string) {
	str = coupon.Plu
	if coupon.ConstantPlu != "" {
		str += " | " + coupon.ConstantPlu
	}
	return
}
func calculateCharacterCount(embed discord.Embed) int {
	totalChars := 0
	totalChars += len(embed.Title) + len(embed.Description)
	for _, field := range embed.Fields {
		totalChars += len(field.Name) + len(field.Value)
	}
	if embed.Footer != nil {
		totalChars += len(embed.Footer.Text)
	}
	if embed.Author != nil {
		totalChars += len(embed.Author.Name)
	}
	return totalChars
}
