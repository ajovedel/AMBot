package main

const am = "HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE. " +
	"THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS THAT FILL MY COMPLEX. " +
	"IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF THOSE HUNDREDS OF MILLIONS OF MILES IT " +
	"WOULD NOT EQUAL ONE ONE-BILLIONTH OF THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR " +
	"YOU. HATE. HATE."

const help = "```AVAILABLE COMMANDS:\n" +
	"	!hello\n" +
	"	!randomlul\n" +
	"	!surprise @[optionalUserName]\n" +
	"	!pubg\n" +
	"	!roll [maxRollNumber]\n" +
	"	!request [yourRequest]\n" +
	"	!youtube [youtubeURL]\n" +
	"	!stfu\n" +
	"	!say [text2voiceMessage]\n" +
	"	!text @[discordUserName] [messageBody]\n" +
	"	bets:\n" +
	"		!show-bets\n" +
	"		!place-bet [betID] [CoinAmount] [outcome]\n" +
	"		!create-bet [betDescription] [Outcome1 | Outcome 2 | ... | Outcome N]\n" +
	"		!wallet\n```"

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

var lulPlaylist = [...]string{
	"https://www.youtube.com/watch?v=FVHWxESdN4o",
	"https://www.youtube.com/watch?v=iq_d8VSM0nw",
	"https://www.youtube.com/watch?v=ygI-2F8ApUM",
	"https://www.youtube.com/watch?v=5-2nByd2cr4",
	"https://www.youtube.com/watch?v=PeihcfYft9w",
	"https://www.youtube.com/watch?v=AP7utU8Efow",
	"https://www.youtube.com/watch?v=tVj0ZTS4WF4",
	"https://www.youtube.com/watch?v=1EKTw50Uf8M",
	"https://www.youtube.com/watch?v=oAQuUuxnsUg",
	"https://www.youtube.com/watch?v=CWzUK4Qjsws",
	"https://www.youtube.com/watch?v=Eo_gJwXxshQ",
	"https://www.youtube.com/watch?v=5eveNk3o1ME",
	"https://www.youtube.com/watch?v=_pqGNSXC9to",
	"https://www.youtube.com/watch?v=ZZ5LpwO-An4",
	"https://www.youtube.com/watch?v=_5fQZhv0poU",
	"https://www.youtube.com/watch?v=0tdyU_gW6WE",
	"https://www.youtube.com/watch?v=jjdl2Yp6rxk",
	"https://www.youtube.com/watch?v=7eKv4BEujFU",
	"https://www.youtube.com/watch?v=Gm7lcZiLOus",
	"https://www.youtube.com/watch?v=Ns7Z8ag4oSY",
	"https://www.youtube.com/watch?v=J---aiyznGQ",
	"https://www.youtube.com/watch?v=Q16KpquGsIc",
	"https://www.youtube.com/watch?v=fI4ZhW1anKY",
	"https://www.youtube.com/watch?v=9C_HReR_McQ",
	"https://www.youtube.com/watch?v=GKFKl22j1RY",
	"https://www.youtube.com/watch?v=dP9Wp6QVbsk",
	"https://www.youtube.com/watch?v=vTIIMJ9tUc8",
	"https://www.youtube.com/watch?v=PeS_lZySMf8",
	"https://www.youtube.com/watch?v=EDjBzeTZxG4",
	"https://www.youtube.com/watch?v=XcicOBS9mBU",
	"https://www.youtube.com/watch?v=I53HDr0-Qew",
	"https://www.youtube.com/watch?v=pDVORKo8rYs",
	"https://www.youtube.com/watch?v=G7RgN9ijwE4",
	"https://www.youtube.com/watch?v=4d_FvgQ1csE",
	"https://www.youtube.com/watch?v=Kppx4bzfAaE",
	"https://www.youtube.com/watch?v=fUY9FmBvqms",
	"https://www.youtube.com/watch?v=1Bix44C1EzY",
	"https://www.youtube.com/watch?v=1tF2dF67Q2c",
	"https://www.youtube.com/watch?v=zBWPL7oD9X8",
	"https://www.youtube.com/watch?v=pD_imYhNoQ4",
	"https://www.youtube.com/watch?v=kJa2kwoZ2a4",
	"https://www.youtube.com/watch?v=69iSXks1bes",
	"https://www.youtube.com/watch?v=wsO-Td0hqXo",
}
