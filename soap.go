package magento

import (
	"encoding/xml"
	"net/url"
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
		Header: NewHeader(),
		Body:   NewBody(),
	}
}

// http://stackoverflow.com/questions/16202170/marshalling-xml-go-xmlname-xmlns
type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Header *Header `xml:"Header"`
	Body   *Body   `xml:"Body"`
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
