package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	//dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
)

func beginOutput(game *game) {

	setDisplayOrders(game)

	m_ID := createHash(game)
	fmt.Println("M_ID", m_ID)

	csvName := "out/" + m_ID + ".csv"

	outputFile, outputFileErr := os.Create(csvName)
	if outputFileErr != nil {
		fmt.Println("OH NOE")
	}
	w := csv.NewWriter(outputFile)

	records := [][]string{
		{
			"m_ID",
			"Map",
			"Team",
			"steam",
			"Name",
			"Rating",
			"Kills",
			"Assists",
			"Deaths",
			"ADR",
			"KAST",
			"Impact",
			"CT",
			"T",
			"ADP",
			"SuppR",
			"SuppX",
			"UD",
			"EF",
			"F_Ass",
			"Util",
			"HS",
			"AWP_K",
			"F_Kills",
			"F_Deaths",
			"Entries",
			"Saves",
			"Trades",
			"Traded",
			"2k",
			"3k",
			"4k",
			"5k",
			"1v1",
			"1v2",
			"1v3",
			"1v4",
			"1v5",
			"Rounds",
			"RF",
			"RA",
			"Damage",
			"XTaken",
			"ATD",
			"ADP-CT",
			"ADP-T",
			"Smokes",
			"Flashes",
			"Fires",
			"Nades",
			"FireX",
			"NadeX",
			"EFT",
			"RWK",
			"IWR",
			"KPA",
			"tOL",
			"ctOK",
			"ctOL",
			"tRounds",
			"tRF",
			"ctAWP",
			"ctK",
		},
	}

	//this shit needa die
	teamA := []string{m_ID, game.mapName, "", "", game.teams[game.teamOrder[0]].name, "", "1", "", "", "", "", "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].ctRW), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].tRW), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].ctR), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].tR), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].ud), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].ef), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].fass), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].util), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].pistolsW), "", "", "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].saves), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].traded), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].deaths), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name]._4v5w), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name]._5v4w), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name].clutches), "", "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name]._4v5s), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[0]].name]._5v4s), strconv.Itoa(game.totalRounds), strconv.Itoa(game.teams[game.teamOrder[0]].score), strconv.Itoa(game.totalRounds - game.teams[game.teamOrder[0]].score)}
	records = append(records, [][]string{teamA}...)

	for _, steam := range game.playerOrder {
		player := game.totalPlayerStats[steam]
		if player.teamClanName == game.teamOrder[0] {
			playerOutput := []string{
				m_ID,
				game.mapName,
				game.teams[game.teamOrder[0]].name,
				"sid" + strconv.FormatUint(player.steamID, 10),
				player.name,
				fmt.Sprintf("%.2f", player.rating),
				strconv.Itoa(int(player.kills)),
				strconv.Itoa(int(player.assists)),
				strconv.Itoa(int(player.deaths)),
				strconv.Itoa(int(player.adr)),
				fmt.Sprintf("%.2f", player.kast),
				fmt.Sprintf("%.2f", player.impactRating),
				fmt.Sprintf("%.2f", player.ctRating),
				fmt.Sprintf("%.2f", player.tRating),
				fmt.Sprintf("%.2f", player.deathPlacement),
				strconv.Itoa(int(player.suppRounds)),
				strconv.Itoa(int(player.suppDamage)),
				strconv.Itoa(int(player.utilDmg)),
				strconv.Itoa(int(player.ef)),
				strconv.Itoa(int(player.fAss)),
				strconv.Itoa(int(player.utilThrown)),
				strconv.Itoa(int(player.hs)),
				strconv.Itoa(int(player.awpKills)),
				strconv.Itoa(int(player.ok)),
				strconv.Itoa(int(player.ol)),
				strconv.Itoa(int(player.entries)),
				strconv.Itoa(int(player.saves)),
				strconv.Itoa(int(player.trades)),
				strconv.Itoa(int(player.traded)),
				strconv.Itoa(int(player._2k)),
				strconv.Itoa(int(player._3k)),
				strconv.Itoa(int(player._4k)),
				strconv.Itoa(int(player._5k)),
				strconv.Itoa(int(player.cl_1)),
				strconv.Itoa(int(player.cl_2)),
				strconv.Itoa(int(player.cl_3)),
				strconv.Itoa(int(player.cl_4)),
				strconv.Itoa(int(player.cl_5)),
				strconv.Itoa(int(player.rounds)),
				strconv.Itoa(int(player.RF)),
				strconv.Itoa(int(player.RA)),
				strconv.Itoa(int(player.damage)),
				strconv.Itoa(int(player.damageTaken)),
				strconv.Itoa(int(player.atd)),
				fmt.Sprintf("%.2f", player.ctADP),
				fmt.Sprintf("%.2f", player.tADP),
				strconv.Itoa(int(player.smokeThrown)),
				strconv.Itoa(int(player.flashThrown)),
				strconv.Itoa(int(player.firesThrown)),
				strconv.Itoa(int(player.nadesThrown)),
				strconv.Itoa(int(player.infernoDmg)),
				strconv.Itoa(int(player.nadeDmg)),
				fmt.Sprintf("%.0f", player.enemyFlashTime),
				strconv.Itoa(int(player.rwk)),
				fmt.Sprintf("%.2f", player.iiwr),
				fmt.Sprintf("%.2f", player.killPointAvg),
				strconv.Itoa(int(player.tOL)),
				strconv.Itoa(int(player.ctOK)),
				strconv.Itoa(int(player.ctOL)),
				strconv.Itoa(int(player.tRounds)),
				strconv.Itoa(int(player.tRF)),
				strconv.Itoa(int(player.ctAWP)),
				strconv.Itoa(int(player.ctKills)),
			}
			records = append(records, [][]string{playerOutput}...)
		}
	}

	teamB := []string{m_ID, game.mapName, "", "", game.teams[game.teamOrder[1]].name, "", "", "", "1", "", "", "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].ctRW), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].tRW), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].ctR), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].tR), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].ud), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].ef), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].fass), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].util), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].pistolsW), "", "", "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].saves), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].traded), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].deaths), "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name]._4v5w), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name]._5v4w), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name].clutches), "", "", strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name]._4v5s), strconv.Itoa(game.totalTeamStats[game.teams[game.teamOrder[1]].name]._5v4s), strconv.Itoa(game.totalRounds), strconv.Itoa(game.teams[game.teamOrder[1]].score), strconv.Itoa(game.totalRounds - game.teams[game.teamOrder[1]].score)}
	records = append(records, [][]string{teamB}...)

	for _, steam := range game.playerOrder {
		player := game.totalPlayerStats[steam]
		if player.teamClanName == game.teamOrder[1] {
			playerOutput := []string{
				m_ID,
				game.mapName,
				game.teams[game.teamOrder[1]].name,
				"sid" + strconv.FormatUint(player.steamID, 10),
				player.name,
				fmt.Sprintf("%.2f", player.rating),
				strconv.Itoa(int(player.kills)),
				strconv.Itoa(int(player.assists)),
				strconv.Itoa(int(player.deaths)),
				strconv.Itoa(int(player.adr)),
				fmt.Sprintf("%.2f", player.kast),
				fmt.Sprintf("%.2f", player.impactRating),
				fmt.Sprintf("%.2f", player.ctRating),
				fmt.Sprintf("%.2f", player.tRating),
				fmt.Sprintf("%.2f", player.deathPlacement),
				strconv.Itoa(int(player.suppRounds)),
				strconv.Itoa(int(player.suppDamage)),
				strconv.Itoa(int(player.utilDmg)),
				strconv.Itoa(int(player.ef)),
				strconv.Itoa(int(player.fAss)),
				strconv.Itoa(int(player.utilThrown)),
				strconv.Itoa(int(player.hs)),
				strconv.Itoa(int(player.awpKills)),
				strconv.Itoa(int(player.ok)),
				strconv.Itoa(int(player.ol)),
				strconv.Itoa(int(player.entries)),
				strconv.Itoa(int(player.saves)),
				strconv.Itoa(int(player.trades)),
				strconv.Itoa(int(player.traded)),
				strconv.Itoa(int(player._2k)),
				strconv.Itoa(int(player._3k)),
				strconv.Itoa(int(player._4k)),
				strconv.Itoa(int(player._5k)),
				strconv.Itoa(int(player.cl_1)),
				strconv.Itoa(int(player.cl_2)),
				strconv.Itoa(int(player.cl_3)),
				strconv.Itoa(int(player.cl_4)),
				strconv.Itoa(int(player.cl_5)),
				strconv.Itoa(int(player.rounds)),
				strconv.Itoa(int(player.RF)),
				strconv.Itoa(int(player.RA)),
				strconv.Itoa(int(player.damage)),
				strconv.Itoa(int(player.damageTaken)),
				strconv.Itoa(int(player.atd)),
				fmt.Sprintf("%.2f", player.ctADP),
				fmt.Sprintf("%.2f", player.tADP),
				strconv.Itoa(int(player.smokeThrown)),
				strconv.Itoa(int(player.flashThrown)),
				strconv.Itoa(int(player.firesThrown)),
				strconv.Itoa(int(player.nadesThrown)),
				strconv.Itoa(int(player.infernoDmg)),
				strconv.Itoa(int(player.nadeDmg)),
				fmt.Sprintf("%.0f", player.enemyFlashTime),
				strconv.Itoa(int(player.rwk)),
				fmt.Sprintf("%.2f", player.iiwr),
				fmt.Sprintf("%.2f", player.killPointAvg),
				strconv.Itoa(int(player.tOL)),
				strconv.Itoa(int(player.ctOK)),
				strconv.Itoa(int(player.ctOL)),
				strconv.Itoa(int(player.tRounds)),
				strconv.Itoa(int(player.tRF)),
				strconv.Itoa(int(player.ctAWP)),
				strconv.Itoa(int(player.ctKills)),
			}
			records = append(records, [][]string{playerOutput}...)
		}
	}

	//this shit needa die
	result := "" + game.teams[game.teamOrder[0]].name + " " + strconv.Itoa(game.teams[game.teamOrder[0]].score) + " - " + strconv.Itoa(game.teams[game.teamOrder[1]].score) + " " + game.teams[game.teamOrder[1]].name
	resultLine := []string{"1", game.mapName, result}
	records = append(records, [][]string{resultLine}...)

	for i, _ := range records {
		w.Write(records[i])
	}
	w.Flush()
}

