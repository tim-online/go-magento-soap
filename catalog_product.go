package magento

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/aodin/date"
)

const (
	catalogProductListAction                  = "catalogProductList"
	catalogProductCreateAction                = "catalogProductCreate"
	catalogProductUpdateAction                = "catalogProductUpdate"
	catalogProductInfoAction                  = "catalogProductInfo"
	ID                         IdentifierType = "ID"
	SKU                        IdentifierType = "SKU"
)

func NewCatalogProductService(client *Client) *CatalogProductService {
	return &CatalogProductService{Client: client}
}

type CatalogProductService struct {
	Client *Client
}

func (s *CatalogProductService) List(requestBody *CatalogProductListRequest, ctx context.Context) (*CatalogProductListResponse, error) {
	responseBody := NewCatalogProductListResponse()
	response := NewResponse().WithData(responseBody)
	requestBody.SessionID = s.Client.GetSession()
	request := NewRequest().WithData(requestBody)

	// create a new HTTP request
	httpReq, err := s.Client.NewRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// submit the request
	_, err = s.Client.Do(httpReq, response)
	return responseBody, err
}

func NewCatalogProductListRequest() *CatalogProductListRequest {
	return &CatalogProductListRequest{
		XMLName: xml.Name{
			Space: xmlns,
			Local: catalogProductListAction,
		},
	}
}

type CatalogProductListRequest struct {
	XMLName xml.Name `xml:"catalogProductList"`

	SessionID *Session
	Filters   *Filters `xml:filters,omitempty`
	StoreView string   `xml:"storeView,omitempty"`
}

func NewCatalogProductListResponse() *CatalogProductListResponse {
	return &CatalogProductListResponse{}
}

type CatalogProductListResponse struct {
	StoreView CatalogProductEntityArray `xml:"storeView"`
}

type CatalogProductEntityArray struct {
	// XMLName xml.Name

	Items []CatalogProductEntity `xml:"item"`
}

type CatalogProductEntity struct {
	XMLName xml.Name `xml:"item"`

	ProductID   int    `xml:"product_id"`
	Sku         string `xml:"sku"`
	Name        string `xml:"name"`
	Set         string `xml:"set"`
	Type        string `xml:"type"`
	CategoryIDs []int  `xml:"category_ids>item"`
	WebsiteIDs  []int  `xml:"website_ids>item"`
}

func (s *CatalogProductService) Create(requestBody *CatalogProductCreateRequest, ctx context.Context) (*CatalogProductCreateResponse, error) {
	responseBody := NewCatalogProductCreateResponse()
	response := NewResponse().WithData(responseBody)
	requestBody.SessionID = s.Client.GetSession()
	request := NewRequest().WithData(requestBody)

	// create a new HTTP request
	httpReq, err := s.Client.NewRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// submit the request
	_, err = s.Client.Do(httpReq, response)
	return responseBody, err
}

func NewCatalogProductCreateRequest() *CatalogProductCreateRequest {
	return &CatalogProductCreateRequest{
		XMLName: xml.Name{
			Space: xmlns,
			Local: catalogProductCreateAction,
		},
	}
}

type CatalogProductCreateRequest struct {
	XMLName xml.Name `xml:"catalogProductCreate"`

	SessionID   *Session
	Type        string                      `xml:"type"`
	Set         string                      `xml:"set"`
	Sku         string                      `xml:"sku"`
	ProductData *CatalogProductCreateEntity `xml:"productData"`
	StoreView   string                      `xml:"storeView,omitempty"`
}

func NewCatalogProductCreateResponse() *CatalogProductCreateResponse {
	return &CatalogProductCreateResponse{}
}

type CatalogProductCreateResponse struct {
	Result int `xml:"result"`
}

// The "websites" and "website_ids" or "categories" and "category_ids"
// parameters are interchangeable. In other words, you can specify an array of
// website IDs (int) and then you don't need to specify the array of website
// codes (string) and vice versa
type CatalogProductCreateEntity struct {
	Categories           []string                                   `xml:"categories"`
	Websites             []string                                   `xml:"websites"`
	Name                 string                                     `xml:"name"`
	Description          string                                     `xml:"description"`
	ShortDescription     string                                     `xml:"short_description"`
	Weight               float64                                    `xml:"weight"`
	URLKey               string                                     `xml:"url_key"`
	URLPath              string                                     `xml:"url_path"`
	Visibility           string                                     `xml:"visibility"`
	CategoryIDs          []string                                   `xml:"category_ids"`
	WebsiteIDs           []string                                   `xml:"website_ids"`
	HasOptions           bool                                       `xml:"has_options"`
	GiftMessageAvailable bool                                       `xml:"gist_message_available"`
	Price                float64                                    `xml:"price"`
	SpecialPrice         float64                                    `xml:"special_price"`
	SpecialFromDate      date.Date                                  `xml:"special_from_date"`
	SpecialToDate        date.Date                                  `xml:"special_to_date"`
	TaxClassID           int                                        `xml:"tax_class_id"`
	TierPrice            []CatalogProductTierPriceEntity            `xml:"tier_price"`
	MetaTitle            string                                     `xml:"meta_title"`
	MetaKeyword          string                                     `xml:"meta_keyword"`
	MetaDescription      string                                     `xml:"meta_description"`
	CustomDesign         string                                     `xml:"custom_design"`
	CustomLayoutUpdate   string                                     `xml:"custom_layout_update"`
	OptionsContainer     string                                     `xml:"options_container"`
	AdditionalAttributes []CatalogProductAdditionalAttributesEntity `xml:"additional_attributes"`
	StockData            []CatalogInventoryStockItemUpdateEntity    `xml:"stock_data"`
}

