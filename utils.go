package magento

import (
	"encoding/xml"
	"time"
)

type TimeWithoutTimeZone struct {
	time.Time
}

func (t *TimeWithoutTimeZone) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	layout := "2006-01-02 15:04:05"
	t2, err := time.Parse(layout, value)
	if err == nil {
		*t = TimeWithoutTimeZone{Time: t2}
	}
	return err
}