func createHash(game *game) string {
	fmt.Println("headerTickL", game.tickLength)
	hashValue := fmt.Sprint(game.tickLength)
	totalDamage := 0
	totalUD := 0
	playerInitial := ""

	for _, player := range game.totalPlayerStats {
		totalDamage += player.damage
		totalUD += player.utilDmg
		playerInitial += string(player.name[0])
	}

	s := strings.Split(playerInitial, "")
	sort.Strings(s)
	playerInitial = strings.Join(s, "")

	fmt.Println("tick", hashValue)
	hashValue += fmt.Sprint(totalDamage) + playerInitial

	return randomizeHash(hashValue, totalUD)
}

func randomizeHash(hashValue string, seedVal int) string {
	rand.Seed(int64(seedVal))

	hashValueRune := []rune(hashValue)
	rand.Shuffle(len(hashValueRune), func(i, j int) {
		hashValueRune[i], hashValueRune[j] = hashValueRune[j], hashValueRune[i]
	})

	return string(hashValueRune)
}

func setDisplayOrders(game *game) {
	if game.winnerClanName != "" {
		game.teamOrder = append(game.teamOrder, game.winnerClanName)
		for _, team := range game.teams {
			if game.teamOrder[0] != team.name {
				game.teamOrder = append(game.teamOrder, team.name)
			}
		}
	} else {
		//just sort alphabetically
		for teamID, _ := range game.teams {
			if len(game.teamOrder) == 0 {
				game.teamOrder = append(game.teamOrder, teamID)
			} else {
				if game.teams[game.teamOrder[0]].name < game.teams[teamID].name {
					game.teamOrder = append(game.teamOrder, teamID)
				} else {
					game.teamOrder = append(game.teamOrder, game.teamOrder[0])
					game.teamOrder[0] = teamID
				}
			}
		}
	}

	for _, teamClanName := range game.teamOrder {
		offset := len(game.playerOrder)
		for steam, player := range game.totalPlayerStats {
			if player.teamClanName == teamClanName {
				if len(game.playerOrder) > offset {
					//subsetI := len(game.playerOrder) - offset
					for index, _ := range game.playerOrder[offset:] {
						if player.rating > game.totalPlayerStats[game.playerOrder[index+offset]].rating {
							game.playerOrder = append(game.playerOrder[:index+offset+1], game.playerOrder[index+offset:]...)
							game.playerOrder[index+offset] = steam
							break
						} else if (index+offset)+1 == len(game.playerOrder) {
							game.playerOrder = append(game.playerOrder, steam)
							break
						} else {
							continue
						}
					}
				} else {
					game.playerOrder = append(game.playerOrder, steam)
				}
			}
		}
	}
	fmt.Println(game.teamOrder)
	fmt.Println(game.playerOrder)
}
