package main

import (
	"HelpingPixl/burgerking"
	"HelpingPixl/config"
	"HelpingPixl/discord"
	"HelpingPixl/networking"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"syscall"
)

func main() {
	if slices.Contains(os.Args, "--debug") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	config.Load()

	burgerking.Load()

	err := burgerking.ScheduleDailyRefresh()
	if err != nil {
		slog.Error("An error occurred while starting the coupon refresh scheduler. ", err)
		os.Exit(1)
		return
	}

	go func() {
		err = discord.Launch()
		if err != nil {
			slog.Error("An error occurred while launching discord, exiting. ", err)
			os.Exit(1)
			return
		}
	}()

	go func() {
		gin.SetMode(gin.ReleaseMode)
		if config.Config.WebServerAPI.EnableAPI {
			r := networking.SetupRouter()
			err = r.Run("0.0.0.0:" + config.Config.WebServerAPI.APIPort)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-signalChannel

	slog.Info("Exiting program...")
}
