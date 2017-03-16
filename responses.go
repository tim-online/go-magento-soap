package magento

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. API error responses are expected to have either no response
// body, or a XML response body that maps to ErrorResponse. Any other response
// body will be silently ignored.
func CheckResponse(r *http.Response) error {
	errorResponse := &ErrorResponse{Response: r}

	err := checkContentType(r)
	if err != nil {
		errorResponse.Message = err.Error()
	}

	// if c := r.StatusCode; c >= 200 && c <= 299 {
	// 	return nil
	// }

	// read data and copy it back
	data, err := ioutil.ReadAll(r.Body)
	r.Body = nopCloser{bytes.NewReader(data)}
	if err != nil {
		return errorResponse
	}

	if len(data) == 0 {
		return errorResponse
	}

	// convert xml to struct
	err = xml.Unmarshal(data, errorResponse)
	if err != nil {
		errorResponse.Message = fmt.Sprintf("Malformed xml response")
		return errorResponse
	}

	if errorResponse.Message != "" {
		log.Printf("%+v", errorResponse)
		log.Println(errorResponse.Message)
		log.Println(errorResponse.Reason)
		log.Println(errorResponse.Code)
		return errorResponse
	}

	return nil
}

// An ErrorResponse reports the error caused by an API request
// <?xml version="1.0" encoding="UTF-8"?>
// <SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
//   <SOAP-ENV:Body>
//     <SOAP-ENV:Fault>
//       <faultcode>Sender</faultcode>
//       <faultstring>Invalid XML</faultstring>
//     </SOAP-ENV:Fault>
//   </SOAP-ENV:Body>
// </SOAP-ENV:Envelope>type ErrorResponse struct {
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Fault code
	Code string `xml:"Body>Fault>faultcode"`

	// Fault message
	Message string `xml:"Body>Fault>faultstring"`

	// Reason
	Reason string `xml:"Body>Fault>Reason>Text"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d (%v %v)",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Reason)
}

func checkContentType(response *http.Response) error {
	// check content-type (application/soap+xml; charset=utf-8)
	header := response.Header.Get("Content-Type")
	contentType := strings.Split(header, ";")[0]
	if contentType != "text/xml" {
		return fmt.Errorf("Expected Content-Type \"text/xml\", got \"%s\"", contentType)
	}

	return nil
}
