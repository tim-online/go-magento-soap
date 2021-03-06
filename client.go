package magento

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	libraryVersion = "0.0.1"
	userAgent      = "go-magento-soap/" + libraryVersion
	mediaType      = "text/xml"
	charset        = "utf-8"
	xmlns          = "https://api.nmbrs.nl/soap/v2.1/EmployeeService"
	sessionTimeout = 3600 * time.Second
)

// Client manages communication with Unit4 Multivers API
type Client struct {
	// SOAP client used to communicate with the API.
	client *http.Client

	// Url pointing to base Unit4 Multivers API
	Endpoint *url.URL

	// Credentials
	apiUser string
	apiKey  string

	// Debugging flag
	Debug bool

	// User agent for client
	UserAgent string

	// Holds current session
	session *Session

	// Optional function called after every successful request made to the DO APIs
	onRequestCompleted RequestCompletionCallback

	// Services
	CatalogProduct *CatalogProductService
	Session        *SessionService
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// NewClient returns a new Unit4 Multivers API client
func NewClient(httpClient *http.Client, baseURL *url.URL, apiUser string, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		client:    httpClient,
		Endpoint:  nil,
		UserAgent: userAgent,
		Debug:     false,
	}

	c.SetEndpoint(baseURL)
	c.SetApiUser(apiUser)
	c.SetApiKey(apiKey)

	// Services
	c.CatalogProduct = NewCatalogProductService(c)
	c.Session = NewSessionService(c)

	return c
}

func (c *Client) SetDebug(debug bool) {
	c.Debug = debug
}

func (c *Client) SetSandbox(sandbox bool) {
	if sandbox == true {
		// u, _ := url.ParseRequestURI(companies.SandboxEndpoint)
		// c.Companies.Endpoint = u
	} else {
		// u, _ := url.ParseRequestURI(companies.Endpoint)
		// c.Companies.Endpoint = u
	}
}

func (c *Client) SetEndpoint(baseURL *url.URL) {
	// set base url for use in http client
	c.Endpoint = baseURL
}

func (c *Client) NewRequest(ctx context.Context, body *Request) (*http.Request, error) {
	u := c.GetEndpoint()

	buf := new(bytes.Buffer)
	if body != nil {
		err := xml.NewEncoder(buf).Encode(body.Envelope)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		return nil, err
	}

	// optionally pass along context
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	req.Header.Add("Content-Type", fmt.Sprintf("%s; charset=%s", mediaType, charset))
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("SOAPAction", "urn:Action")

	return req, nil
}

func (c *Client) GetEndpoint() *url.URL {
	return c.Endpoint
}

// Do sends an API request and returns the API response. The API response is XML decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, responseBody *Response) (*http.Response, error) {
	if c.Debug == true {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(dump))
	}

	httpResp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, httpResp)
	}

	// close body io.Reader
	defer func() {
		if rerr := httpResp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if c.Debug == true {
		dump, _ := httputil.DumpResponse(httpResp, true)
		log.Println(string(dump))
	}

	// check if the response isn't an error
	err = CheckResponse(httpResp)
	if err != nil {
		return httpResp, err
	}

	// check the provided interface parameter
	if httpResp == nil {
		return httpResp, err
	}

	// interface implements io.Writer: write Body to it
	// if w, ok := response.Envelope.(io.Writer); ok {
	// 	_, err := io.Copy(w, httpResp.Body)
	// 	return httpResp, err
	// }

	// try to decode body into interface parameter
	err = xml.NewDecoder(httpResp.Body).Decode(responseBody.Envelope)
	if err != nil {
		errorResponse := &ErrorResponse{Response: httpResp}
		errorResponse.Message = err.Error()
		return httpResp, errorResponse
	}

	return httpResp, nil
}

func (c *Client) ApiUser() string {
	return c.apiUser
}

func (c *Client) SetApiUser(apiUser string) {
	c.apiUser = apiUser
}

func (c *Client) ApiKey() string {
	return c.apiKey
}

func (c *Client) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) GetSession() *Session {
	if c.session == nil {
		c.session = c.Login()
	}

	if c.session.IsExpired() {
		c.session = c.Login()
	}

	return c.session
}

func (c *Client) Login() *Session {
	now := time.Now()
	request := NewLoginRequest().
		WithApiUser(c.ApiUser()).
		WithApiKey(c.ApiKey())

	resp, err := c.Session.Login(request, nil)
	if err != nil {
		panic(err)
	}

	return &Session{
		token:  resp.LoginReturn,
		expiry: now.Add(sessionTimeout),
	}
}
