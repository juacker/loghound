package clf

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Clf represents a common log format entry
type Clf struct {
	RemoteHost    string    `json:"remote_host"`
	RemoteLogname string    `json:"remote_logname"`
	AuthUser      string    `json:"auth_user"`
	Date          time.Time `json:"date"`
	Request       *Request  `json:"request"`
	Status        int       `json:"status"`
	Bytes         int       `json:"bytes"`
}

// Request represents the request field in a Clf entry
type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

var clfParser = regexp.MustCompile(`^(?P<remotehost>\S+) (?P<remotelogname>\S+) (?P<authuser>\S+) \[(?P<date>[^\]]+)\] "(?P<method>[A-Z]+) (?P<path>[^ "]+)? HTTP/[0-9.]+" (?P<status>[0-9]{3}) (?P<bytes>[0-9]+|-)`)

// Parse a Clf entry
func Parse(s string) (*Clf, error) {
	match := clfParser.FindStringSubmatch(s)

	date, err := time.Parse("02/Jan/2006:15:04:05 -0700", match[4])
	if err != nil {
		return nil, fmt.Errorf("fail parsing date in common log format: %s", match[4])
	}

	status, err := strconv.Atoi(match[7])
	if err != nil {
		return nil, fmt.Errorf("fail parsing status in common log format: %s", match[7])
	}

	bytes, err := strconv.Atoi(match[8])
	if err != nil {
		return nil, fmt.Errorf("fail parsing bytes in common log format: %s", match[8])
	}

	return &Clf{
		RemoteHost:    match[1],
		RemoteLogname: match[2],
		AuthUser:      match[3],
		Date:          date,
		Request: &Request{
			Method: match[5],
			Path:   match[6],
		},
		Status: status,
		Bytes:  bytes,
	}, nil
}
