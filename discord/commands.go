package discord

import (
	dc "github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
)

const (
	ClearCommand     = "clear"
	BSProfileCommand = "bs-profile"
	BSSnipeCommand   = "bs-snipe"
	BKCouponsCommand = "bk-coupons"
	BKRefreshCommand = "bk-refresh"
	BKUploadCommand  = "bk-upload"
)

var (
	ClearCommandBoundMin = 2
	ClearCommandBoundMax = 1<<15 - 1

	GlobalCommands = []dc.ApplicationCommandCreate{
		dc.SlashCommandCreate{
			Name:        ClearCommand,
			Description: "Clears the channel or a desired chunk of messages",
			Options: []dc.ApplicationCommandOption{
				dc.ApplicationCommandOptionInt{
					Name:        "count",
					Description: "The number of messages to clear",
					Required:    true,
					MaxValue:    &ClearCommandBoundMax,
					MinValue:    &ClearCommandBoundMin,
				},
			},
			DefaultMemberPermissions: json.NewNullablePtr[dc.Permissions](dc.PermissionAdministrator),
		},
		dc.SlashCommandCreate{
			Name:        BSProfileCommand,
			Description: "Displays the BeatLeader profile of the given player",
			Options: []dc.ApplicationCommandOption{
				dc.ApplicationCommandOptionString{
					Name:        "profile_id",
					Description: "The steam id of the player to query",
					Required:    true,
				},
			},
		},
		dc.SlashCommandCreate{
			Name:        BSSnipeCommand,
			Description: "Generates a playlist with songs that are easy to snipe",
			Options: []dc.ApplicationCommandOption{
				dc.ApplicationCommandOptionString{
					Name:         "player",
					Description:  "You! (IGN)",
					Required:     true,
					Autocomplete: true,
				},
				dc.ApplicationCommandOptionString{
					Name:         "target",
					Description:  "The player to snipe (IGN)",
					Required:     true,
					Autocomplete: true,
				},
				dc.ApplicationCommandOptionInt{
					Name:        "leaderboard",
					Description: "The leaderboard to use",
					Required:    true,
					Choices: []dc.ApplicationCommandOptionChoiceInt{
						{
							Name:              "ScoreSaber",
							NameLocalizations: nil,
							Value:             1,
						},
						{
							Name:              "BeatLeader",
							NameLocalizations: nil,
							Value:             0,
						},
					},
				},
			},
		},
		dc.SlashCommandCreate{
			Name:        BKCouponsCommand,
			Description: "Offers a list of current BurgerKing® Coupons",
		},
		dc.SlashCommandCreate{
			Name:                     BKRefreshCommand,
			Description:              "Updates the saved BurgerKing Coupons",
			DefaultMemberPermissions: json.NewNullablePtr(dc.PermissionAdministrator),
		},
		dc.SlashCommandCreate{
			Name:                     BKUploadCommand,
			Description:              "Adds paper coupons to the db",
			DefaultMemberPermissions: json.NewNullablePtr(dc.PermissionAdministrator),
			Options: []dc.ApplicationCommandOption{
				dc.ApplicationCommandOptionAttachment{
					Name:        "image",
					Description: "Paper Coupon Image",
					Required:    true,
				},
			},
		},
	}
)
