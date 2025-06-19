package shop

import "fmt"

const (
	MsgStart = "✌️Привет! Я — твой помощник для управления списками покупок. 📋\n" +
		"\n" +
		"Введи название списка и просто накидывай товары\n" +
		"\n" +
		"Так ты всегда будешь помнить, что купить 🛒\n" +
		"\n" +
		"🔸 Чтобы начать, просто напиши /add или выбери нужное действие из меню."

	MsgVoznya = "Пошла возня"
)

func formListMessage(items []string) string {
	itemsText := "🛒 *Ваш список покупок:*\n"
	for i, item := range items {
		itemsText += fmt.Sprintf("%d. %s\n", i+1, item)
	}
        return itemsText
}
