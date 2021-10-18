package main

import (
	"bufio"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var bot, e = tgbot.NewBotAPI(getAPI())
var quests = createQuest()
var isGame = false
var questObject questions

type questions struct {
	quest   string
	answer1 string
	answer2 string
	answer3 string
	want    string
}

func createQuest() []questions {
	slice := []questions{
		{quest: "Кто из президентов США написал свой собственный рассказ про Шерлока Холмса?", answer1: "Джон Кеннеди", answer2: "Франклин Рузвельт", answer3: "Рональд Рейган", want: "Франклин Рузвельт"},
		{quest: "Какую пошлину ввели в XII  веке в Англии для того чтобы заставить мужчин пойти на войну?", answer1: "Налог на тунеядство", answer2: "Налог на трусость", answer3: "Налог на отсутствие сапог", want: "Налог на трусость"},
		{quest: "Откуда пошло выражение «деньги не пахнут?", answer1: "От подателей за провоз парфюмерии", answer2: "От сборов за нестиранные носки", answer3: "От налога на туалеты", want: "От налога на туалеты"},
		{quest: "Туристы, приезжающие на Майорку, обязаны заплатить налог…", answer1: "На плавки", answer2: "На пальмы", answer3: "На солнце", want: "На солнце"},
		{quest: "Один известный писатель рассказывал, что списал образ старушки-вредины со своей бывшей жены. При этом бабулька оказалась удивительно похожей на Коко Шанель. На голове у нее всегда была шляпка со складной тульей, благодаря которой она и получила прозвище", answer1: "Шапокляк", answer2: "Красная Шапочка", answer3: "Мадам Баттерфляй", want: "Шапокляк"},
	}
	return slice
}

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

func sendMess(msg tgbot.MessageConfig, text string) {
	msg.Text = text
	_, err := bot.Send(msg)
	checkErr(err)
}

func randQuest(msg tgbot.MessageConfig) {
	questObject = quests[rand.Intn(len(quests))]
	var menu = tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton(questObject.answer1),
		),
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton(questObject.answer2),
		),
		tgbot.NewKeyboardButtonRow(
			tgbot.NewKeyboardButton(questObject.answer3),
		),
	)
	msg.Text = questObject.quest
	msg.ReplyMarkup = menu
	_, err := bot.Send(msg)
	checkErr(err)
}

func main() {
	checkErr(e)
	rand.Seed(time.Now().Unix())

	u := tgbot.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := tgbot.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			switch strings.ToLower(update.Message.Command()) {
			case "start":
				sendMess(msg, "Привет")
			case "go":
				isGame = true
				randQuest(msg)
			case "stop":
				isGame = false
				msg.Text = "Игра была остановлена"
				msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
				_, err := bot.Send(msg)
				checkErr(err)
			default:
				sendMess(msg, "Я не знаю такую команду...")
			}
		} else if isGame {
			if update.Message.Text == questObject.want {
				sendMess(msg, "Верно")
				randQuest(msg)
			} else {
				sendMess(msg, "Неверно")
				randQuest(msg)
			}
		} else {
			continue
		}
	}
}
