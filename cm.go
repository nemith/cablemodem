package cablemodem

import (
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// CableModem is a common interface for returning cable modem stats.
type Modem interface {
	SignalData() *SignalData
}

// DownstreamChannel is all data/stats for a downstream cabel modem channel.
type DownstreamChannel struct {
	ID                     int
	Freq                   int // in Hz
	Power                  int // in dBmV
	Modulation             string
	SNR                    int // in dB
	UnerroredCodewords     int
	CorrectableCodewords   int
	UncorrectableCodewords int
}

// UpstreamChannel is all data form a upstream channel.
type UpstreamChannel struct {
	ID               int
	Freq             int // in Hz
	Power            int // in dBmV
	Modulation       string
	RangingServiceID int
	SymbolRate       float64 // in Msym/sec
	RangingStatus    string
}

// SignalData stores all upstream and downstream channels
type SignalData struct {
	Downstream []DownstreamChannel
	Upstream   []UpstreamChannel
}

type LogPriority int

const (
	Emergency LogPriority = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Informational
	Debugging
)

type Log struct {
	Time       time.Time
	Priority   LogPriority
	Code       string
	Message    string
	Attributes []string
}

// fetchDoc will fetch a webpage and return the root node for it.
func fetchDoc(u url.URL) (*html.Node, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return root, nil
}
