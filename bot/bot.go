package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/geekloper/discord-bot-ip-whitelister/config"
	"github.com/geekloper/discord-bot-ip-whitelister/logger"
)

var s *discordgo.Session

func InitBot() {
	var err error
	botToken := config.GetEnv("BOT_TOKEN", true)

	s, err = discordgo.New("Bot " + botToken)
	if err != nil {
		logger.Fatal("Invalid bot parameters: ", err)
	}

	s.AddHandler(HandleInteractions)
	s.AddHandler(HandleReady)
}

func OpenSession() error {
	return s.Open()
}

func CloseSession() {
	s.Close()
}

func RemoveCommands(botGuildID string) {
	logger.Info("Removing commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, botGuildID)
	if err != nil {
		logger.Fatal("Could not fetch registered commands: ", err)
	}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, botGuildID, v.ID)
		if err != nil {
			logger.Fatal("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func RegisterCommands(botGuildID string) {
	logger.Info("Adding commands...")
	for _, v := range Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, botGuildID, v)
		if err != nil {
			logger.Fatal("Cannot create '%v' command: %v", v.Name, err)
		}
	}
}

func HandleInteractions(session *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
		h(session, i)
	}
}

func HandleReady(session *discordgo.Session, r *discordgo.Ready) {
	logger.Info("Logged in as", "username", session.State.User.Username, "discriminator", session.State.User.Discriminator)
}
