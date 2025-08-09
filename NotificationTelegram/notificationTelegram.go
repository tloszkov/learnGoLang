package sendtelegramnotification

import (
	loadenv "learnGoLang/LoadEnv"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegramNotification(message string) {
	// Bot token - replace with your actual bot token
	botToken := loadenv.TELEGRAM_BOT_TOKEN

	// Chat ID - you need to get this by messaging your bot first
	// To get your chat ID, message your bot and visit:
	// https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates
	var chatID int64 = loadenv.TELEGRAM_CHAT_ID // Your actual numeric chat ID

	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	log.Printf("Bot created successfully: @%s", bot.Self.UserName)

	// Send "Hello world!" to the specified chat ID automatically
	msg := tgbotapi.NewMessage(chatID, message)

	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Printf("Successfully sent '%s' to chat ID %d (Message ID: %d)", msg.Text, chatID, sentMsg.MessageID)
	}
}
