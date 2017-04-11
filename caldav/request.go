package caldav

import (
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/wothing/caldav-go/http"
	"github.com/wothing/caldav-go/icalendar"
	"github.com/wothing/caldav-go/utils"
	"github.com/wothing/caldav-go/webdav"
)

var _ = log.Print

// an CalDAV request object
type Request webdav.Request

// downcasts the request to the WebDAV interface
func (r *Request) WebDAV() *webdav.Request {
	return (*webdav.Request)(r)
}

// creates a new CalDAV request object
func NewRequest(method string, urlstr string, icaldata ...interface{}) (*Request, error) {
	if buffer, err := icalToReader(icaldata...); err != nil {
		return nil, utils.NewError(NewRequest, "unable to encode icalendar data", icaldata, err)
	} else if r, err := http.NewRequest(method, urlstr, buffer); err != nil {
		return nil, utils.NewError(NewRequest, "unable to create request", urlstr, err)
	} else {
		if buffer != nil {
			// set the content type to XML if we have a body
			r.Native().Header.Set("Content-Type", "text/calendar; charset=UTF-8")
		}
		return (*Request)(r), nil
	}
}

func icalToReader(icaldata ...interface{}) (io.Reader, error) {
	var buffer []string
	for _, icaldatum := range icaldata {
		if encoded, err := icalendar.Marshal(icaldatum); err != nil {
			return nil, utils.NewError(icalToReader, "unable to encode as icalendar data", icaldatum, err)
		} else {
			//			log.Printf("OUT: %+v", encoded)
			buffer = append(buffer, encoded)
		}
	}
	if len(buffer) > 0 {
		var encoded = strings.Join(buffer, "\n")
		return bytes.NewBuffer([]byte(encoded)), nil
	} else {
		return nil, nil
	}
}
