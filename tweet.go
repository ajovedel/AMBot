package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func setTwitterHandlers(discordSession *discordgo.Session) {

	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	//trumpTwitterID := int64(4916819346)
	trumpTwitterID := int64(25073877)

	// loop every 20 mins
	for {

		tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			UserID: trumpTwitterID,
			Count:  10,
		})

		if err != nil {
			println("tweet error")
			return
		}

		// loop over tweets in reverse, to post older tweets first and newer tweets last
		for i := len(tweets) - 1; i >= 0; i-- {
			fmt.Println(tweets[i])
			if checkIfTweetHasBeenPosted(tweets[i].ID) {
				continue
			} else {
				postTweet(discordSession, tweets[i])
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(20 * time.Minute)
	}

}

func checkIfTweetHasBeenPosted(tweetID int64) bool {

	posted := false

	query := fmt.Sprintf("SELECT posted FROM tweetHistory WHERE tweetID = %d", tweetID)

	row, err := ambotDB.Query(query)

	defer row.Close()
	if err != nil {
		return false
	}

	for row.Next() {
		err := row.Scan(&posted)
		if err != nil {
			fmt.Printf("Error scanning vars: %s\n", err)
			return false
		}
	}

	if posted {
		fmt.Println("POSTED IS TRUE")
	} else {
		fmt.Println("POSTED IS FALSE")
	}

	return posted
}

func postTweet(discordSession *discordgo.Session, tweet twitter.Tweet) {
	fmt.Println("Posting...")
	query := "INSERT INTO tweetHistory (tweetID, posted) VALUES (?, 1)"

	tx, err := ambotDB.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(tweet.ID)
	if err != nil {
		fmt.Println(err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	if affectedRows != 1 {
		return
	}

	_, err = discordSession.ChannelMessageSend("348922901531066371", tweet.Text)
	if err != nil {
		log.Printf("Derp, %s", err)
	}

}
