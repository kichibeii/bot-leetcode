package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

type Stat struct {
	QuestionID         int    `json:"question_id"`
	QuestionTitle      string `json:"question__title"`
	QuestionTitleSlug  string `json:"question__title_slug"`
	TotalACS           int    `json:"total_acs"`
	TotalSubmitted     int    `json:"total_submitted"`
	FrontendQuestionID int    `json:"frontend_question_id"`
	IsNewQuestion      bool   `json:"is_new_question"`
}

// difficulty
const (
	DifficultyEasy   = 1
	DifficultyMedium = 2
	DifficultyHard   = 3
)

func main() {
	s := gocron.NewScheduler(time.UTC)
	_, err := s.Every(1).Hour().Do(func() {
		loopFunction()
	})
	if err != nil {
		panic(err)
	}

	s.StartAsync()

	// start the select without parameter for
	// continous running
	select {}
}

// main job of function
func loopFunction() {
	// random difficulty
	rand.Seed(time.Now().UnixNano())
	randomDifficulty := rand.Intn(3) + 1

	// open file
	jsonFile, err := os.ReadFile(fmt.Sprintf("question-%d.json", randomDifficulty))
	if err != nil {
		panic(err)
	}

	var listQuestion []Stat
	err = json.Unmarshal(jsonFile, &listQuestion)
	if err != nil {
		panic(err)
	}

	// pick random data
	numberCase := rand.Intn(len(listQuestion))

	question := listQuestion[numberCase]

	now := time.Now()
	dateString := now.Format("2006-01-02")

	link := "https://leetcode.com/problems/" + question.QuestionTitleSlug

	difficultyString := ""
	switch randomDifficulty {
	case 1:
		difficultyString = "Easy"
	case 2:
		difficultyString = "Medium"
	case 3:
		difficultyString = "Hard"
	}

	message := fmt.Sprintf("Date : %s \nLink : %s \nDifficulty  : %s \n \n", dateString, link, difficultyString)

	// post to discord
	fileToken, err := os.Open("token.txt")
	if err != nil {
		PrintError("open file token", err)
		return
	}
	defer fileToken.Close()

	dataToken, err := ioutil.ReadAll(fileToken)
	if err != nil {
		PrintError("read token file", err)
		return
	}

	token := string(dataToken)
	token = strings.TrimSpace(token)
	keyAuth := "Bot " + token

	discord, err := discordgo.New(keyAuth)
	if err != nil {
		PrintError("new discord", err)
		return
	}

	err = discord.Open()
	if err != nil {
		PrintError("discord open", err)
		return
	}
	defer discord.Close()

	channelID := "1165116574290673785"
	_, err = discord.ChannelMessageSend(channelID, message)
	if err != nil {
		PrintError("send message", err)
		return
	}
}

func PrintError(msg string, err error) {
	fmt.Println(msg)
	fmt.Println(err)
}
