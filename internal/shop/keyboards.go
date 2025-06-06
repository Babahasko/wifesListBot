package shop

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"
)

func InitPurchaseKB(b *bot.Bot) *reply.ReplyKeyboard{
	AddPurchaseKeyboard := reply.New(
		reply.WithPrefix("reply_keyboard"),
		reply.IsSelective(),
		reply.IsOneTimeKeyboard(),
	).
		Button("Button", b, bot.MatchTypeExact, onReplyKeyboardSelect).
		Row().
		Button("Cancel", b, bot.MatchTypeExact, onReplyKeyboardSelect)
        return AddPurchaseKeyboard
}

func onReplyKeyboardSelect(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You selected: " + string(update.Message.Text),
	})
}
