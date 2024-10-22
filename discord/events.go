package discord

import (
	"HelpingPixl/burgerking"
	"HelpingPixl/config"
	"HelpingPixl/discord/modules/beatsaber"
	bk "HelpingPixl/discord/modules/burgerking"
	"HelpingPixl/models"
	"HelpingPixl/utils"
	"bytes"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/go-resty/resty/v2"
	_ "github.com/json-iterator/go"
	jsoniter "github.com/json-iterator/go"
	"image"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var ongoingProcess = make(map[string]bool)

func OnReady(e *events.Ready) {
	slog.Info(fmt.Sprintf("(✓) Discord Bot Logged In (%s, %d guild(s))", e.User.Username, len(e.Guilds)))
}

func OnComponentInteract(e *events.ComponentInteractionCreate) {
	id := e.Data.CustomID()
	if strings.HasPrefix(id, "coupon-chooser-") {
		data := e.StringSelectMenuInteractionData()

		couponId := data.Values[0]

		coupon := burgerking.CachedCoupons.GetById(couponId)
		if coupon != nil {
			_ = e.CreateMessage(bk.BuildCouponMsg(coupon, e.Message.ID))
		} else {
			_ = e.CreateMessage(GetErrorEmbed("This coupon could not be found! Is it expired or from a past day?", true))
		}
	}
}

func OnAutocomplete(e *events.AutocompleteInteractionCreate) {
	switch e.Data.CommandName {
	case BSSnipeCommand:
		go func() {
			if err := e.AutocompleteResult(buildPlayerAutocomplete(e)); err != nil {
				slog.Error(err.Error())
			}
		}()
	}
}

func OnInteractionCreate(e *events.ApplicationCommandInteractionCreate) {
	data := e.SlashCommandInteractionData()

	switch data.CommandName() {
	case ClearCommand:
		count, existing := data.OptInt("count")
		if !existing {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.InternalError, true))
			return
		}
		go func() {
			err := clearMessages(count, e)
			if err != nil {
				_, _ = Bot.Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), GetUpdateErrorEmbed(config.Config.Formatting.InternalError+"\n```"+err.Error()+"```"))
				return
			}
		}()
	case BSProfileCommand:
		id, existing := data.OptString("profile_id")
		if !existing {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.NotEnoughArguments, true))
			return
		}
		_ = e.DeferCreateMessage(false)
		get, err := http.Get(fmt.Sprintf("https://render.beatleader.xyz/screenshot/800x600/myprofile/general/u/%s", id))
		if err != nil {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.InternalError, true))
			return
		}
		_, _ = Bot.Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.NewMessageUpdateBuilder().AddFiles(&discord.File{
			Name:   id + ".png",
			Reader: get.Body,
		}).Build())
		_ = get.Body.Close()
	case BSSnipeCommand:
		timeStart := time.Now().UnixNano()

		self, existing := data.OptString("player")
		if !existing {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.NotEnoughArguments, true))
			return
		}
		target, existing := data.OptString("target")
		if !existing {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.NotEnoughArguments, true))
			return
		}

		leaderboard, existing := data.OptInt("leaderboard")
		if !existing {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.NotEnoughArguments, true))
			return
		}

		if _, ok := ongoingProcess[e.User().ID.String()]; ok {
			_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.AlreadySniping, true))
			return
		}

		go func() {
			ongoingProcess[e.User().ID.String()] = true

			snipe, hold, errStr, isDelayed := beatsaber.SnipeHoldPlaylist(e, self, target, leaderboard)
			if errStr != "" {
				if isDelayed {
					_, _ = Bot.Rest().UpdateInteractionResponse(AppID, e.Token(), GetUpdateErrorEmbed(config.Config.Formatting.InternalError))
				}
				delete(ongoingProcess, e.User().ID.String())
				_ = e.CreateMessage(GetErrorEmbed(errStr, true))
				return
			}
			json := jsoniter.ConfigCompatibleWithStandardLibrary

			byteSnipePl, err := json.MarshalIndent(snipe, "", "   ")
			if err != nil {
				delete(ongoingProcess, e.User().ID.String())
				_, _ = Bot.Rest().UpdateInteractionResponse(AppID, e.Token(), GetUpdateErrorEmbed(config.Config.Formatting.InternalError))
				slog.Error("Error while marshaling the snipe playlist")
				return
			}
			byteHoldPl, err := json.MarshalIndent(hold, "", "   ")
			if err != nil {
				delete(ongoingProcess, e.User().ID.String())
				_, _ = Bot.Rest().UpdateInteractionResponse(AppID, e.Token(), GetUpdateErrorEmbed(config.Config.Formatting.InternalError))
				slog.Error("Error while marshaling the hold playlist")
				return
			}
			timeEnd := time.Now().UnixNano()

			files := []*discord.File{
				discord.NewFile(fmt.Sprintf("snipe_%s.json", snipe.Stats.TargetName), fmt.Sprintf(config.Config.BeatSaber.SnipeFileDescription, snipe.Stats.TargetName, snipe.Stats.SelfName), bytes.NewReader(byteSnipePl)),
				discord.NewFile(fmt.Sprintf("hold_%s.json", snipe.Stats.TargetName), fmt.Sprintf(config.Config.BeatSaber.HoldFileDescription, snipe.Stats.TargetName, snipe.Stats.SelfName), bytes.NewReader(byteHoldPl)),
			}

			_, _ = Bot.Rest().UpdateInteractionResponse(AppID, e.Token(), GetUpdateSuccessFileEmbed(fmt.Sprintf(config.Config.Formatting.PlaylistMsg, snipe.Stats.SelfConsidered, snipe.Stats.SelfName, snipe.Stats.TargetConsidered, snipe.Stats.TargetName, snipe.Stats.SnipeCount, snipe.Stats.HoldCount, (timeEnd-timeStart)/1e6), files...))
			delete(ongoingProcess, e.User().ID.String())
		}()
	case BKCouponsCommand:
		_ = e.CreateMessage(bk.BuildCouponCompMsg(&burgerking.CachedCoupons))
	case BKRefreshCommand:
		collectedCoupons, addCount, ms, err := burgerking.Crawl()
		if err != nil {
			slog.Error("Error while crawling Coupons: ", err.Error())
			break
		}
		burgerking.SaveCoupons(collectedCoupons)
		_ = e.CreateMessage(GetSuccessEmbed(fmt.Sprintf("**%d** coupons were found while fetching, **%d** of them were additional!\nThe process took **%dms**.", len(collectedCoupons), addCount, ms)))
	case BKUploadCommand:
		file, ok := e.SlashCommandInteractionData().OptAttachment("image")
		if ok {
			if strings.HasPrefix(*file.ContentType, "image") {
				startTime := time.Now().UnixNano()

				_ = e.DeferCreateMessage(true)

				get, err := resty.New().R().Get(file.URL)
				if err != nil {
					_ = e.CreateMessage(GetErrorEmbed(config.Config.Formatting.InternalError, true))
					return
				}

				img, _, err := image.Decode(bytes.NewReader(get.Body()))

				codes := burgerking.FindQRCodes(img)

				if len(codes) > 0 {
					_, _ = Bot.Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), GetUpdateSuccessEmbed(fmt.Sprintf("Found %d coupons!\nTook **%dms**", len(codes), (time.Now().UnixNano()-startTime)/1e6)))
				} else {
					_, _ = Bot.Rest().UpdateInteractionResponse(e.ApplicationID(), e.Token(), GetUpdateErrorEmbed("No coupons found!"))
				}
			}

		}
	}

}

