package burgerking

import (
	"github.com/makiuchi-d/gozxing"
	qrcode2 "github.com/makiuchi-d/gozxing/multi/qrcode"
	"image"
	"image/draw"
)

func FindQRCodes(img image.Image) (codes []*gozxing.Result) {
	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, img.Bounds(), img, image.Point{}, draw.Src)

	src := gozxing.NewLuminanceSourceFromImage(rgbaImg)
	bitmap, _ := gozxing.NewBinaryBitmap(gozxing.NewGlobalHistgramBinarizer(src))
	qrReader := qrcode2.NewQRCodeMultiReader()

	codes, _ = qrReader.DecodeMultiple(bitmap, map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	})

	return
}
