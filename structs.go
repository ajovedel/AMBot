package main

import (
	"fmt"
)

type betEvent struct {
	name       string
	outcomes   []string
	creator    string
	gamblers   map[string]bet
	active     bool
	betEventID int
}

type bet struct {
	Money   int
	Outcome string
}
type amCoinBank struct {
	name  string
	money int
}

func (myBetEvent *betEvent) discordPrettyPrint() string {

	finalString := ""
	finalString = finalString + myBetEvent.name + " | "
	for _, outcome := range myBetEvent.outcomes {
		if outcome == myBetEvent.outcomes[len(myBetEvent.outcomes)-1] {
			finalString = finalString + outcome + " | "
		} else {
			finalString = finalString + outcome + " - "
		}
	}
	finalString = finalString + myBetEvent.creator + " | "
	finalString = finalString + fmt.Sprintf("%v", myBetEvent.gamblers) + " | "
	finalString = finalString + fmt.Sprintf("%t", myBetEvent.active) + " | "
	finalString = finalString + fmt.Sprintf("%d", myBetEvent.betEventID)

	return finalString
}