func clearMessages(count int, e *events.ApplicationCommandInteractionCreate) error {
	_ = e.CreateMessage(GetDeferEmbed("Deleting messages...", true))

	for count > 0 {
		toClear := 100
		if !(count > 100) {
			toClear = count
		}
		count -= 100

		messages, err := Bot.Rest().GetMessages(e.Channel().ID(), 0, 0, 0, toClear)
		if err != nil {
			return err
		}
		err = Bot.Rest().BulkDeleteMessages(e.Channel().ID(), utils.Map[discord.Message, snowflake.ID](messages, func(t discord.Message) snowflake.ID { return t.ID }))
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	_, _ = Bot.Rest().UpdateInteractionResponse(e.ApplicationID(), e.ApplicationCommandInteraction.Token(), GetUpdateSuccessEmbed("Messages have been deleted successfully!"))

	return nil
}
func buildPlayerAutocomplete(e *events.AutocompleteInteractionCreate) (choices []discord.AutocompleteChoice) {
	var query string
	focused := e.AutocompleteInteraction.Data.Focused()
	if (e.Data.CommandName == BSProfileCommand && focused.Name == "profile_id") || (e.Data.CommandName == BSSnipeCommand && (focused.Name == "player" || focused.Name == "target")) {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		data, err := json.Marshal(focused.Value)
		if err != nil {
			return []discord.AutocompleteChoice{}
		}
		query = strings.Trim(string(data), "\"")
	} else {
		return []discord.AutocompleteChoice{
			discord.AutocompleteChoiceInt{
				Name:  "Error",
				Value: 0,
			},
		}
	}
	slog.Info(query)

	if !(len(query) > 3) {
		return []discord.AutocompleteChoice{}
	}

	players, err := utils.FetchToStruct[models.BLPlayersResponse](fmt.Sprintf(models.BLBase+models.BLPlayers, query))
	if err != nil {
		return []discord.AutocompleteChoice{}
	}

	for _, resp := range players.Data {
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  resp.Name,
			Value: resp.Name,
		})
	}
	return
}
