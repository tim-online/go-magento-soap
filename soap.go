package magento

import (
	"encoding/xml"
	"net/url"
	"time"
)

func NewRequest() *Request {
	return &Request{
		Envelope: NewEnvelope(),
	}
}

// <SOAP-ENV:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:Magento" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/">
//    <SOAP-ENV:Header/>
//    <SOAP-ENV:Body>
//       <urn:catalogProductList SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
//          <sessionId xsi:type="xsd:string">?</sessionId>
//          <filters xsi:type="urn:filters">
//             <!--You may enter the following 2 items in any order-->
//             <!--Optional:-->
//             <filter xsi:type="urn:associativeArray" soapenc:arrayType="urn:associativeEntity[]"/>
//             <!--Optional:-->
//             <complex_filter xsi:type="urn:complexFilterArray" soapenc:arrayType="urn:complexFilter[]"/>
//          </filters>
//          <storeView xsi:type="xsd:string">?</storeView>
//       </urn:catalogProductList>
//    </SOAP-ENV:Body>
// </SOAP-ENV:Envelope>

type Request struct {
	Envelope *Envelope `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Action   *url.URL  `xml:""`
}

func (r *Request) WithData(data interface{}) *Request {
	r.Envelope.Body.Data = data
	return r
}

func NewResponse() *Response {
	return &Response{
		Envelope: NewEnvelope(),
	}
}

type Response struct {
	Envelope *Envelope `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
}

func (r *Response) WithData(data interface{}) *Response {
	r.Envelope.Body.Data = data
	return r
}

func NewEnvelope() *Envelope {
	return &Envelope{
		Xsi: "http://www.w3.org/2001/XMLSchema-instance",
		Xsd: "http://www.w3.org/2001/XMLSchema",
		// SoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
		Urn:     "urn:Magento",
		Soapenc: "http://schemas.xmlsoap.org/soap/encoding/",

		Header: NewHeader(),
		Body:   NewBody(),
	}
}

// http://stackoverflow.com/questions/16202170/marshalling-xml-go-xmlname-xmlns
type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	// SoapEnv string   `xml:"xmlns:SOAP-ENV,attr"`
	Urn     string `xml:"xmlns:urn,attr"`
	Soapenc string `xml:"xmlns:soapenc,attr"`

	Header *Header `xml:"Header"`
	Body   *Body   `xml:"Body"`
	Test   string  `xml:"test"`
}

func NewHeader() *Header {
	return &Header{
		Data: nil,
	}
}

type Header struct {
	Data interface{}
}

type Body struct {
	// If the XML element contains a sub-element that hasn't matched any
	// of the above rules and the struct has a field with tag ",any",
	// unmarshal maps the sub-element to that struct field.
	Data interface{} `xml:",any"`
}

func NewBody() *Body {
	return &Body{}
}

func NewSessionID(value string) *SessionID {
	return &SessionID{
		token:  value,
		Expiry: time.Now(),
	}
}

type SessionID struct {
	XMLName xml.Name `xml:"sessionId"`

	token  string    `xml:"-"`
	Expiry time.Time `xml:-`
}

func (s *SessionID) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(s.token, start)
}
