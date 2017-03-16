package magento

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a XML response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	// @TODO: figure out nmbrs errors
	errorResponse := &ErrorResponse{Response: r}

	err := checkContentType(r)
	if err != nil {
		errorResponse.Message = err.Error()
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errorResponse
	}

	if len(data) == 0 {
		return errorResponse
	}

	// convert xml to struct
	err = xml.Unmarshal(data, errorResponse)
	if err != nil {
		errorResponse.Message = fmt.Sprintf("Malformed json response")
		return errorResponse
	}

	return errorResponse
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
	if contentType != "application/soap+xml" {
		return fmt.Errorf("Expected Content-Type \"application/soap+xml\", got \"%s\"", contentType)
	}

	return nil
}
