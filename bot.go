package main

import (
	"bufio"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

var bot,e=tgbot.NewBotAPI(getAPI())

func getAPI() string {
	var res string
	file, err := os.Open("API.txt")
	checkErr(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = scanner.Text()
	}
	checkErr(scanner.Err())
	return strings.TrimSpace(res)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	checkErr(e)

	u:=tgbot.NewUpdate(0)
	u.Timeout= 60

	updates:=bot.GetUpdatesChan(u)

	for update:=range updates{
		msg:=tgbot.NewMessage(update.Message.Chat.ID,"")

		if !update.Message.IsCommand(){
			continue
		}
		switch strings.ToLower(update.Message.Command()){
		case "start":
			msg.Text="Привет"
		case "go":
			msg.Text="Go"
		default:
			msg.Text="Я не знаю такую команду..."
		}
		_,err:=bot.Send(msg)
		checkErr(err)
	}
}
