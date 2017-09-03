package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var airHornBuffer = make([][]byte, 0)
var youtubeBuffer = make([][]byte, 0)

func setHandlers(discordSession *discordgo.Session) {
	discordSession.AddHandler(messageListenAndRespond)
}

func messageListenAndRespond(s *discordgo.Session, m *discordgo.MessageCreate) {

	messageWithoutUserMentions := strings.ToLower(m.ContentWithMentionsReplaced())

	fmt.Printf("Message is '%s' from '%s'\n", messageWithoutUserMentions, m.Author.Username)

	// ignore message posted by AM
	if m.Author.ID == s.State.User.ID {
		return
	}

	// ignore message if is for everyone
	if m.MentionEveryone && !strings.Contains(messageWithoutUserMentions, "!text") {
		fmt.Println("Ignored: Message is for everyone")
		return
	}

	splitMessageWithUserMentions := strings.Fields(m.Content)
	firstUserMentioned := strings.SplitAfterN(messageWithoutUserMentions, " ", 2)[0]

	// check if you are the user being mentioned at the BEGINNING of the message
	if !strings.Contains(firstUserMentioned, "@ambot") {
		fmt.Println("Ignored: @ambot is not mentioned at the beginnning of the message")
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

	// get channel and guild
	msgChannel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	guildID := msgChannel.GuildID

	msgGuild, err := s.State.Guild(guildID)
	if err != nil {
		return
	}

	fmt.Printf("Message received: %s: %s \n", m.Author.ID, messageWithoutUserMentions)

	// check message and reply/react to message
	/***** HELLO *****/
	if strings.Contains(messageCommand, "!hello") {
		s.ChannelMessageSend(m.ChannelID, am)

		/***** HELP MENU *****/
	} else if strings.Contains(messageCommand, "!help") {
		s.ChannelMessageSend(m.ChannelID, help)

		/***** PUBG LOCATIONS *****/
	} else if strings.Contains(messageCommand, "!pubg") {
		rand.Seed(time.Now().Unix())
		s.ChannelMessageSend(m.ChannelID, pubgLocations[rand.Intn(len(pubgLocations))])

		/***** LOG REQUESTS *****/
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

	} else if strings.Contains(messageCommand, "!randomlul") {
		rand.Seed(time.Now().Unix())
		s.ChannelMessageSend(m.ChannelID, lulPlaylist[rand.Intn(len(lulPlaylist))])

		/***** YOUTUBE STREAMING *****/
	} else if strings.Contains(messageCommand, "!youtube") {

		vc, err := joinUserVoiceChannel(s, m.Author.ID)
		if err != nil {
			fmt.Printf("ERR is %s", err)
			return
		}

		vid := splitMessageWithUserMentions[2]

		// download youtube vid
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
							//s.VoiceReady = false
							vc.Disconnect()
							break
						}
						fmt.Println("Audio error: ", err)
					}
					vc.OpusSend <- opus
				}
			}
		}

		/***** TEXT MESSAGES *****/
	} else if strings.Contains(messageCommand, "!text") {
		if len(splitMessage) < 4 {
			s.ChannelMessageSend(m.ChannelID, "You forgot your message. Moron.")
			return
		}
		toUser := splitMessage[2]
		smsMsg := strings.SplitAfterN(messageWithoutUserMentions, " ", 4)

		textSuccess := false

		if toUser == "@everyone" {
			for user := range directory {
				textSuccess = sendSMS(user, m.Author.Username, smsMsg[len(smsMsg)-1])
			}
		} else {
			textSuccess = sendSMS(toUser, m.Author.Username, smsMsg[len(smsMsg)-1])
		}

		if !textSuccess {
			s.ChannelMessageSend(m.ChannelID, "User not found in directory. Fool.")
			return
		} else {
			s.ChannelMessageSend(m.ChannelID, "Message sent.")
		}

		/*** AIRHORN USER ***/
	} else if strings.Contains(messageCommand, "!surprise") {
		fmt.Println("AIRHORN!!")
		airhornUser := m.Author.ID

		if len(splitMessage) == 2 {
			airhornUser = m.Author.ID
		} else if len(splitMessage) == 3 {
			airhornUser = splitMessageWithUserMentions[2]
			for _, user := range m.Mentions {
				if user.Username != "AMBot" {
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

		/***** CREATE BET EVENTS *****/
	} else if strings.Contains(messageCommand, "!create-bet") {
		eventNameStart := false
		eventName := ""

		// parse event name
		for _, word := range splitMessage {
			if strings.Contains(word, messageCommand) {
				eventNameStart = true
			} else if word == "|" {
				eventNameStart = false
			} else if eventNameStart {
				eventName = eventName + word + " "
			}
		}

		// parse event name, outcomes, creator and store in DB
		outcomes := strings.Split(messageWithoutUserMentions, "|")
		myBetEvent := new(betEvent)
		myBetEvent.outcomes = make([]string, len(outcomes)-1)
		for i, outcome := range outcomes {
			fmt.Printf("int is: %d\n", i)
			if i == 0 {

			} else {
				myBetEvent.outcomes[i-1] = strings.TrimSpace(outcome)
			}
		}
		myBetEvent.name = strings.TrimSpace(eventName)
		myBetEvent.creator = m.Author.Username
		myBetEvent.active = true

		fmt.Printf("event is: %+v\n", myBetEvent)
		insertBetEventQuery(myBetEvent)
		getBetEventQuery(1)

		/****** PLACE BETS *****/
	} else if strings.Contains(messageCommand, "!place-bet") {
		if len(splitMessage) != 5 {
			s.ChannelMessageSend(m.ChannelID, "Your bet is not properly formatted. Imbecile")
			return
		}

		// parse betEventID and betAmount
		betEventID, err := strconv.Atoi(splitMessage[2])
		if err != nil {
			fmt.Printf("Error parsing betting ID. Make sure it uses digits only: %s\n", err)
			return
		}
		betAmount, err := strconv.Atoi(splitMessage[3])
		if err != nil {
			fmt.Printf("Error parsing betting amount. Make sure it uses digits only: %s\n", err)
			return
		}
		betOutcome := ""
		for i := 4; i < len(splitMessage); i++ {
			betOutcome = betOutcome + splitMessage[i]
		}
		if err != nil {
			fmt.Printf("Error parsing betting outcome. Make sure it is one of accepted outcomes: %s\n", err)
			return
		}

		gamblerBet := map[string]bet{strings.ToLower(m.Author.Username): {Money: betAmount, Outcome: betOutcome}}

		updatedGamble, err := updateBetEventGamblerQuery(betEventID, gamblerBet)

		if updatedGamble == 1 {
			s.ChannelMessageSend(m.ChannelID, "Your bet has been placed. Good luck!")
		} else if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			fmt.Printf("error is: %s\n", err)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Your transaction could not be completed. Please try again later")
		}

		/***** SHOW ALL BET EVENTS */
	} else if strings.Contains(messageCommand, "!show-bets") {
		var allBetEvents []betEvent
		allBetEventsStr := ""
		allBetEvents = getAllBetEventsQuery()

		allBetEventsStr = allBetEventsStr + "```\n"
		allBetEventsStr = allBetEventsStr + "Event Name | Outcomes | Bet Creator | Gamblers | Active | Bet Event ID\n"
		for _, myBetEvent := range allBetEvents {
			allBetEventsStr = allBetEventsStr + fmt.Sprintf("%s\n", myBetEvent.discordPrettyPrint())
		}
		allBetEventsStr = allBetEventsStr + "```"

		s.ChannelMessageSend(m.ChannelID, allBetEventsStr)

		/***** GET YOUR WALLET *****/
	} else if strings.Contains(messageCommand, "!wallet") {
		userAmCoins, err := getAmCoins(strings.ToLower(m.Author.Username))
		if err != nil {
			fmt.Printf("error getting coins: %s\n", err)
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s: %d`", m.Author.Username, userAmCoins))
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// join the voice channel
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range airHornBuffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	vc.Disconnect()

	return nil
}
