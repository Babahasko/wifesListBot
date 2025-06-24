package bot

import (
	"fmt"
	"shopping_bot/internal/repository"
)

const (
	MsgStart = "✌️Привет! Я — твой помощник для управления списками покупок. 📋\n" +
		"\n" +
		"Введи название списка и просто накидывай товары\n" +
		"\n" +
		"Так ты всегда будешь помнить, что купить 🛒\n" +
		"\n" +
		"🔸 Чтобы начать, просто напиши /add или выбери нужное действие из меню."

	MsgWriteListName = "Введите название списка📋"
	MsgNoLists = "У вас ещё нет ни одного списка! Создайте для начала командой /add_list"
	MsgYourLists = "Ваши списки"
	MsgVoznya = "Пошла возня"

	MsgWriteItemName = "Отправляйте названия покупок.\n" +
		"\n" +
		"Чтобы закончить формировать список вы можете:\n" +
		"1.Нажать кнопку `✅Завершить список`\n" +
		"2.Ввести команду /end\n"
	MsgPitty = "Жаль, что вы прервались!"
)

func formListMessage(listName string, items []*repository.ShoppingItem) string {
	itemsText := fmt.Sprintf("🛒 %s\n", listName)
	for i, item := range items {
		itemsText += fmt.Sprintf("%d. %s\n", i+1, item.Name)
	}
	return itemsText
}
