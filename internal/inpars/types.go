package inpars

import "time"

// APIError представляет ошибку API
type APIError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
}

// Meta содержит метаданные ответа
type Meta struct {
	Limit           int `json:"limit,omitempty"`
	TotalCount      int `json:"totalCount,omitempty"`
	UpdateLimit     int `json:"updateLimit,omitempty"`
	UpdateRemaining int `json:"updateRemaining,omitempty"`
	RateLimit       int `json:"rateLimit,omitempty"`
	RateRemaining   int `json:"rateRemaining,omitempty"`
	RateReset       int `json:"rateReset,omitempty"`
}

// Estate представляет объявление о недвижимости
type Estate struct {
	ID         int      `json:"id"`
	RegionID   int      `json:"regionId"`
	CityID     int      `json:"cityId"`
	MetroID    int      `json:"metroId,omitempty"`
	TypeAd     int      `json:"typeAd"`     // 1-сдам, 2-продам, 3-сниму, 4-куплю
	SectionID  int      `json:"sectionId"` // ID раздела недвижимости
	CategoryID int      `json:"categoryId"`
	Title      string   `json:"title"`
	Address    string   `json:"address"`
	Floor      int      `json:"floor,omitempty"`
	Floors     int      `json:"floors,omitempty"`
	Sq         float64  `json:"sq,omitempty"`         // Площадь
	SqLand     float64  `json:"sqLand,omitempty"`     // Площадь участка
	SqLiving   float64  `json:"sqLiving,omitempty"`   // Жилая площадь
	SqKitchen  float64  `json:"sqKitchen,omitempty"`  // Площадь кухни
	Cost       int      `json:"cost"`                 // Стоимость
	Text       string   `json:"text"`                 // Описание
	Images     []string `json:"images"`               // Ссылки на фото
	Lat        float64  `json:"lat"`                  // Широта
	Lng        float64  `json:"lng"`                  // Долгота
	Name       string   `json:"name"`                 // Имя продавца
	Phones     []int64  `json:"phones"`               // Телефоны
	URL        string   `json:"url"`                  // Ссылка на источник
	Agent      int      `json:"agent"`                // 0-собственник, 1-агент, 2-застройщик
	Source     string   `json:"source"`               // Название источника
	SourceID   int      `json:"sourceId"`             // ID источника
	Created    string   `json:"created"`              // Дата создания
	Updated    string   `json:"updated"`              // Дата обновления

	// Дополнительные поля (требуют expand параметра)
	Region         string      `json:"region,omitempty"`
	City           string      `json:"city,omitempty"`
	Type           string      `json:"type,omitempty"`
	Section        string      `json:"section,omitempty"`
	Category       string      `json:"category,omitempty"`
	Metro          string      `json:"metro,omitempty"`
	Material       string      `json:"material,omitempty"`
	RentTime       int         `json:"rentTime,omitempty"`       // 0-не указан, 1-длительно, 2-посуточно
	IsNew          bool        `json:"isNew,omitempty"`          // Новостройка
	Rooms          int         `json:"rooms,omitempty"`          // Количество комнат
	PhoneProtected bool        `json:"phoneProtected,omitempty"` // Подменный номер
	ParseID        string      `json:"parseId,omitempty"`        // ID на источнике
	IsApartments   bool        `json:"isApartments,omitempty"`   // Апартаменты
	RentTerms      *RentTerms  `json:"rentTerms,omitempty"`      // Условия аренды
	House          *House      `json:"house,omitempty"`          // Информация о доме
}

// RentTerms условия аренды
type RentTerms struct {
	Commission       int `json:"commission,omitempty"`       // Комиссия
	CommissionType   int `json:"commissionType,omitempty"`   // 1-процент, 2-фикс.сумма
	Deposit          int `json:"deposit,omitempty"`          // Залог
	Utilities        int `json:"utilities,omitempty"`        // 1-арендатор, 2-включено
	UtilitiesMeters  int `json:"utilitiesMeters,omitempty"`  // Счетчики
	UtilitiesPrice   int `json:"utilitiesPrice,omitempty"`   // Стоимость ЖКУ
}

// House информация о доме
type House struct {
	BuildYear       int `json:"buildYear,omitempty"`
	CargoLifts      int `json:"cargoLifts,omitempty"`
	PassengerLifts  int `json:"passengerLifts,omitempty"`
}

// EstateListResponse ответ на запрос списка объявлений
type EstateListResponse struct {
	Data []Estate `json:"data"`
	Meta Meta     `json:"meta"`
}

// EstateResponse ответ на запрос одного объявления
type EstateResponse struct {
	Data Estate `json:"data"`
	Meta Meta   `json:"meta"`
}

// Region представляет регион
type Region struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// RegionListResponse ответ на запрос списка регионов
type RegionListResponse struct {
	Data []Region `json:"data"`
	Meta Meta     `json:"meta"`
}

// City представляет город
type City struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	RegionID int    `json:"regionId"`
}

// CityListResponse ответ на запрос списка городов
type CityListResponse struct {
	Data []City `json:"data"`
	Meta Meta   `json:"meta"`
}

// GetTypeAdName возвращает текстовое название типа объявления
func GetTypeAdName(typeAd int) string {
	switch typeAd {
	case 1:
		return "Сдам"
	case 2:
		return "Продам"
	case 3:
		return "Сниму"
	case 4:
		return "Куплю"
	default:
		return "Неизвестно"
	}
}

// GetSellerTypeName возвращает текстовое название типа продавца
func GetSellerTypeName(agent int) string {
	switch agent {
	case 0:
		return "Собственник"
	case 1:
		return "Агент"
	case 2:
		return "Застройщик"
	default:
		return "Неизвестно"
	}
}

// GetCreatedTime возвращает время создания объявления
func (e *Estate) GetCreatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, e.Created)
}

// GetUpdatedTime возвращает время обновления объявления
func (e *Estate) GetUpdatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, e.Updated)
}

// FormatCost форматирует стоимость в читаемый вид
func (e *Estate) FormatCost() string {
	if e.Cost == 0 {
		return "Не указана"
	}
	return formatPrice(e.Cost)
}

// formatPrice форматирует цену с разделителями тысяч
func formatPrice(price int) string {
	if price == 0 {
		return "0"
	}

	str := ""
	for i, digit := range reverseString(intToString(price)) {
		if i > 0 && i%3 == 0 {
			str = " " + str
		}
		str = string(digit) + str
	}
	return str + " ₽"
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+(n%10))) + result
		n /= 10
	}
	return result
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
