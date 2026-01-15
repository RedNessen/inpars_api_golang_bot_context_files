package telegram

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
)

// Bot представляет Telegram бота
type Bot struct {
	api       *tgbotapi.BotAPI
	chatIDs   map[int64]bool // Список активных чатов
}

// NewBot создает новый экземпляр Telegram бота
func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("Telegram bot authorized as @%s", api.Self.UserName)

	return &Bot{
		api:     api,
		chatIDs: make(map[int64]bool),
	}, nil
}

// Start запускает бота для обработки сообщений
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.handleMessage(update.Message)
	}

	return nil
}

// handleMessage обрабатывает входящие сообщения
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	log.Printf("Received message from %d: %s", chatID, text)

	// Добавляем чат в список активных
	b.chatIDs[chatID] = true

	switch {
	case text == "/start":
		b.sendStartMessage(chatID)
	case text == "/help":
		b.sendHelpMessage(chatID)
	case text == "/stop":
		delete(b.chatIDs, chatID)
		msg := tgbotapi.NewMessage(chatID, "Уведомления о новых объявлениях остановлены. Используйте /start для возобновления.")
		b.api.Send(msg)
	default:
		msg := tgbotapi.NewMessage(chatID, "Используйте /help для просмотра доступных команд.")
		b.api.Send(msg)
	}
}

// sendStartMessage отправляет приветственное сообщение
func (b *Bot) sendStartMessage(chatID int64) {
	text := `🏠 Добро пожаловать в бот мониторинга объявлений InPars!

Я буду присылать вам уведомления о новых объявлениях об аренде недвижимости.

Используйте /help для просмотра доступных команд.`

	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

// sendHelpMessage отправляет сообщение со списком команд
func (b *Bot) sendHelpMessage(chatID int64) {
	text := `📋 Доступные команды:

/start - Начать получать уведомления
/stop - Остановить уведомления
/help - Показать это сообщение

Бот автоматически мониторит новые объявления и отправляет их вам.`

	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

// SendEstate отправляет информацию об объявлении всем активным пользователям
func (b *Bot) SendEstate(estate *inpars.Estate) error {
	message := b.formatEstateMessage(estate)

	for chatID := range b.chatIDs {
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "HTML"
		msg.DisableWebPagePreview = false

		_, err := b.api.Send(msg)
		if err != nil {
			log.Printf("Failed to send message to %d: %v", chatID, err)
			// Если пользователь заблокировал бота, удаляем его из списка
			if strings.Contains(err.Error(), "blocked") || strings.Contains(err.Error(), "forbidden") {
				delete(b.chatIDs, chatID)
			}
			continue
		}
	}

	return nil
}

// formatEstateMessage форматирует информацию об объявлении для Telegram
func (b *Bot) formatEstateMessage(estate *inpars.Estate) string {
	var sb strings.Builder

	// Заголовок
	sb.WriteString(fmt.Sprintf("<b>🏠 %s</b>\n\n", estate.Title))

	// Тип объявления и продавец
	sb.WriteString(fmt.Sprintf("📌 %s", inpars.GetTypeAdName(estate.TypeAd)))
	if estate.Agent > 0 {
		sb.WriteString(fmt.Sprintf(" • %s", inpars.GetSellerTypeName(estate.Agent)))
	}
	sb.WriteString("\n\n")

	// Цена
	sb.WriteString(fmt.Sprintf("💰 <b>%s</b>", estate.FormatCost()))
	if estate.TypeAd == 1 && estate.RentTime == 2 {
		sb.WriteString(" / сутки")
	} else if estate.TypeAd == 1 {
		sb.WriteString(" / месяц")
	}
	sb.WriteString("\n\n")

	// Адрес
	if estate.Address != "" {
		sb.WriteString(fmt.Sprintf("📍 %s\n", estate.Address))
	}

	// Метро
	if estate.Metro != "" {
		sb.WriteString(fmt.Sprintf("🚇 %s\n", estate.Metro))
	}

	// Характеристики
	sb.WriteString("\n<b>Характеристики:</b>\n")

	if estate.Rooms > 0 {
		sb.WriteString(fmt.Sprintf("🛏 Комнат: %d\n", estate.Rooms))
	}

	if estate.Sq > 0 {
		sb.WriteString(fmt.Sprintf("📐 Площадь: %.1f м²\n", estate.Sq))
	}

	if estate.Floor > 0 && estate.Floors > 0 {
		sb.WriteString(fmt.Sprintf("🏢 Этаж: %d/%d\n", estate.Floor, estate.Floors))
	} else if estate.Floor > 0 {
		sb.WriteString(fmt.Sprintf("🏢 Этаж: %d\n", estate.Floor))
	}

	if estate.Material != "" {
		sb.WriteString(fmt.Sprintf("🧱 Материал: %s\n", estate.Material))
	}

	// Условия аренды
	if estate.RentTerms != nil {
		sb.WriteString("\n<b>Условия аренды:</b>\n")
		if estate.RentTerms.Deposit > 0 {
			sb.WriteString(fmt.Sprintf("💳 Залог: %s\n", formatPrice(estate.RentTerms.Deposit)))
		}
		if estate.RentTerms.Commission > 0 {
			if estate.RentTerms.CommissionType == 1 {
				sb.WriteString(fmt.Sprintf("💵 Комиссия: %d%%\n", estate.RentTerms.Commission))
			} else {
				sb.WriteString(fmt.Sprintf("💵 Комиссия: %s\n", formatPrice(estate.RentTerms.Commission)))
			}
		}
	}

	// Описание (ограниченное)
	if estate.Text != "" {
		text := estate.Text
		if len(text) > 300 {
			text = text[:297] + "..."
		}
		sb.WriteString(fmt.Sprintf("\n📝 %s\n", text))
	}

	// Контакты
	if estate.Name != "" {
		sb.WriteString(fmt.Sprintf("\n👤 Контакт: %s\n", estate.Name))
	}

	if len(estate.Phones) > 0 && estate.Phones[0] > 0 {
		sb.WriteString(fmt.Sprintf("📞 Телефон: +%d\n", estate.Phones[0]))
	}

	// Ссылка на объявление
	if estate.URL != "" {
		sb.WriteString(fmt.Sprintf("\n🔗 <a href=\"%s\">Посмотреть объявление</a>\n", estate.URL))
	}

	// Источник
	sb.WriteString(fmt.Sprintf("\n📌 Источник: %s", estate.Source))

	return sb.String()
}

// formatPrice форматирует цену с разделителями тысяч
func formatPrice(price int) string {
	if price == 0 {
		return "0 ₽"
	}

	str := ""
	n := price
	for i := 0; n > 0; i++ {
		if i > 0 && i%3 == 0 {
			str = " " + str
		}
		str = string(rune('0'+(n%10))) + str
		n /= 10
	}
	return str + " ₽"
}

// GetActiveChatIDs возвращает список активных chat ID
func (b *Bot) GetActiveChatIDs() []int64 {
	chatIDs := make([]int64, 0, len(b.chatIDs))
	for id := range b.chatIDs {
		chatIDs = append(chatIDs, id)
	}
	return chatIDs
}

// HasActiveChats проверяет, есть ли активные чаты
func (b *Bot) HasActiveChats() bool {
	return len(b.chatIDs) > 0
}
