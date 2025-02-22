package main

import (
	//"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
    "io/ioutil"
    "path/filepath"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	common "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

//TODO
//"Catch up on the score" - dont remember what this is lol

//FUNCTIONAL CHANGES
//add verification for if a round event has triggered so far in the round (avoid double roundEnds)
//check for game start without pistol (if we have bad demo)
//Add backend support
//Add lurker/anchor stuff
//Add team economy round stats (ecos, forces, etc)

//CLEAN CODE
//create a outputPlayer function to clean up output.go
//convert rating calculations to a function
//actually use killValues lmao

const DEBUG = false

//const suppressNormalOutput = false

//globals
const printChatLog = true
const printDebugLog = true

const tradeCutoff = 4 // in seconds
var multikillBonus = [...]float64{0, 0, 0.3, 0.7, 1.2, 2}
var clutchBonus = [...]float64{0, 0.2, 0.6, 1.2, 2, 3}
var killValues = map[string]float64{
	"attacking":     1.2, //base values
	"defending":     1.0,
	"bombDefense":   1.0,
	"retake":        1.2,
	"chase":         0.8,
	"exit":          0.6,
	"t_consolation": 0.5,
	"gravy":         0.6,
	"punish":        0.8,
	"entry":         0.8, //multipliers
	"t_opener":      0.3,
	"ct_opener":     0.5,
	"trade":         0.3,
	"flashAssist":   0.2,
	"assist":        0.15,
}

type game struct {
	//winnerID         int
	winnerClanName   string
	rounds           []*round
	potentialRound   *round
	teams            map[string]*team
	flags            flag
	mapName          string
	tickRate         int
	tickLength       int
	roundsToWin      int //30 or 16
	totalPlayerStats map[uint64]*playerStats
	totalTeamStats   map[string]*teamStats
	playerOrder      []uint64
	teamOrder        []string
	totalRounds      int
}

type flag struct {
	//all our sentinals and shit
	hasGameStarted            bool
	isGameLive                bool
	isGameOver                bool
	prePlant                  bool
	postPlant                 bool
	postWinCon                bool
	roundIntegrityStart       int
	roundIntegrityEnd         int
	roundIntegrityEndOfficial int

	//for the round (gets reset on a new round) maybe should be in a new struct
	tAlive        int
	ctAlive       int
	tMoney        bool
	tClutchVal    int
	ctClutchVal   int
	tClutchSteam  uint64
	ctClutchSteam uint64
	openingKill   bool
}

type team struct {
	//id    int //meaningless?
	name  string
	score int
}

type teamStats struct {
	winPoints      float64
	impactPoints   float64
	tWinPoints     float64
	ctWinPoints    float64
	tImpactPoints  float64
	ctImpactPoints float64
	_4v5w          int
	_4v5s          int
	_5v4w          int
	_5v4s          int
	pistols        int
	pistolsW       int
	saves          int
	clutches       int
	traded         int
	fass           int
	ef             int
	ud             int
	util           int
	ctR            int
	ctRW           int
	tR             int
	tRW            int
	deaths         int

	//kinda garbo
	normalizer int
}

type round struct {
	//round value
	roundNum            int8
	startingTick        int
	endingTick          int
	playerStats         map[uint64]*playerStats
	teamStats           map[string]*teamStats
	initTerroristCount  int
	initCTerroristCount int
	winnerClanName      string
	//winnerID            int //this is the unique ID which should not change BUT IT DOES
	winnerENUM         int //this effectively represents the side that won: 2 (T) or 3 (CT)
	integrityCheck     bool
	planter            uint64
	defuser            uint64
	serverNormalizer   int
	serverImpactPoints float64
}

type playerStats struct {
	name    string
	steamID uint64
	//teamID  int
	teamENUM     int
	teamClanName string
	side         int
	rounds       int
	//playerPoints float32
	//teamPoints float32
	damage         int
	kills          uint8
	assists        uint8
	deaths         uint8
	deathTick      int
	deathPlacement float64
	ticksAlive     int
	trades         int
	traded         int
	ok             int
	ol             int
	cl_1           int
	cl_2           int
	cl_3           int
	cl_4           int
	cl_5           int
	_2k            int
	_3k            int
	_4k            int
	_5k            int
	nadeDmg        int
	infernoDmg     int
	utilDmg        int
	ef             int
	fAss           int
	enemyFlashTime float64
	hs             int
	kastRounds     float64
	saves          int
	entries        int
	killPoints     float64
	impactPoints   float64
	winPoints      float64
	awpKills       int
	RF             int
	RA             int
	nadesThrown    int
	firesThrown    int
	flashThrown    int
	smokeThrown    int
	damageTaken    int
	suppRounds     int
	suppDamage     int
	rwk            int

	//derived
	utilThrown   int
	atd          int
	kast         float64
	killPointAvg float64
	iiwr         float64
	adr          float64
	drDiff       float64
	kR           float64
	tr           float64 //trade ratio
	impactRating float64
	rating       float64

	//side specific
	tDamage               int
	ctDamage              int
	tImpactPoints         float64
	tWinPoints            float64
	tOK                   int
	tOL                   int
	ctImpactPoints        float64
	ctWinPoints           float64
	ctOK                  int
	ctOL                  int
	tKills                uint8
	tDeaths               uint8
	tKAST                 float64
	tKASTRounds           float64
	tADR                  float64
	ctKills               uint8
	ctDeaths              uint8
	ctKAST                float64
	ctKASTRounds          float64
	ctADR                 float64
	tTeamsWinPoints       float64
	ctTeamsWinPoints      float64
	tWinPointsNormalizer  int
	ctWinPointsNormalizer int
	tRounds               int
	ctRounds              int
	ctRating              float64
	tRating               float64
	tADP                  float64
	ctADP                 float64

	tRF   int
	ctAWP int

	//kinda garbo
	teamsWinPoints      float64
	winPointsNormalizer int

	//"flags"
	health             int
	tradeList          map[uint64]int
	mostRecentFlasher  uint64
	mostRecentFlashVal float64
	damageList         map[uint64]int
}

func main() {
    input_dir := "in"
    files, _ := ioutil.ReadDir(input_dir)

    for _, file := range files {
        filename := file.Name()
        if strings.HasSuffix(filename, ".dem") {
            fmt.Println("processing", file.Name())
            processDemo(filepath.Join(input_dir, filename))
        }
    }
    
    var input string
	fmt.Scanln(&input)
}

func initGameObject() *game {
	g := game{}
	g.rounds = make([]*round, 0)
	g.potentialRound = &round{}

	g.flags.hasGameStarted = false
	g.flags.isGameLive = false
	g.flags.isGameOver = false
	g.flags.prePlant = true
	g.flags.postPlant = false
	g.flags.postWinCon = false
	//these three vars to check if we have a complete round
	g.flags.roundIntegrityStart = -1
	g.flags.roundIntegrityEnd = -1
	g.flags.roundIntegrityEndOfficial = -1

	return &g
}

func processDemo(demoName string) {

	game := initGameObject()

	/*
	   var errLog = "";
	*/

	f, err := os.Open(demoName)
	//f, err := os.Open("league1.dem")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	defer p.Close()

	//must parse header to get header info
	header, err := p.ParseHeader()
	if err != nil {
		panic(err)
	}

	//set map name
	game.mapName = strings.Title((header.MapName)[3:])
	fmt.Println("Map is", game.mapName)

	//set tick rate
	game.tickRate = int(math.Round(p.TickRate()))
	fmt.Println("Tick rate is", game.tickRate)

	game.tickLength = header.PlaybackTicks

	//creating a file to dump chat log into
	fmt.Printf("Creating chatLog file.\n")
	os.Mkdir("out", 0777)
	chatFile, chatFileErr := os.Create("out/chatLog.txt")
	if chatFileErr != nil {
		fmt.Printf("Error in opening chatLog file.\n")
	}
	defer chatFile.Close()

	fmt.Printf("Creating debug file.\n")
	debugFile, debugFileErr := os.Create("out/debug.txt")
	if debugFileErr != nil {
		fmt.Printf("Error in opening debug file.\n")
	}
	defer debugFile.Close()

	//---------------FUNCTIONS---------------

	initGameStart := func() {
		game.flags.hasGameStarted = true
		game.flags.isGameLive = true
		fmt.Println("GAME HAS STARTED!!!")

		game.teams = make(map[string]*team)

		teamTemp := p.GameState().TeamTerrorists()
		game.teams[teamTemp.ClanName()] = &team{name: validateTeamName(game, teamTemp.ClanName())}
		teamTemp = p.GameState().TeamCounterTerrorists()
		game.teams[teamTemp.ClanName()] = &team{name: validateTeamName(game, teamTemp.ClanName())}

		//to handle short and long matches
		if p.GameState().ConVars()["mp_maxrounds"] != "30" {
			maxRounds, fuckOFF := strconv.Atoi(p.GameState().ConVars()["mp_maxrounds"])
			if fuckOFF == nil {
				game.roundsToWin = maxRounds/2 + 1
			} else {
				//ADD TO ERROR LOG
				game.roundsToWin = maxRounds/2 + 1
				//maybe this gives us a way to check for short vs long match
			}
		} else {
			game.roundsToWin = 16 //we will assume long match in case convar is not set
		}

	}

	//reset various flags
	resetRoundFlags := func() {
		game.flags.prePlant = true
		game.flags.postPlant = false
		game.flags.postWinCon = false
		game.flags.tClutchVal = 0
		game.flags.ctClutchVal = 0
		game.flags.tClutchSteam = 0
		game.flags.ctClutchSteam = 0
		game.flags.tMoney = false
		game.flags.openingKill = true
	}

	initTeamPlayer := func(team *common.TeamState, currRoundObj *round) {
		for _, teamMember := range team.Members() {
			player := &playerStats{name: teamMember.Name, steamID: teamMember.SteamID64, side: int(team.Team()), teamENUM: team.ID(), teamClanName: validateTeamName(game, team.ClanName()), health: 100, tradeList: make(map[uint64]int), damageList: make(map[uint64]int)}
			currRoundObj.playerStats[player.steamID] = player
		}
	}

	initRound := func() {
		game.flags.roundIntegrityStart = p.GameState().TotalRoundsPlayed() + 1
		fmt.Println("We are starting round", game.flags.roundIntegrityStart)

		newRound := &round{roundNum: int8(game.flags.roundIntegrityStart), startingTick: p.GameState().IngameTick()}
		newRound.playerStats = make(map[uint64]*playerStats)
		newRound.teamStats = make(map[string]*teamStats)

		//set players in playerStats for the round
		terrorists := p.GameState().TeamTerrorists()
		counterTerrorists := p.GameState().TeamCounterTerrorists()

		initTeamPlayer(terrorists, newRound)
		initTeamPlayer(counterTerrorists, newRound)

		//set teams in teamStats for the round
		newRound.teamStats[validateTeamName(game, p.GameState().TeamTerrorists().ClanName())] = &teamStats{tR: 1}
		newRound.teamStats[validateTeamName(game, p.GameState().TeamCounterTerrorists().ClanName())] = &teamStats{ctR: 1}

		// Reset round
		game.potentialRound = newRound

		// if len(game.rounds) < game.flags.roundIntegrityStart {
		// 	//game.potentialRound = currRoundObj
		// } else {
		// 	//game.rounds[roundIntegrityStart - 1] = currRoundObj
		// 	fmt.Println(game.rounds[game.flags.roundIntegrityStart - 1].integrityCheck)
		// 	//game.rounds = append(game.rounds, currRoundObj)
		// }

		//track the number of people alive for clutch checking and record keeping
		game.flags.tAlive = len(terrorists.Members())
		game.flags.ctAlive = len(counterTerrorists.Members())
		game.potentialRound.initTerroristCount = game.flags.tAlive
		game.potentialRound.initCTerroristCount = game.flags.ctAlive

		resetRoundFlags()
	}

	processRoundOnWinCon := func(winnerClanName string) {
		game.flags.roundIntegrityEnd = p.GameState().TotalRoundsPlayed() + 1
		fmt.Println("We are processing round win con stuff", game.flags.roundIntegrityEnd)

		game.flags.prePlant = false
		game.flags.postPlant = false
		game.flags.postWinCon = true

		//set winner
		game.potentialRound.winnerClanName = winnerClanName
		game.teams[game.potentialRound.winnerClanName].score += 1
		//fmt.Println("We think this team won", game.teams[game.potentialRound.winnerID].name)
		//check clutch

	}

	processRoundFinal := func(lastRound bool) {
		game.potentialRound.endingTick = p.GameState().IngameTick()
		game.flags.roundIntegrityEndOfficial = p.GameState().TotalRoundsPlayed()
		if lastRound {
			game.flags.roundIntegrityEndOfficial += 1
			game.totalRounds = game.flags.roundIntegrityEndOfficial
		}
		fmt.Println("We are processing round final stuff", game.flags.roundIntegrityEndOfficial)
		fmt.Println(len(game.rounds))

		//we have the entire round uninterrupted
		if game.flags.roundIntegrityStart == game.flags.roundIntegrityEnd && game.flags.roundIntegrityEnd == game.flags.roundIntegrityEndOfficial {
			game.potentialRound.integrityCheck = true

			//check team stats
			if game.potentialRound.teamStats[game.potentialRound.winnerClanName].pistols == 1 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName].pistolsW = 1
			}
			if game.potentialRound.teamStats[game.potentialRound.winnerClanName]._4v5s == 1 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName]._4v5w = 1
			} else if game.potentialRound.teamStats[game.potentialRound.winnerClanName]._5v4s == 1 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName]._5v4w = 1
			}
			if game.potentialRound.teamStats[game.potentialRound.winnerClanName].tR == 1 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName].tRW = 1
			} else if game.potentialRound.teamStats[game.potentialRound.winnerClanName].ctR == 1 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName].ctRW = 1
			}

			//set the clutch
			if game.potentialRound.winnerENUM == 2 && game.flags.tClutchSteam != 0 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName].clutches = 1
				game.potentialRound.playerStats[game.flags.tClutchSteam].impactPoints += clutchBonus[game.flags.tClutchVal]
				switch game.flags.tClutchVal {
				case 1:
					game.potentialRound.playerStats[game.flags.tClutchSteam].cl_1 = 1
				case 2:
					game.potentialRound.playerStats[game.flags.tClutchSteam].cl_2 = 1
				case 3:
					game.potentialRound.playerStats[game.flags.tClutchSteam].cl_3 = 1
				case 4:
					game.potentialRound.playerStats[game.flags.tClutchSteam].cl_4 = 1
				case 5:
					game.potentialRound.playerStats[game.flags.tClutchSteam].cl_5 = 1
				}
			} else if game.potentialRound.winnerENUM == 3 && game.flags.ctClutchSteam != 0 {
				game.potentialRound.teamStats[game.potentialRound.winnerClanName].clutches = 1
				game.potentialRound.playerStats[game.flags.ctClutchSteam].impactPoints += clutchBonus[game.flags.ctClutchVal]
				switch game.flags.ctClutchVal {
				case 1:
					game.potentialRound.playerStats[game.flags.ctClutchSteam].cl_1 = 1
				case 2:
					game.potentialRound.playerStats[game.flags.ctClutchSteam].cl_2 = 1
				case 3:
					game.potentialRound.playerStats[game.flags.ctClutchSteam].cl_3 = 1
				case 4:
					game.potentialRound.playerStats[game.flags.ctClutchSteam].cl_4 = 1
				case 5:
					game.potentialRound.playerStats[game.flags.ctClutchSteam].cl_5 = 1
				}
			}

			//add multikills & saves & misc
			for _, player := range (game.potentialRound).playerStats {
				if player.deaths == 0 {
					player.kastRounds = 1
					if player.teamENUM != game.potentialRound.winnerENUM {
						player.saves = 1
						game.potentialRound.teamStats[player.teamClanName].saves = 1
					}
				}
				game.potentialRound.playerStats[player.steamID].impactPoints += player.killPoints
				game.potentialRound.playerStats[player.steamID].impactPoints += float64(player.damage) / float64(250)
				game.potentialRound.playerStats[player.steamID].impactPoints += multikillBonus[player.kills]

				switch player.kills {
				case 2:
					player._2k = 1
				case 3:
					player._3k = 1
				case 4:
					player._4k = 1
				case 5:
					player._5k = 1
				}

				if player.teamENUM == game.potentialRound.winnerENUM {
					player.winPoints = player.impactPoints
					player.RF = 1
					if player.name == "iNSANEmayne" {
						debugMsg := fmt.Sprintf("---------%.2f win points for brod on round %d.------------\n", player.impactPoints, game.potentialRound.roundNum)
						debugFile.WriteString(debugMsg)
						debugFile.Sync()
					}
				} else {
					player.RA = 1
				}
			}

			//add our valid round
			game.rounds = append(game.rounds, game.potentialRound)
		}

		//endRound function functionality

	}

	//-------------ALL OUR EVENTS---------------------

	p.RegisterEventHandler(func(e events.RoundStart) {
		fmt.Printf("Round Start\n")
	})

	p.RegisterEventHandler(func(e events.RoundFreezetimeEnd) {
		fmt.Printf("Round Freeze Time End\n")
		pistol := false

		//we are going to check to see if the first pistol is actually starting
		membersT := p.GameState().TeamTerrorists().Members()
		membersCT := p.GameState().TeamCounterTerrorists().Members()
		if len(membersT) != 0 && len(membersCT) != 0 {
			if membersT[0].Money()+membersT[0].MoneySpentThisRound() == 800 && membersCT[0].Money()+membersCT[0].MoneySpentThisRound() == 800 {
				//start the game
				if !game.flags.hasGameStarted {
					initGameStart()
				}

				//track the pistol
				pistol = true
			}
		}
		fmt.Println("Has the Game Started?", game.flags.hasGameStarted)

		if game.flags.isGameLive {
			//init round stats
			initRound()
			if pistol {
				for _, team := range game.potentialRound.teamStats {
					team.pistols = 1
				}
			}

		}

	})

	p.RegisterEventHandler(func(e events.RoundEnd) {
		fmt.Println("Round", p.GameState().TotalRoundsPlayed()+1, "End", e.WinnerState.ClanName(), "won", "this determined from e.WinnerState.ClanName()")

		fmt.Println("e.WinnerState.ID()", e.WinnerState.ID(), "and", "e.Winner", e.Winner, "and", "e.WinnerState.Team()", e.WinnerState.Team())

		validWinner := true
		if e.Winner < 2 {
			validWinner = false
			//and set the integrity flag to false

		} else if e.Winner == 2 {
			game.flags.tMoney = true
		} else {
			//we need to check if the game is over

		}

		//we want to actually process the round
		if game.flags.isGameLive && validWinner && game.flags.roundIntegrityStart == p.GameState().TotalRoundsPlayed()+1 {
			game.potentialRound.winnerENUM = int(e.Winner)
			processRoundOnWinCon(validateTeamName(game, e.WinnerState.ClanName()))

			//check last round
			roundWinnerScore := game.teams[validateTeamName(game, e.WinnerState.ClanName())].score
			roundLoserScore := game.teams[validateTeamName(game, e.LoserState.ClanName())].score
			fmt.Println("winner Rounds", roundWinnerScore)
			fmt.Println("loser Rounds", roundLoserScore)

			if game.roundsToWin == 16 {
				//check for normal win
				if roundWinnerScore == 16 && roundLoserScore < 15 {
					//normal win
					game.winnerClanName = game.potentialRound.winnerClanName
					processRoundFinal(true)
				} else if roundWinnerScore > 15 { //check for OT win
					overtime := ((roundWinnerScore+roundLoserScore)-30-1)/6 + 1
					//OT win
					if (roundWinnerScore-15-1)/3 == overtime {
						game.winnerClanName = game.potentialRound.winnerClanName
						processRoundFinal(true)
					}
				}
			} else if game.roundsToWin == 9 {
				//check for normal win
				if roundWinnerScore == 9 && roundLoserScore < 8 {
					//normal win
					game.winnerClanName = game.potentialRound.winnerClanName
					processRoundFinal(true)
				} else if roundWinnerScore == 8 && roundLoserScore == 8 { //check for tie
					//tie
					game.winnerClanName = game.potentialRound.winnerClanName
					processRoundFinal(true)
				}
			}
		}

		//check last round
		//or check overtime win

	})

	//round end official doesnt fire on the last round
	p.RegisterEventHandler(func(e events.RoundEndOfficial) {
		fmt.Printf("Round End Official\n")

		if game.flags.isGameLive && game.flags.roundIntegrityEnd == p.GameState().TotalRoundsPlayed() {
			processRoundFinal(false)
		}
	})

	// Register handler on kill events
	p.RegisterEventHandler(func(e events.Kill) {
		if game.flags.isGameLive && isDuringExpectedRound(game, p) {
			pS := game.potentialRound.playerStats
			tick := p.GameState().IngameTick()

			killerExists := false
			victimExists := false
			assisterExists := false
			if e.Killer != nil && pS[e.Killer.SteamID64] != nil {
				killerExists = true
			}
			if e.Victim != nil && pS[e.Victim.SteamID64] != nil {
				victimExists = true
			}
			if e.Assister != nil && pS[e.Assister.SteamID64] != nil {
				assisterExists = true
			}

			killValue := 1.0
			multiplier := 1.0
			traded := false
			assisted := false
			flashAssisted := false

			//death logic (traded here)
			if victimExists {
				pS[e.Victim.SteamID64].deaths += 1
				pS[e.Victim.SteamID64].deathTick = tick
				if e.Victim.Team == 2 {
					game.flags.tAlive -= 1
					pS[e.Victim.SteamID64].deathPlacement = float64(game.potentialRound.initTerroristCount - game.flags.tAlive)
					//pS[e.Victim.SteamID64].tADP = float64(game.potentialRound.initTerroristCount - game.flags.tAlive)
				} else if e.Victim.Team == 3 {
					game.flags.ctAlive -= 1
					pS[e.Victim.SteamID64].deathPlacement = float64(game.potentialRound.initCTerroristCount - game.flags.ctAlive)
					//pS[e.Victim.SteamID64].ctADP = float64(game.potentialRound.initCTerroristCount - game.flags.ctAlive)
				} else {
					//else log an error
				}

				//do 4v5 calc
				if game.flags.openingKill && game.potentialRound.initCTerroristCount+game.potentialRound.initTerroristCount == 10 {
					//the 10th player died
					_4v5Team := pS[e.Victim.SteamID64].teamClanName
					game.potentialRound.teamStats[_4v5Team]._4v5s = 1
					for teamName, team := range game.potentialRound.teamStats {
						if teamName != _4v5Team {
							team._5v4s = 1
						}
					}
				}

				//add support damage
				for suppSteam, suppDMG := range pS[e.Victim.SteamID64].damageList {
					if killerExists && suppSteam != e.Killer.SteamID64 {
						pS[suppSteam].suppDamage += suppDMG
						if pS[suppSteam].suppDamage > 60 {
							pS[suppSteam].suppRounds = 1
						}
					} else if !killerExists {
						pS[suppSteam].suppDamage += suppDMG
						if pS[suppSteam].suppDamage > 60 {
							pS[suppSteam].suppRounds = 1
						}
					}

				}

				//check clutch start

				if !game.flags.postWinCon {
					if game.flags.tAlive == 1 && game.flags.tClutchVal == 0 {
						game.flags.tClutchVal = game.flags.ctAlive
						membersT := p.GameState().TeamTerrorists().Members()
						for _, terrorist := range membersT {
							if terrorist.IsAlive() && e.Victim.SteamID64 != terrorist.SteamID64 {
								game.flags.tClutchSteam = terrorist.SteamID64
								fmt.Println("Clutch opportunity:", terrorist.Name, game.flags.tClutchVal)
							}
						}
					}
					if game.flags.ctAlive == 1 && game.flags.ctClutchVal == 0 {
						game.flags.ctClutchVal = game.flags.tAlive
						membersCT := p.GameState().TeamCounterTerrorists().Members()
						for _, counterTerrorist := range membersCT {
							if counterTerrorist.IsAlive() && e.Victim.SteamID64 != counterTerrorist.SteamID64 {
								game.flags.ctClutchSteam = counterTerrorist.SteamID64
								fmt.Println("Clutch opportunity:", counterTerrorist.Name, game.flags.ctClutchVal)
							}
						}
					}
				}

				pS[e.Victim.SteamID64].ticksAlive = tick - game.potentialRound.startingTick
				for deadGuySteam, deadTick := range (*game.potentialRound).playerStats[e.Victim.SteamID64].tradeList {
					if tick-deadTick < tradeCutoff*game.tickRate {
						pS[deadGuySteam].traded = 1
						pS[deadGuySteam].kastRounds = 1
					}
				}
			}

			//assist logic
			if assisterExists && victimExists && e.Assister.TeamState.ID() != e.Victim.TeamState.ID() {
				//this logic needs to be replaced
				pS[e.Assister.SteamID64].assists += 1
				pS[e.Assister.SteamID64].kastRounds = 1
				pS[e.Assister.SteamID64].suppRounds = 1
				assisted = true
				if e.AssistedFlash {
					pS[e.Assister.SteamID64].fAss += 1
					flashAssisted = true
				} else if float64(p.GameState().IngameTick()) < pS[e.Victim.SteamID64].mostRecentFlashVal {
					//this will trigger if there is both a flash assist and a damage assist
					pS[pS[e.Victim.SteamID64].mostRecentFlasher].fAss += 1
					pS[pS[e.Victim.SteamID64].mostRecentFlasher].suppRounds = 1
					flashAssisted = true
				}

			}

			//kill logic (trades here)
			if killerExists && victimExists && e.Killer.TeamState.ID() != e.Victim.TeamState.ID() {
				pS[e.Killer.SteamID64].kills += 1
				pS[e.Killer.SteamID64].kastRounds = 1
				pS[e.Killer.SteamID64].rwk = 1
				pS[e.Killer.SteamID64].tradeList[e.Victim.SteamID64] = tick
				if e.Weapon.Type == 309 {
					pS[e.Killer.SteamID64].awpKills += 1
					if e.Killer.Team == 3 {
						pS[e.Killer.SteamID64].ctAWP += 1
					}
				}
				if e.IsHeadshot {
					pS[e.Killer.SteamID64].hs += 1
				}
				for _, deadTick := range (*game.potentialRound).playerStats[e.Victim.SteamID64].tradeList {
					if tick-deadTick < tradeCutoff*game.tickRate {
						pS[e.Killer.SteamID64].trades += 1
						traded = true
						break
					}
				}

				killerTeam := e.Killer.Team
				if game.flags.prePlant {
					//normal base value
					if killerTeam == 2 {
						//taking site by T
						killValue = 1.2
					} else if killerTeam == 3 {
						//site Defense by CT
						killValue = 1
					}
				} else if game.flags.postPlant {
					//site D or retake
					if killerTeam == 2 {
						//site Defense by T
						killValue = 1
					} else if killerTeam == 3 {
						//retake
						killValue = 1.2
					}
				} else if game.flags.postWinCon {
					//exit or chase
					if game.potentialRound.winnerENUM == 2 { //Ts win
						if killerTeam == 2 { //chase
							killValue = 0.8
						}
						if killerTeam == 3 { //exit
							killValue = 0.6
						}
					} else if game.potentialRound.winnerENUM == 3 { //CTs win
						if killerTeam == 2 { //T kill in lost round
							killValue = 0.5
						}
						if killerTeam == 3 { //CT kill in won round
							if game.flags.tMoney {
								killValue = 0.6
							} else {
								killValue = 0.8
							}
						}
					}
				}

				if game.flags.openingKill {
					game.flags.openingKill = false

					pS[e.Killer.SteamID64].ok = 1
					pS[e.Victim.SteamID64].ol = 1

					if killerTeam == 2 { //T entry/opener {
						if game.flags.prePlant {
							multiplier += 0.8
							pS[e.Killer.SteamID64].entries = 1
						} else {
							multiplier += 0.3
						}
					} else if killerTeam == 3 { //CT opener
						multiplier += 0.5
					}

				} else if traded {
					multiplier += 0.3
				}

				if flashAssisted { //flash assisted kill
					multiplier += 0.2
				}
				if assisted { //assisted kill
					killValue -= 0.15
					pS[e.Assister.SteamID64].impactPoints += 0.15
				}

				killValue *= multiplier

				ecoRatio := float64(e.Victim.EquipmentValueCurrent()) / float64(e.Killer.EquipmentValueCurrent())
				ecoMod := 1.0
				if ecoRatio > 4 {
					ecoMod += 0.25
				} else if ecoRatio > 2 {
					ecoMod += 0.14
				} else if ecoRatio < 0.25 {
					ecoMod -= 0.25
				} else if ecoRatio < 0.5 {
					ecoMod -= 0.14
				}
				killValue *= ecoMod

				pS[e.Killer.SteamID64].killPoints += killValue
				if e.Killer.Name == "iNSANEmayne" {
					debugMsg := fmt.Sprintf("---------%.2f kill points for brod on round %d.------------\n", killValue, game.potentialRound.roundNum)
					debugFile.WriteString(debugMsg)
					debugFile.Sync()
				}

			}

		}
		var hs string
		if e.IsHeadshot {
			hs = " (HS)"
		}
		var wallBang string
		if e.PenetratedObjects > 0 {
			wallBang = " (WB)"
		}
		fmt.Printf("%s <%v%s%s> %s at %d\n", e.Killer, e.Weapon, hs, wallBang, e.Victim, p.GameState().IngameTick())
	})

	p.RegisterEventHandler(func(e events.PlayerHurt) {
		//fmt.Printf("Player Hurt\n")
		if game.flags.isGameLive {
			equipment := e.Weapon.Type
			if e.Player != nil {
				game.potentialRound.playerStats[e.Player.SteamID64].damageTaken += e.HealthDamageTaken
			}
			if e.Player != nil && e.Attacker != nil && e.Player.Team != e.Attacker.Team {
				game.potentialRound.playerStats[e.Attacker.SteamID64].damage += e.HealthDamageTaken

				//add to damage list for supp damage calc
				game.potentialRound.playerStats[e.Player.SteamID64].damageList[e.Attacker.SteamID64] += e.HealthDamageTaken

				if equipment >= 500 && equipment <= 506 {
					game.potentialRound.playerStats[e.Attacker.SteamID64].utilDmg += e.HealthDamageTaken
					if equipment == 506 {
						game.potentialRound.playerStats[e.Attacker.SteamID64].nadeDmg += e.HealthDamageTaken
					}
					if equipment == 502 || equipment == 503 {
						game.potentialRound.playerStats[e.Attacker.SteamID64].infernoDmg += e.HealthDamageTaken
					}
				}
			}
		}
	})

	p.RegisterEventHandler(func(e events.PlayerFlashed) {
		//fmt.Printf("Player Flashed\n")
		tick := float64(p.GameState().IngameTick())
		blindTicks := e.FlashDuration().Seconds() * 128.0
		if game.flags.isGameLive && e.Player != nil && e.Attacker != nil {
			victim := e.Player
			flasher := e.Attacker
			if flasher.Team != victim.Team && blindTicks > 128.0 && victim.IsAlive() && (float64(victim.FlashDuration) < (blindTicks/128.0 + 1)) {
				game.potentialRound.playerStats[flasher.SteamID64].ef += 1
				game.potentialRound.playerStats[flasher.SteamID64].enemyFlashTime += (blindTicks / 128.0)
				if tick+blindTicks > game.potentialRound.playerStats[victim.SteamID64].mostRecentFlashVal {
					game.potentialRound.playerStats[victim.SteamID64].mostRecentFlashVal = tick + blindTicks
					game.potentialRound.playerStats[victim.SteamID64].mostRecentFlasher = flasher.SteamID64
				}

			}
			// if flasher.Name != "" {
			// 	debugMsg := fmt.Sprintf("%s flashed %s for %.2f at %d. He was %f blind.\n", flasher, victim, blindTicks/128, int(tick), victim.FlashDuration)
			// 	debugFile.WriteString(debugMsg)
			// 	debugFile.Sync()
			// }

		}
	})

	p.RegisterEventHandler(func(e events.PlayerJump) {
		//fmt.Printf("Player Jumped\n")
	})

	p.RegisterEventHandler(func(e events.BombPlanted) {
		fmt.Printf("Bomb Planted\n")
		if game.flags.isGameLive && !game.flags.postWinCon {
			game.flags.prePlant = false
			game.flags.postPlant = true
			game.flags.tMoney = true
			game.potentialRound.planter = e.BombEvent.Player.SteamID64
		}
	})

	p.RegisterEventHandler(func(e events.BombDefused) {
		fmt.Println("Bomb Defused by", e.BombEvent.Player.Name)
		if game.flags.isGameLive {
			game.flags.prePlant = false
			game.flags.postPlant = false
			game.flags.postWinCon = true
			game.potentialRound.playerStats[e.BombEvent.Player.SteamID64].impactPoints += 0.5
		}
	})

	p.RegisterEventHandler(func(e events.BombExplode) {
		fmt.Printf("Bomb Exploded\n")
		if game.flags.isGameLive {
			game.flags.prePlant = false
			game.flags.postPlant = false
			game.flags.postWinCon = true
			game.potentialRound.playerStats[game.potentialRound.planter].impactPoints += 0.5
		}
	})

	p.RegisterEventHandler(func(e events.GrenadeProjectileThrow) {
		//fmt.Println("Grenade Thrown", e.Projectile.WeaponInstance.Type)
		if game.flags.isGameLive {
			if e.Projectile.WeaponInstance.Type == 506 {
				game.potentialRound.playerStats[e.Projectile.Thrower.SteamID64].nadesThrown += 1
			} else if e.Projectile.WeaponInstance.Type == 505 {
				game.potentialRound.playerStats[e.Projectile.Thrower.SteamID64].smokeThrown += 1
			} else if e.Projectile.WeaponInstance.Type == 504 {
				game.potentialRound.playerStats[e.Projectile.Thrower.SteamID64].flashThrown += 1
			} else if e.Projectile.WeaponInstance.Type == 502 || e.Projectile.WeaponInstance.Type == 503 {
				game.potentialRound.playerStats[e.Projectile.Thrower.SteamID64].firesThrown += 1
			}

		}
	})

	p.RegisterEventHandler(func(e events.PlayerDisconnected) {
		//fmt.Println("Player DC", e.Player)

		//update alive players
		if game.flags.isGameLive {
			game.flags.tAlive = 0
			game.flags.ctAlive = 0

			membersT := p.GameState().TeamTerrorists().Members()
			for _, terrorist := range membersT {
				if terrorist.IsAlive() {
					game.flags.tAlive += 1
				}
			}
			membersCT := p.GameState().TeamCounterTerrorists().Members()
			for _, counterTerrorist := range membersCT {
				if counterTerrorist.IsAlive() {
					game.flags.ctAlive += 1
				}
			}
		}

	})

	if printChatLog {
		p.RegisterEventHandler(func(e events.ChatMessage) {
			chatMsg := fmt.Sprintf("%s: \"%s\" at %d\n", e.Sender, e.Text, p.GameState().IngameTick())
			chatFile.WriteString(chatMsg)
			chatFile.Sync()
			//check(err)
		})
	}

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}

	//----END OF MATCH PROCESSING----
	//we want to iterate through rounds backwards to make sure their are no repeats

	endOfMatchProcessing(game)

	fmt.Println("Demo is complete!")
	//cleanup()

}
