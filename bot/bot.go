package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func SendMessage(msg tgbotapi.Chattable) {
    if _, err := Bot.Send(msg); err != nil {
        log.Printf("Error Sending to bot: %s", err)
    }
}

func EditMessageAndKeyboard(chatId int64, msgId int, text string, replyMarkup tgbotapi.InlineKeyboardMarkup) {
    editMsg := tgbotapi.NewEditMessageText(
        chatId,
        msgId,
        text,
    )
    SendMessage(editMsg)
    if len(replyMarkup.InlineKeyboard) > 0 {
        editKeyboard := tgbotapi.NewEditMessageReplyMarkup(
            chatId,
            msgId,
            replyMarkup,
        )
        SendMessage(editKeyboard)
    }
}

func EditMessage(chatId int64, msgId int, text string) {
    editMsg := tgbotapi.NewEditMessageText(
        chatId,
        msgId,
        text,
    )
    SendMessage(editMsg)
}
