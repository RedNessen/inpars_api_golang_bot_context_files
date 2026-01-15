package inpars

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	BaseURL     = "https://inpars.ru/api/v2"
	DefaultLimit = 50
)

// Client представляет клиент для работы с API InPars
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// NewClient создает новый клиент API
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token:   token,
		baseURL: BaseURL,
	}
}

// getAuthHeader возвращает заголовок авторизации для Basic Auth
func (c *Client) getAuthHeader() string {
	// Формат: base64(token:)
	auth := c.token + ":"
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encoded
}

// doRequest выполняет HTTP запрос с авторизацией
func (c *Client) doRequest(method, endpoint string, params url.Values) ([]byte, error) {
	reqURL := c.baseURL + endpoint
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.getAuthHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Проверка на ошибки rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return nil, fmt.Errorf("rate limit exceeded: %s", apiErr.Message)
		}
		return nil, fmt.Errorf("rate limit exceeded (429)")
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil {
			return nil, fmt.Errorf("API error (%d): %s", apiErr.Status, apiErr.Message)
		}
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	return body, nil
}

// GetEstateList получает список объявлений
func (c *Client) GetEstateList(params *EstateListParams) (*EstateListResponse, error) {
	urlParams := params.ToURLValues()

	body, err := c.doRequest("GET", "/estate", urlParams)
	if err != nil {
		return nil, err
	}

	var response EstateListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetEstate получает информацию об одном объявлении
func (c *Client) GetEstate(id int) (*EstateResponse, error) {
	endpoint := fmt.Sprintf("/estate/%d", id)

	body, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response EstateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetRegions получает список регионов
func (c *Client) GetRegions() (*RegionListResponse, error) {
	body, err := c.doRequest("GET", "/region", nil)
	if err != nil {
		return nil, err
	}

	var response RegionListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// EstateListParams параметры для получения списка объявлений
type EstateListParams struct {
	SortBy      string   // updated_desc, updated_asc, created_desc, created_asc, id_desc, id_asc
	LastID      int      // ID последнего объявления для пагинации
	TimeStart   int64    // UNIX timestamp начала выборки
	TimeEnd     int64    // UNIX timestamp конца выборки
	RegionID    []int    // ID регионов
	CityID      []int    // ID городов
	MetroID     []int    // ID станций метро
	TypeAd      []int    // Тип: 1-сдам, 2-продам, 3-сниму, 4-куплю
	SectionID   []int    // ID разделов
	CategoryID  []int    // ID категорий
	SellerType  []int    // Продавец: 1-собственник, 2-агент, 3-застройщик
	WithPhoto   *int     // 0-без фото, 1-с фото
	IsNew       *int     // 0-вторичка, 1-новостройка
	CostMin     int      // Цена от
	CostMax     int      // Цена до
	FloorMin    int      // Этаж от
	FloorMax    int      // Этаж до
	SqMin       float64  // Площадь от
	SqMax       float64  // Площадь до
	SourceID    []int    // ID источников (1-avito, 2-cian, 5-youla и т.д.)
	Fields      []string // Поля для возврата
	Expand      []string // Дополнительные поля
	Limit       int      // Лимит объектов (по умолчанию 500, макс 1000)
}

// ToURLValues преобразует параметры в url.Values
func (p *EstateListParams) ToURLValues() url.Values {
	values := url.Values{}

	if p.SortBy != "" {
		values.Set("sortBy", p.SortBy)
	}
	if p.LastID > 0 {
		values.Set("lastId", strconv.Itoa(p.LastID))
	}
	if p.TimeStart > 0 {
		values.Set("timeStart", strconv.FormatInt(p.TimeStart, 10))
	}
	if p.TimeEnd > 0 {
		values.Set("timeEnd", strconv.FormatInt(p.TimeEnd, 10))
	}
	if len(p.RegionID) > 0 {
		values.Set("regionId", intsToString(p.RegionID))
	}
	if len(p.CityID) > 0 {
		values.Set("cityId", intsToString(p.CityID))
	}
	if len(p.MetroID) > 0 {
		values.Set("metroId", intsToString(p.MetroID))
	}
	if len(p.TypeAd) > 0 {
		values.Set("typeAd", intsToString(p.TypeAd))
	}
	if len(p.SectionID) > 0 {
		values.Set("sectionId", intsToString(p.SectionID))
	}
	if len(p.CategoryID) > 0 {
		values.Set("categoryId", intsToString(p.CategoryID))
	}
	if len(p.SellerType) > 0 {
		values.Set("sellerType", intsToString(p.SellerType))
	}
	if p.WithPhoto != nil {
		values.Set("withPhoto", strconv.Itoa(*p.WithPhoto))
	}
	if p.IsNew != nil {
		values.Set("isNew", strconv.Itoa(*p.IsNew))
	}
	if p.CostMin > 0 {
		values.Set("costMin", strconv.Itoa(p.CostMin))
	}
	if p.CostMax > 0 {
		values.Set("costMax", strconv.Itoa(p.CostMax))
	}
	if p.FloorMin > 0 {
		values.Set("floorMin", strconv.Itoa(p.FloorMin))
	}
	if p.FloorMax > 0 {
		values.Set("floorMax", strconv.Itoa(p.FloorMax))
	}
	if p.SqMin > 0 {
		values.Set("sqMin", fmt.Sprintf("%.1f", p.SqMin))
	}
	if p.SqMax > 0 {
		values.Set("sqMax", fmt.Sprintf("%.1f", p.SqMax))
	}
	if len(p.SourceID) > 0 {
		values.Set("sourceId", intsToString(p.SourceID))
	}
	if len(p.Fields) > 0 {
		values.Set("fields", stringSliceToString(p.Fields))
	}
	if len(p.Expand) > 0 {
		values.Set("expand", stringSliceToString(p.Expand))
	}
	if p.Limit > 0 {
		values.Set("limit", strconv.Itoa(p.Limit))
	}

	return values
}

// Вспомогательные функции для преобразования слайсов в строки
func intsToString(ints []int) string {
	if len(ints) == 0 {
		return ""
	}
	result := strconv.Itoa(ints[0])
	for i := 1; i < len(ints); i++ {
		result += "," + strconv.Itoa(ints[i])
	}
	return result
}

func stringSliceToString(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	result := strings[0]
	for i := 1; i < len(strings); i++ {
		result += "," + strings[i]
	}
	return result
}
