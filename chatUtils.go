package main

import (
	"fmt"
	"io"
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

	messageWithoutUserMentions := strings.ToLower(m.ContentWithMentionsReplaced())
	splitMessageWithUserMentions := strings.Fields(m.Content)
	firstUserMentioned := strings.SplitAfterN(messageWithoutUserMentions, " ", 2)[0]

	// check if you are the user being mentioned at the BEGINNING of the message
	if !strings.Contains(firstUserMentioned, "@ambot") {
		fmt.Printf("IGNORED333333")
		return
	}

	// parse command
	splitMessage := strings.Fields(messageWithoutUserMentions)
	if len(splitMessage) == 1 {
		s.ChannelMessageSend(m.ChannelID, "No command given. Idiot.")
		return
	} else if !strings.Contains(splitMessage[1], "!") {
		s.ChannelMessageSend(m.ChannelID, "No command given. Idiot.")
		return
	}
	messageCommand := splitMessage[1]
	fmt.Printf("command is %s\n", messageCommand)

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

	fmt.Printf("Message received: %s\n", messageWithoutUserMentions)

	// check message and reply/react to message
	if strings.Contains(messageCommand, "!hello") {
		s.ChannelMessageSend(m.ChannelID, am)

	} else if strings.Contains(messageCommand, "!help") {
		s.ChannelMessageSend(m.ChannelID, help)

	} else if strings.Contains(messageCommand, "!pubg") {
		rand.Seed(time.Now().Unix())
		s.ChannelMessageSend(m.ChannelID, pubgLocations[rand.Intn(len(pubgLocations))])

	} else if strings.Contains(messageCommand, "!request") {
		f, err := os.OpenFile("requests.log", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("ERR is %s", err)
		}
		defer f.Close()

		if _, err = f.WriteString(messageWithoutUserMentions + "\n"); err != nil {
			fmt.Printf("ERR is %s", err)
		}
		s.ChannelMessageSend(m.ChannelID, "Request has been logged and will be reviewed.")

		/*

			} else if strings.Contains(messageWithoutUserMentions, "axel") && strings.Contains(m.Content, "awesome") {
				err = s.MessageReactionAdd(msgChannel.ID, m.ID, "💩")
				if err != nil {
					fmt.Printf("ERR Is: %s", err)
				}

			} else if strings.Contains(messageWithoutUserMentions, "luis") && strings.Contains(m.Content, "awesome") {
				err = s.MessageReactionAdd(msgChannel.ID, m.ID, "👍")
				if err != nil {
					fmt.Printf("ERR is: %s", err)
				}

			} else if strings.Contains(messageWithoutUserMentions, "pedro") && strings.Contains(m.Content, "awesome") {
				err = s.MessageReactionAdd(msgChannel.ID, m.ID, "😂")
				if err != nil {
					fmt.Printf("ERR is: %s", err)
				}
		*/

	} else if strings.Contains(messageCommand, "!randomlul") {
		rand.Seed(time.Now().Unix())
		s.ChannelMessageSend(m.ChannelID, lulPlaylist[rand.Intn(len(lulPlaylist))])

	} else if strings.Contains(messageCommand, "!youtube") {
		// Open a youtube stream

		fmt.Println("LEL")

		vc, err := joinUserVoiceChannel(s, m.Author.ID)
		if err != nil {
			fmt.Printf("ERR is %s", err)

			return
		}

		fmt.Printf("Vid is '%x'\n", splitMessage[2])
		//fmt.Printf("Vid is '%s'\n", reflect.TypeOf(splitMessage[2]))

		fmt.Printf("vid is '%x'\n", "https://www.youtube.com/watch?v=ZqNpXJwgO8o")

		//vid := make(string, len(splitMessage[2]))
		vid := splitMessageWithUserMentions[2]
		//copy(vid, splitMessage[2])
		fmt.Printf("Vid is '%x'\n", vid)

		//yt, err := youtubePy(vid.)
		yt, err := youtubePy(vid)

		if err != nil {
			fmt.Printf("ERR is: %s", err)
			return
		}

		// Create opus stream
		stream, err := convertToOpus(yt)
		if err != nil {
			fmt.Printf("ERR is %s", err)
			return
		}

		for _, vs := range msgGuild.VoiceStates {
			if vs.UserID == m.Author.ID {
				for {
					opus, err := readOpus(stream)
					if err != nil {
						if err == io.ErrUnexpectedEOF || err == io.EOF {
							fmt.Printf("ERR is: %s", err)
							break
						}
						fmt.Println("Audio error: ", err)
					}
					vc.OpusSend <- opus
				}
			}
		}

	} else if strings.Contains(messageCommand, "!text") {
		if len(splitMessage) < 4 {
			s.ChannelMessageSend(m.ChannelID, "You forgot your message. Moron.")
			return
		}
		user := splitMessage[2]
		smsMsg := strings.SplitAfterN(messageWithoutUserMentions, " ", 4)
		fmt.Printf("msg is: %s\n", messageWithoutUserMentions)
		fmt.Printf("Uuser is: %s\n", user)
		fmt.Printf("sms message is: %s\n", smsMsg[len(smsMsg)-1])

		textSuccess := sendSMS(user, smsMsg[len(smsMsg)-1])

		if !textSuccess {
			s.ChannelMessageSend(m.ChannelID, "User not found in directory. Fool.")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Message sent.")

	} else if strings.Contains(messageCommand, "!surprise") {
		fmt.Println("AIRHORN!!")
		airhornUser := m.Author.ID

		if len(splitMessage) == 2 {
			airhornUser = m.Author.ID
		} else if len(splitMessage) == 3 {
			airhornUser = splitMessageWithUserMentions[2]
			for _, user := range m.Mentions {
				if user.Username != "AMBot" {
					fmt.Printf("airhorn userrrr: %s\n", user.Username)
					airhornUser = user.ID
				}
			}
		} else {
			return
		}
		fmt.Printf("airhorn user is: %s\n", airhornUser)
		for _, vs := range msgGuild.VoiceStates {
			fmt.Printf("vs.UserID %s\n", vs.UserID)
			if vs.UserID == airhornUser {
				err = playSound(s, guildID, vs.ChannelID)
				if err != nil {
					fmt.Printf("ERR is: %s", err)
				}
			}
		}

	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	fmt.Printf("Buffer is: %v", buffer)

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