type CatalogProductTierPriceEntity struct {
	CustomerGroupID string  `xml:"customer_group_id"`
	Website         string  `xml:"website"`
	Qty             int     `xml:"qty"`
	Price           float64 `xml:"price"`
}

type CatalogProductAdditionalAttributesEntity struct {
	MultiData  []associativeMultiEntity `xml:"multi_data"`
	SingleData []associativeEntity      `xml:"single_date"`
}

// @TODO: implement me
type associativeMultiEntity struct {
}

// @TODO: implement me
type associativeEntity struct {
}

type CatalogInventoryStockItemUpdateEntity struct {
	Qty                     int  `xml:"qty"`
	IsInstock               bool `xml:"is_in_stock"`
	ManageStock             bool `xml:"manage_stock"`
	UseConfigManageStock    bool `xml:"use_config_manage_stock"`
	MinQty                  int  `xml:"min_qty"`
	UseConfigMinQty         bool `xml:"use_config_min_qty"`
	MinSaleQty              int  `xml:"min_sale_qty"`
	UseConfigMinSaleQty     bool `xml:"use_config_min_sale_qty"`
	MaxSaleQty              int  `xml:"max_sale_qty"`
	UseConfigMaxSaleQty     int  `xml:"use_config_max_sale_qty"`
	IsQtyDecimal            bool `xml:"is_qty_decimal"`
	Backorders              bool `xml:"backorders"`
	UseConfigBackorders     bool `xml:"use_config_backorders"`
	NotifyStockQty          bool `xml:"notify_stock_qty"`
	UseConfigNotifyStockQty bool `xml:"use_config_notify_stock_qty"`
}

func (s *CatalogProductService) Update(requestBody *CatalogProductUpdateRequest, ctx context.Context) (*CatalogProductUpdateResponse, error) {
	responseBody := NewCatalogProductUpdateResponse()
	response := NewResponse().WithData(responseBody)
	requestBody.SessionID = s.Client.GetSession()
	request := NewRequest().WithData(requestBody)

	// create a new HTTP request
	httpReq, err := s.Client.NewRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// submit the request
	_, err = s.Client.Do(httpReq, response)
	return responseBody, err
}

func NewCatalogProductUpdateRequest() *CatalogProductUpdateRequest {
	return &CatalogProductUpdateRequest{
		XMLName: xml.Name{
			Space: xmlns,
			Local: catalogProductUpdateAction,
		},
	}
}

type CatalogProductUpdateRequest struct {
	XMLName xml.Name `xml:"catalogProductUpdate"`

	SessionID      *Session
	Product        string                      `xml:"product"`
	ProductID      string                      `xml:"productId"`
	ProductData    *CatalogProductCreateEntity `xml:"productData"`
	StoreView      string                      `xml:"storeView,omitempty"`
	IdentifierType IdentifierType              `xml:"identifierType"`
}

func NewCatalogProductUpdateResponse() *CatalogProductUpdateResponse {
	return &CatalogProductUpdateResponse{}
}

type CatalogProductUpdateResponse struct {
	Result bool `xml:'result'`
}

type IdentifierType string

func (s *CatalogProductService) Info(requestBody *CatalogProductInfoRequest, ctx context.Context) (*CatalogProductInfoResponse, error) {
	responseBody := NewCatalogProductInfoResponse()
	response := NewResponse().WithData(responseBody)
	requestBody.SessionID = s.Client.GetSession()
	request := NewRequest().WithData(requestBody)

	// create a new HTTP request
	httpReq, err := s.Client.NewRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// submit the request
	_, err = s.Client.Do(httpReq, response)
	return responseBody, err
}

func NewCatalogProductInfoRequest() *CatalogProductInfoRequest {
	return &CatalogProductInfoRequest{
		XMLName: xml.Name{
			Space: xmlns,
			Local: catalogProductInfoAction,
		},
	}
}

type CatalogProductInfoRequest struct {
	XMLName xml.Name `xml:"catalogProductInfo"`

	SessionID      *Session
	Product        string                            `xml:"product"`
	ProductID      string                            `xml:"productId"`
	StoreView      string                            `xml:"storeView,omitempty"`
	Attributes     []CatalogProductRequestAttributes `xml:"attributes,omitempty"`
	IdentifierType IdentifierType                    `xml:"identifierType"`
}

func NewCatalogProductInfoResponse() *CatalogProductInfoResponse {
	return &CatalogProductInfoResponse{}
}

type CatalogProductInfoResponse struct {
	XMLName xml.Name                   `xml:"catalogProductInfoResponse"`
	Info    CatalogProductReturnEntity `xml:"info"`
}

type CatalogProductRequestAttributes struct {
	Attributes           []string `xml:"attributes"`
	AdditionalAttributes []string `xml:"additional_attributes"`
}

type CatalogProductReturnEntity struct {
	CatalogProductCreateEntity

	ProductID string              `xml:"product_id"`
	Set       int                 `xml:"set"`
	Type      string              `xml:"type"`
	Sku       string              `xml:"sku"`
	UpdatedAt TimeWithoutTimeZone `xml:"updated_at"`
	CreatedAt time.Time           `xml:"created_at"`
	TypeID    string              `xml:"type_id"`
}
