package discord

import (
	"github.com/disgoorg/disgo/discord"
)

func GetErrorEmbed(description string, ephemeral bool) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("❌ Error").SetColor(16711680).SetDescription(description).Build()).SetEphemeral(true).SetEphemeral(ephemeral).Build()
}
func GetSuccessEmbed(description string) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("✅ Success").SetColor(65280).SetDescription(description).Build()).SetEphemeral(true).Build()
}
func GetSuccessFileEmbed(description string, ephemeral bool, files ...*discord.File) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("✅ Success").SetColor(65280).SetDescription(description).Build()).SetEphemeral(ephemeral).AddFiles(files...).Build()
}
func GetDeferEmbed(description string, ephemeral bool) discord.MessageCreate {
	return discord.NewMessageCreateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("⏲️ Ongoing Action").SetColor(16776960).SetDescription(description).Build()).SetEphemeral(ephemeral).Build()
}
func GetUpdateSuccessEmbed(description string) discord.MessageUpdate {
	return discord.NewMessageUpdateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("✅ Success").SetColor(65280).SetDescription(description).Build()).Build()
}
func GetUpdateSuccessFileEmbed(description string, files ...*discord.File) discord.MessageUpdate {
	return discord.NewMessageUpdateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("✅ Success").SetColor(65280).SetDescription(description).Build()).AddFiles(files...).Build()
}
func GetUpdateErrorEmbed(description string) discord.MessageUpdate {
	return discord.NewMessageUpdateBuilder().SetEmbeds(discord.NewEmbedBuilder().SetTitle("❌ Error").SetColor(16711680).SetDescription(description).Build()).Build()
}
