package cablemodem

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

// CableModem is a common interface for returning cable modem stats.
type Modem interface {
	Status() (*Status, error)
	SignalData() (*SignalData, error)
}

type Status struct {
	Uptime time.Duration `json:"uptime"`
}

// DownstreamChannel is all data/stats for a downstream cabel modem channel.
type DownstreamChannel struct {
	ID                     int    `json:"id"`
	Freq                   int    `json:"freq"`  // in Hz
	Power                  int    `json:"power"` // in dBmV
	Modulation             string `json:"modulation"`
	SNR                    int    `json:"snr"` // in dB
	UnerroredCodewords     int    `json:"uncorrected_codewords"`
	CorrectableCodewords   int    `json:"correctable_codewords"`
	UncorrectableCodewords int    `json:"uncorrectable_codewords"`
}

// UpstreamChannel is all data form a upstream channel.
type UpstreamChannel struct {
	ID               int     `json:"id"`
	Freq             int     `json:"freq"`  // in Hz
	Power            int     `json:"power"` // in dBmV
	Modulation       string  `json:"modulation"`
	RangingServiceID int     `json:"ranging_service_id"`
	SymbolRate       float64 `json:"symbol_rate"` // in Msym/sec
	RangingStatus    string  `json:"ranging_status"`
}

// SignalData stores all upstream and downstream channels
type SignalData struct {
	Downstream []DownstreamChannel `json:"downstream"`
	Upstream   []UpstreamChannel   `json:"upstream"`
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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Couldn't fetch page '%s: Response code '%s'",
			u.String(), resp.Status)
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return root, nil
}
