package discord

import (
	"HelpingPixl/config"
	"context"
	"github.com/disgoorg/disgo"
	bot2 "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var Bot bot2.Client
var AppID snowflake.ID

func Launch() error {
	bot, err := disgo.New(config.Config.Discord.BotToken, bot2.WithGatewayConfigOpts(
		gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildMessages, gateway.IntentMessageContent),
		gateway.WithPresenceOpts(gateway.WithPlayingActivity("BeatSaber")),
	),
		bot2.WithEventListenerFunc(OnInteractionCreate),
		bot2.WithEventListenerFunc(OnReady),
		bot2.WithEventListenerFunc(OnAutocomplete),
		bot2.WithEventListenerFunc(OnComponentInteract),
	)
	if err != nil {
		slog.Error("Error creating a discord-bot session.", slog.Any("err", err))
		return err
	}

	defer bot.Close(context.TODO())

	if _, err = bot.Rest().SetGlobalCommands(bot.ApplicationID(), GlobalCommands); err != nil {
		slog.Error("Error while registering commands", slog.Any("err", err))
		return err
	} else {
		slog.Info("(✓) Discord Bot Registered Commands")
	}

	if err = bot.OpenGateway(context.TODO()); err != nil {
		slog.Error("Error while connecting to gateway", slog.Any("err", err))
		return err
	}

	slog.Info("(✓) Discord Bot started!")

	AppID = bot.ApplicationID()
	Bot = bot

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s

	return nil
}
