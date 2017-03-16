package magento

import "encoding/xml"

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(struct {
		string `xml:",cdata"`
	}{string(c)}, start)
	return nil
}
