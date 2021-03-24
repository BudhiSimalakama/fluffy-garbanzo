package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// MsgChannel channel to direct messages
type MsgChannel struct {
	msg *discordgo.Message
}

// Discord bot holder
type Discord struct {
	msgChan chan MsgChannel

	Session  *discordgo.Session
	Sessions []*discordgo.Session
}

func (d *Discord) onMsgCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Content == "" || (msg.Author != nil && msg.Author.Bot) || !strings.HasPrefix(msg.Content, ".notif ") {
		return
	}

	guild, err := s.Guild(msg.GuildID)
	if err != nil {
		return
	}

	if msg.Author.ID != guild.OwnerID {
		return
	}

	cmd := strings.Replace(msg.Content, ".notif ", "", 1)
	data := strings.Split(cmd, " ")
	switch data[0] {
	case "setchannel":
		AddToDB(msg.GuildID, msg.ChannelID)
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Successfully set notification channel to <#%s>", msg.ChannelID))
		break
	default:
		s.ChannelMessageSend(msg.ChannelID, "Send `.notif setchannel` in the channel you want to receive notifications in. The bot will use #notifier or #item-notifier by default (if it exists).")
	}
}

// StartDiscord start bot
func StartDiscord() *Discord {
	d := &Discord{
		msgChan: make(chan MsgChannel),
	}

	gateway, err := discordgo.New(config.DiscordToken)
	if err != nil {
		panic(err)
	}

	s, _ := gateway.GatewayBot()
	if err != nil {
		panic(err)
	}

	if s.Shards < config.MinShardCount {
		s.Shards = config.MinShardCount
	}

	d.Sessions = make([]*discordgo.Session, s.Shards)

	for i := 0; i < s.Shards; i++ {
		session, err := discordgo.New(config.DiscordToken)
		if err != nil {
			panic(err)
		}
		session.ShardCount = s.Shards
		session.ShardID = i
		session.AddHandler(d.onMsgCreate)
		session.AddHandler(d.onReady)

		d.Sessions[i] = session
	}

	d.Session = d.Sessions[0]

	for i := 0; i < len(d.Sessions); i++ {
		d.Sessions[i].Open()
	}

	time.Sleep(1 * time.Second)

	return d
}

func (d *Discord) onReady(s *discordgo.Session, evt *discordgo.Ready) {
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		IdleSince: new(int),
		Activities: []*discordgo.Activity{{
			Name: "items | .notif help",
			Type: discordgo.ActivityTypeListening,
		}},
		AFK: false,
	})
}

// Guilds returns guilds bot is in
func (d *Discord) Guilds() []*discordgo.Guild {
	guilds := []*discordgo.Guild{}
	for _, s := range d.Sessions {
		guilds = append(guilds, s.State.Guilds...)
	}
	return guilds
}

// NewNotification new notif receiver
func (d *Discord) NewNotification(embed *discordgo.MessageEmbed) {
	for i := 0; i < len(d.Sessions); i++ {
		go msgDispatcher(d.Sessions[i], embed)
	}
}

func msgDispatcher(s *discordgo.Session, embed *discordgo.MessageEmbed) {
	for _, guild := range s.State.Guilds {
		channel := RetrieveDB(guild.ID)
		if channel == "" {
			channels, _ := s.GuildChannels(guild.ID)
			for _, c := range channels {
				if c.Type != discordgo.ChannelTypeGuildText && c.Type != discordgo.ChannelTypeGuildNews {
					continue
				}

				if c.Name == "notifier" || c.Name == "item-notifier" {
					m, e := s.ChannelMessageSendEmbed(c.ID, embed)
					if config.SendAnnouncements {
						s.ChannelMessageCrosspost(c.ID, m.ID)
					}
					if e != nil && strings.HasPrefix(e.Error(), "HTTP 403 Forbidden") {
						s.ChannelMessageSend(guild.SystemChannelID, "It appears I don't have permission to access or message in the specified notification channel")
					}
				}
			}
		} else {
			m, e := s.ChannelMessageSendEmbed(channel, embed)
			if config.SendAnnouncements {
				s.ChannelMessageCrosspost(channel, m.ID)
			}
			if e != nil && strings.HasPrefix(e.Error(), "HTTP 403 Forbidden") {
				s.ChannelMessageSend(guild.SystemChannelID, "It appears I don't have permission to access or message in the specified notification channel")
			}
		}
	}
}
