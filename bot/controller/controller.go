package controller

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/TrevorEdris/hypernova-bot/bot/services"
)

type (
	// Controller has access to all services via services.Container.
	// It is responsbile for the core business logic of utilizing all
	// services in conjunction with one another.
	Controller struct {
		Container *services.Container
	}
)

// NewController creates an instance of a Controller, given a reference to a Container.
func NewController(c *services.Container) Controller {
	return Controller{
		Container: c,
	}
}

// Run registers all the discord command handlers, opens a websocket connection,
// and begins execution.
func (c *Controller) Run(errChan chan error) {
	c.Container.Logger.Info("Running bot")

	// Register all commands and corresponding handlers
	err := c.prepare()
	if err != nil {
		errChan <- err
		return
	}

	// Open a websocket connection to the Discord API
	err = c.Container.DiscordSession.Open()
	if err != nil {
		errChan <- err
		return
	}

	// TODO: Make it easier to extract this state info?
	c.Container.Logger.Info("Successfully opened discord session",
		zap.String("session_id", c.Container.DiscordSession.State.SessionID),
		zap.String("bot_username", c.Container.DiscordSession.State.User.Username+"#"+c.Container.DiscordSession.State.User.Discriminator),
	)
}

// prepare registers all the event handler functions.
func (c *Controller) prepare() error {
	c.Container.DiscordSession.AddHandler(c.onMessageCreate)
	c.Container.DiscordSession.Identify.Intents = discordgo.IntentsGuildMessages
	return nil
}

// onMessageCreate is invoked by the websocket connection belonging to c.Container.DiscordSession.
// When the corresponding event is triggered, the websocket handler will invoke this function.
//
// Note: This is not a "slash command", but merely a function invoked on every single message
//       sent in the discord channels the bot has access to. Because of this, it will be more
//       computationally intensive than using a specific event trigger, such as a "slash command".
//       Modifying this to be a "slash command" will be a future improvement.
func (c *Controller) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	c.Container.Logger.Info("Handling new message event", zap.String("author", m.Author.Username+"#"+m.Author.Discriminator))
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
