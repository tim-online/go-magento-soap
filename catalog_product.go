package magento

import (
	"context"
	"encoding/xml"
	"log"
)

const (
	catalogProductListAction = "catalogProductList"
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
	requestBody.SessionID = s.Client.GetSessionID()
	request := NewRequest().WithData(requestBody)

	// create a new HTTP request
	httpReq, err := s.Client.NewRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	// submit the request
	_, err = s.Client.Do(httpReq, response)
	log.Printf("%+v", responseBody)
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

	SessionID *SessionID
	Filters   *Filters `xml:filters,omitempty`
	StoreView string   `xml:"storeView,omitempty"`
}

func NewCatalogProductListResponse() *CatalogProductListResponse {
	return &CatalogProductListResponse{}
}

type CatalogProductListResponse struct {
	StoreView CatalogProductEntityArray `xml:"storeView"`
}

func (r *CatalogProductListResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Fatal("aaaaaaaaaah")
	return nil
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
	CategoryIds []int  `xml:"category_ids>item"`
	WebsiteIds  []int  `xml:"website_ids>item"`
}
