package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var buffer = make([][]byte, 0)

func setHandlers(discordSession *discordgo.Session) {
	discordSession.AddHandler(messageListenAndRespond)
}

func messageListenAndRespond(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore message posted by AM
	if m.Author.ID == s.State.User.ID {
		return
	}

	// ignore message if is for everyone
	if m.MentionEveryone {
		fmt.Printf("IGNORED11111")
		return
	}

	// check if you are the users
	/*
		if len(m.Mentions) != 1 {
			fmt.Printf("IGNORED22222")
			return
		}*/

	// check AM is being mentioned
	/*
		if m.Mentions[0].ID != s.State.User.ID {
			fmt.Printf("IGNORED333333")
			return
		}*/

	// get guild ID
	// get channel that the message came from.
	msgChannel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	guildID := msgChannel.GuildID

	msgGuild, err := s.State.Guild(guildID)
	if err != nil {
		return
	}

	fmt.Printf("Message received: %s\n", m.ContentWithMentionsReplaced())

	// check message and reply/react to message
	if strings.Contains(m.Content, "hello") {
		s.ChannelMessageSend(m.ChannelID, am)

	} else if strings.Contains(m.Content, "!help") {
		s.ChannelMessageSend(m.ChannelID, help)

	} else if strings.Contains(m.Content, "!pubg") {
		rand.Seed(time.Now().Unix())
		s.ChannelMessageSend(m.ChannelID, pubgLocations[rand.Intn(len(pubgLocations))])

	} else if strings.Contains(m.Content, "!request") {
		f, err := os.OpenFile("requests.log", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("ERR is %s", err)
		}
		defer f.Close()

		if _, err = f.WriteString(m.ContentWithMentionsReplaced()); err != nil {
			fmt.Printf("ERR is %s", err)
		}
		s.ChannelMessageSend(m.ChannelID, "Request has been logged and will be reviewed.")

	} else if strings.Contains(m.Content, "axel") && strings.Contains(m.Content, "awesome") {
		err = s.MessageReactionAdd(msgChannel.ID, m.ID, "👍")
		if err != nil {
			fmt.Printf("ERR Is: %s", err)
		}

	} else if strings.Contains(m.Content, "luis") && strings.Contains(m.Content, "awesome") {
		err = s.MessageReactionAdd(msgChannel.ID, m.ID, "💩")
		if err != nil {
			fmt.Printf("ERR Is: %s", err)
		}

	} else if strings.Contains(m.Content, "pedro") && strings.Contains(m.Content, "awesome") {
		err = s.MessageReactionAdd(msgChannel.ID, m.ID, "😂")
		if err != nil {
			fmt.Printf("ERR Is: %s", err)
		}

	} else if strings.Contains(m.Content, "!randomlul") {
		rand.Seed(time.Now().Unix())
		s.ChannelMessageSend(m.ChannelID, lulPlaylist[rand.Intn(len(lulPlaylist))])

	} else if strings.Contains(m.Content, "!text") {
		splitContent := strings.Fields(m.Message.ContentWithMentionsReplaced())
		user := splitContent[2]
		smsMsg := strings.SplitAfterN(m.Message.ContentWithMentionsReplaced(), " ", 4)
		fmt.Printf("msg is: %s\n", m.Message.ContentWithMentionsReplaced())
		fmt.Printf("User is: %s\n", user)
		fmt.Printf("msg is: %s\n", smsMsg[len(smsMsg)-1])

		sendSMS(user, smsMsg[len(smsMsg)-1])

		s.ChannelMessageSend(m.ChannelID, "Message sent.")

	} else if strings.Contains(m.Content, "!surprise") {
		fmt.Println("AIRHORN!!")
		for _, vs := range msgGuild.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, guildID, vs.ChannelID)
				if err != nil {
					fmt.Printf("ERR Is: %s", err)
				}
			}
		}
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
