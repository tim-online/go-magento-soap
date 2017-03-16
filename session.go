package magento

import (
	"context"
	"encoding/xml"
	"time"
)

const (
	loginAction = "login"
)

func NewSessionService(client *Client) *SessionService {
	return &SessionService{Client: client}
}

type SessionService struct {
	Client *Client
}

func (s *SessionService) Login(requestBody *LoginRequest, ctx context.Context) (*LoginResponse, error) {
	responseBody := NewLoginResponse()
	response := NewResponse().WithData(responseBody)
	// requestBody.SessionID = s.Client.GetSession()
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

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{
		XMLName: xml.Name{
			Space: xmlns,
			Local: loginAction,
		},
	}
}

func (req *LoginRequest) WithApiUser(apiUser string) *LoginRequest {
	req.ApiUser = CDATA(apiUser)
	return req
}

func (req *LoginRequest) WithApiKey(apiKey string) *LoginRequest {
	req.ApiKey = CDATA(apiKey)
	return req
}

type LoginRequest struct {
	XMLName xml.Name `xml:"login"`

	ApiUser CDATA `xml:"username"`
	ApiKey  CDATA `xml:"apiKey"`
}

func NewLoginResponse() *LoginResponse {
	return &LoginResponse{}
}

type LoginResponse struct {
	LoginReturn string `xml:"loginReturn"`
}

func NewSession(value string) *Session {
	return &Session{
		token:  value,
		expiry: time.Now(),
	}
}

type Session struct {
	XMLName xml.Name `xml:"sessionId"`

	token  string    `xml:"-"`
	expiry time.Time `xml:-`
}

func (s *Session) Token() string {
	return s.token
}

func (s *Session) Expiry() time.Time {
	return s.expiry
}

func (s *Session) IsExpired() bool {
	now := time.Now()
	return now.After(s.Expiry())
}

func (s *Session) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(s.token, start)
}
