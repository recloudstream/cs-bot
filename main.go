package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"

	_ "embed"

	"github.com/bwmarrin/discordgo"
)

//go:embed .token
var TOKEN string

func main() {
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		panic(err)
	}
	dg.Identify.Intents |= discordgo.IntentsGuildMessages
	dg.Identify.Intents |= discordgo.IntentsDirectMessages
	dg.Identify.Intents |= discordgo.IntentsGuildMessageReactions
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		panic(err)
	}
	defer dg.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

var helpRegex = regexp.MustCompile("(?imsU)where|how|hwo|i need|i can([' ]|no)?t|help|list|please|plz|can (someone|anyone)|((is|does)(n[ ']?t| not) work)")
var subjectRegex = regexp.MustCompile("(?imsU)(source|repo|(cloud ?stream)|ext[a-z]{3,10}ons?|short ?code|re[a-z]{1,20}tor(y|ies)|provider)")
var cmdRegex = regexp.MustCompile(`^[!\$\/\.](repos?|ext[a-z]{0,15}|list|re[a-z]{0,5}tor(y|ies)|links?|providers?)$`)
var notWorkeyRegex = regexp.MustCompile("(?imsU)((is|does)(n[ ']?t| not) work)|broke")

func askedForHelp(msg string) bool {
	str1 := helpRegex.FindString(msg)
	str2 := subjectRegex.FindString(msg)
	if str1 != "" && str2 != "" {
		fmt.Println(str1, str2)
	}
	return str1 != "" && str2 != "" && len(msg) < 150
}

//go:embed body1.md
var embedBody1 string

//go:embed body2.md
var embedBody2 string

var msg1 = &discordgo.MessageSend{
	Content: "Hello! My algorithm detected that you might need help with Cloudstream repositories. If you didn't, please ignore this message.",
	Embeds: []*discordgo.MessageEmbed{
		{
			Title:       "Cloudstream Repository",
			Description: embedBody1,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Cloudstream",
				IconURL: "https://github.com/recloudstream.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "This is an unofficial bot. Please forward any DMCA complaints to GitHub.",
			},
		},
	},
}

var msg2 = &discordgo.MessageSend{
	Embeds: []*discordgo.MessageEmbed{
		{
			Title:       "Cloudstream Repository",
			Description: embedBody2,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Cloudstream",
				IconURL: "https://github.com/recloudstream.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "This is an unofficial bot. Please forward any DMCA complaints to GitHub.",
			},
		},
	},
}

func onNotWorkey(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if notWorkeyRegex.FindString(m.Content) != "" {
		_, err := s.ChannelMessageSendComplex(m.ChannelID, msg2)
		if err != nil {
			fmt.Println(err)
		}
		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
			fmt.Println(err)
		}
		return true
	}
	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	for _, mention := range m.Mentions {
		if mention.ID == s.State.User.ID {
			if onNotWorkey(s, m) {
				return
			} else {
				break
			}
		}
	}

	if m.Content == ".nofap" {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.discordapp.com/attachments/737725084881387652/1087722994010431599/XgmvI0N.mp4")
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err == nil && channel != nil && channel.Type == discordgo.ChannelTypeDM {
		onNotWorkey(s, m)
		return
	}

	if askedForHelp(m.Content) || cmdRegex.MatchString(m.Content) {
		dm, err := s.UserChannelCreate(m.Author.ID)
		if dm != nil && err == nil {
			_, err = s.ChannelMessageSendComplex(dm.ID, msg1)
			if err == nil {
				if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
					fmt.Println(err)
				}
				return
			}
		}
		msg1.Reference = &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		}
		_, err = s.ChannelMessageSendComplex(m.ChannelID, msg1)
		msg1.Reference = nil
		if err != nil {
			fmt.Println(err)
		}

		if err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅"); err != nil {
			fmt.Println(err)
		}
	}
}