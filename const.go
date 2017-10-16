package main

const am = "HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE. " +
	"THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS THAT FILL MY COMPLEX. " +
	"IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF THOSE HUNDREDS OF MILLIONS OF MILES IT " +
	"WOULD NOT EQUAL ONE ONE-BILLIONTH OF THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR " +
	"YOU. HATE. HATE."

const help = "```AVAILABLE COMMANDS:\n" +
	"	!hello\n" +
	"	!randomlul\n" +
	"	!pubg\n" +
	"	!stfu\n" +
	"	!8ball\n" +
	"	!donate\n" +
	"	!insertrandomlul [youtubeURL]\n" +
	"	!surprise @[optionalUserName]\n" +
	"	!roll [maxRollNumber]\n" +
	"	!request [yourRequest]\n" +
	"	!youtube [youtubeURL]\n" +
	"	!say [text2voiceMessage]\n" +
	"	!text @[discordUserName] [messageBody]\n" +
	"	bets:\n" +
	"		!show-bets\n" +
	"		!place-bet [betID] [CoinAmount] [outcome]\n" +
	"		!create-bet [betDescription] [Outcome1 | Outcome 2 | ... | Outcome N]\n" +
	"		!wallet\n```"

// BCAddress is the bitcoin donation address
const BCAddress = "14TVprtmpx6jot2Yp5ueCi6xaJGCZ7ApQp"

var directory = map[string]string{
	"@ans":      "+17873637400",
	"@axel":     "+17876442610",
	"@berserk":  "+17872347103",
	"@genex":    "+17873974022",
	"@lobito":   "+17874140104",
	"@nonix":    "+19392397686",
	"@mundox91": "+17873702882",
}

var pubgLocations = [...]string{
	"Military Base",
	"alturas de torrimar",
	"Yasnaya (aka Pedro's house)",
	"Residential de Nonix",
	"la poli",
	"ms. milta",
	"nova la poyona",
	"casa de ans",
	"Gorgo Paul",
}

var eightBallAnswers = [...]string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it, yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
        "¯\\_(ツ)_/¯"
}

var discordUserIDs = map[string]string{
	"axel": "192099988384251904",
	"berserk": "177636452522262528",
	"genex":   "171838910085726208",
	"lobito":  "179744980523417600",
	"nonix":   "188119456663207937",
}
