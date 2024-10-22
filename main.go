package main

import (
	"HelpingPixl/burgerking"
	"HelpingPixl/config"
	"HelpingPixl/discord"
	"log/slog"
	"os"
)

func main() {
	//slog.SetLogLoggerLevel(slog.LevelDebug)

	config.Load()

	burgerking.Load()

	err := burgerking.ScheduleDailyRefresh()
	if err != nil {
		slog.Error("An error occurred while starting the coupon refresh scheduler. ", err)
		os.Exit(1)
		return
	}

	err = discord.Launch()
	if err != nil {
		slog.Error("An error occurred while launching discord, exiting. ", err)
		os.Exit(1)
		return
	}
}
