package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

const (
	apiKey = "AIzaSyDQheLAs6b_tyO8eucAW8lAQt8yY96eXKU"
	cx     = "015184425675892934810:2ccrmmdtlku"
)

func customSearchMain(query string) string {
	client := &http.Client{Transport: &transport.APIKey{Key: apiKey}}

	svc, err := customsearch.New(client)
	if err != nil {
		log.Fatal(err)
	}

	rnd := rand.Int63n(90)

	resp, err := svc.Cse.List(query).SearchType("image").ImgSize("medium").Start(rnd).Cx(cx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// rndItem := rand.Intn(10)
	// return resp.Items[rndItem].Link

	if len(resp.Items) == 0 {
		return "Not found :("
	} else {
		i := len(resp.Items)
		rndItem := rand.Intn(i)
		return resp.Items[rndItem].Link
	}

	// for i, result := range resp.Items {
	// 	fmt.Printf("#%d: %s\n", i+1, result.Snippet)
	// 	fmt.Printf("\t%s\n", result.Link)
	// }
}

func main() {

	// https://www.googleapis.com/customsearch/v1?key=ap_key&cx=cx&q=hello&searchType=image&imgSize=xlarge&alt=json&num=10&start=1

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// https://go-telegram-bot-api.github.io/examples/commands/

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.Text = customSearchMain(update.Message.Text)
		// // msg.ReplyToMessageID = update.Message.MessageID

		// if _, err := bot.Send(msg); err != nil {
		// 	log.Panic(err)
		// }

		if update.Message.IsCommand() { // ignore any non-command Messages
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /start /status or /cat."
			case "cat":
				msg.Text = customSearchMain("cat")
				// msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			case "start":
				msg.Text = "Write text to get a related random image"
			default:
				msg.Text = "I don't know that command"
			}
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we should leave it empty.
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
