package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var ambotDB *sql.DB

func init() {
	fmt.Println("Initializing database driver")

	var err error

	if _, err := os.Stat("./db/ambot.db"); os.IsNotExist(err) {
		fmt.Println("Database db/ambot.db does not exist. Please create or copy the database")
		return
	}
	ambotDB, err = sql.Open("sqlite3", "./db/ambot.db")
	if err != nil {
		fmt.Println(err)
	}
}

func getAllBetEventsQuery() []betEvent {
	query := "SELECT COUNT(*) FROM betEvents"

	rows, err := ambotDB.Query(query)
	if err != nil {
		fmt.Printf("Error in query: %s\n", err)
	}
	var numRows int

	for rows.Next() {
		rows.Scan(&numRows)
	}

	fmt.Printf("Rows is %d\n", numRows)
	var allBetEvents []betEvent
	for i := 0; i < numRows; i++ {
		betEventTemp, err := getBetEventQuery(i + 1)
		allBetEvents = append(allBetEvents, betEventTemp)
		if err != nil {
			fmt.Printf("Err querying row %d, err is: %s", i+1, err)
		}
	}

	return allBetEvents

}

func getBetEventQuery(id int) (betEvent, error) {
	var myBetEvent betEvent

	var outcomesBytes []byte
	var gamblersBytes []byte
	var activeBytes []byte

	query := fmt.Sprintf("SELECT rowid, name, creator, outcomes, gamblers, active from betEvents where rowid = %d", id)
	rows, err := ambotDB.Query(query)

	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return myBetEvent, err
	}
	for rows.Next() {
		err := rows.Scan(&myBetEvent.betEventID, &myBetEvent.name, &myBetEvent.creator, &outcomesBytes, &gamblersBytes, &activeBytes)
		if err != nil {
			fmt.Printf("Error scanning vars: %s\n", err)
			return myBetEvent, nil
		}
	}

	// unmarshall the data blobs
	err = json.Unmarshal(outcomesBytes, &myBetEvent.outcomes)
	if err != nil {
		fmt.Printf("Err is: %s\n", err)
		return myBetEvent, err
	}
	err = json.Unmarshal(gamblersBytes, &myBetEvent.gamblers)
	if err != nil {
		fmt.Printf("Err is: %s\n", err)
		return myBetEvent, err
	}
	err = json.Unmarshal(activeBytes, &myBetEvent.active)
	if err != nil {
		fmt.Printf("Err is: %s\n", err)
		return myBetEvent, err
	}

	fmt.Printf("unmarshalled betEvent is: %v\n", myBetEvent)

	return myBetEvent, nil

}

