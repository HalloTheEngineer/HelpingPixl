package networking

import (
	"HelpingPixl/beatsaber"
	"HelpingPixl/config"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

const (
	FilesPath = "./server/files"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	_ = os.MkdirAll(FilesPath, 0755)

	r.GET("/", handleAPIRoot)

	if config.Config.WebServerAPI.EnableAPI {

		fileServer := http.FileServer(http.Dir(FilesPath))

		r.GET("/files/*filepath", func(c *gin.Context) {
			http.StripPrefix("/files", fileServer).ServeHTTP(c.Writer, c.Request)
		})
	}

	bk := r.Group("/burgerking")

	bk.GET("/coupons", handleBKCoupons)

	bs := r.Group("/beatsaber")

	bs.GET("/playlist/:key", handleBSPlaylist)

	r.NoRoute(noRoute)

	slog.Info("(✓) API Engine starting")
	return r
}

func handleBSPlaylist(c *gin.Context) {
	s := c.Param("key")

	switch s {
	case "snipe":
		val := c.Request.URL.Query()

		if !val.Has("self") || !val.Has("target") || !val.Has("leaderboard") {
			c.Status(http.StatusBadRequest)
			return
		}

		self := val.Get("self")
		target := val.Get("target")
		leaderboard := val.Get("leaderboard")

		lInt, err := strconv.Atoi(leaderboard)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		snipePl, _, errStr := beatsaber.SnipeHoldPlaylist(nil, nil, &self, &target, lInt)
		if errStr != "" {
			c.Status(http.StatusInternalServerError)
			return
		}
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		data, err := json.MarshalIndent(snipePl, "", "   ")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "application/json", data)

	case "hold":
		val := c.Request.URL.Query()

		if !val.Has("self") || !val.Has("target") || !val.Has("leaderboard") {
			c.Status(http.StatusBadRequest)
			return
		}

		self := val.Get("self")
		target := val.Get("target")
		leaderboard := val.Get("leaderboard")

		lInt, err := strconv.Atoi(leaderboard)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		_, holdPl, errStr := beatsaber.SnipeHoldPlaylist(nil, nil, &self, &target, lInt)
		if errStr != "" {
			c.Status(http.StatusInternalServerError)
			return
		}
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		data, err := json.MarshalIndent(holdPl, "", "  ")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "application/json", data)
	}
}

func handleBKCoupons(c *gin.Context) {

}

func handleAPIRoot(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

func noRoute(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "https://bit.ly/3BlS71b")
}