func insertBetEventQuery(myBetEvent *betEvent) int64 {

	query := fmt.Sprint("INSERT INTO betEvents (name, creator, outcomes, gamblers, active) " +
		"VALUES (?,?,?,?,?)")

	// marshall slices and structs
	outcomesByte, err := json.Marshal(myBetEvent.outcomes)
	if err != nil {
		fmt.Printf("Error marshalling outcomes: %s\n", err)
		return 0
	}
	gamblerBytes, err := json.Marshal(myBetEvent.gamblers)
	if err != nil {
		fmt.Printf("Error marshalling gamblers: %s\n", err)
		return 0
	}
	activeBytes, err := json.Marshal(myBetEvent.active)
	if err != nil {
		fmt.Printf("Error marshalling active: %s\n", err)
		return 0
	}

	tx, err := ambotDB.Begin()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer stmt.Close()

	res, err := stmt.Exec(myBetEvent.name, myBetEvent.creator, outcomesByte, gamblerBytes, activeBytes)
	if err != nil {
		fmt.Println(err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	tx.Commit()

	return affectedRows
}

func updateBetEventGamblerQuery(betEventID int, gamblerBet map[string]bet) (int64, error) {

	myBetEvent, err := getBetEventQuery(betEventID)
	if err != nil {
		fmt.Printf("Error getting bet: %s\n", err)
		return 0, err
	}

	if myBetEvent.gamblers == nil {
		myBetEvent.gamblers = make(map[string]bet)
	}

	// check if gambler has placed bet before in event
	for k := range gamblerBet {
		if _, ok := myBetEvent.gamblers[k]; ok {
			fmt.Printf("Gambler has already placed a bet. Cannot accept this new bet\n")
			return 0, errors.New("You already placed a bet. Cannot accept this new bet")
		}
	}

	// check if gambler has enough amCoins to place the bet, if so do it.
	for k := range gamblerBet {
		amCoins, err := getAmCoins(k)
		fmt.Printf("My bet is: %v\n", gamblerBet[k])
		fmt.Printf("My wallet is: %d\n", amCoins)

		if err != nil {
			fmt.Printf("Error getting wallet: %s\n", err)
			return 0, err
		} else if gamblerBet[k].Money > amCoins {
			fmt.Printf("Not enough coins for the bet. Looser\n")
			return 0, nil
		} else {
			updateAmCoins(k, amCoins-gamblerBet[k].Money)
		}
	}

	//merge both gamblerBets into one map
	for k, v := range gamblerBet {
		myBetEvent.gamblers[k] = v
	}

	fmt.Printf("Gambler bet after merge is: %v\n", myBetEvent.gamblers)

	gamblerBetBytes, err := json.Marshal(myBetEvent.gamblers)
	if err != nil {
		fmt.Printf("Error marshalling gamblers: %s\n", err)
		return 0, err
	}

	fmt.Printf("gamblerBetBytes is: %v\n", gamblerBetBytes)

	query := fmt.Sprint("UPDATE betEvents SET gamblers = ? WHERE rowid=?")

	tx, err := ambotDB.Begin()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(gamblerBetBytes, betEventID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	tx.Commit()

	return affectedRows, nil

}

func getAmCoins(name string) (int, error) {

	query := fmt.Sprintf("SELECT amCoins FROM coinBank where name = '%s'", name)
	rows, err := ambotDB.Query(query)

	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	amCoins := 0
	for rows.Next() {
		err := rows.Scan(&amCoins)
		if err != nil {
			fmt.Printf("Error scanning vars: %s\n", err)
			return 0, err
		}
	}
	return amCoins, nil
}

func updateAmCoins(name string, newCoinAmount int) (int64, error) {
	query := fmt.Sprint("UPDATE coinBank SET amCoins = ? WHERE name=?")

	tx, err := ambotDB.Begin()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, newCoinAmount)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	tx.Commit()

	return affectedRows, nil
}

func randomYoutubeVid() (string, error) {
	query := "SELECT rowid FROM youtubeVids"

	rows, err := ambotDB.Query(query)
	if err != nil {
		fmt.Printf("Error in query: %s\n", err)
		return "", err
	}

	//var rowIDs []int
	rowIDs := make([]int, 0)
	rowTemp := 0
	rowNum := 0

	for rows.Next() {
		err = rows.Scan(&rowTemp)
		if err != nil {
			fmt.Printf("error in query: %s\n", err)
			return "", err
		}
		rowIDs = append(rowIDs, rowTemp)
		rowNum = rowNum + 1
	}

	fmt.Printf("rowIDs is: %v", rowIDs)

	rows.Close()

	rand.Seed(time.Now().Unix())
	randRowID := rand.Intn(rowNum)

	querySingleVid := fmt.Sprintf("SELECT name FROM youtubeVids where rowid=%d", rowIDs[randRowID])

	rows, err = ambotDB.Query(querySingleVid)
	if err != nil {
		fmt.Printf("Error in query: %s\n", err)
		return "", err
	}

	youtubeVid := ""

	for rows.Next() {
		err = rows.Scan(&youtubeVid)
		if err != nil {
			fmt.Printf("Error in query: %s\n", err)
			return "", err
		}
	}
	rows.Close()

	return youtubeVid, nil
}

func insertYoutubeVid(youtubeURL string) (int64, error) {
	query := fmt.Sprint("INSERT INTO youtubeVids (name) VALUES (?)")

	tx, err := ambotDB.Begin()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(youtubeURL)
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	tx.Commit()

	return affectedRows, nil
}
