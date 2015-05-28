package cablemodem

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func findTable(n *html.Node, tableName string) (*html.Node, error) {
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Th && scrape.Text(n) == tableName {
			return true
		}
		return false
	}

	thNode, ok := scrape.Find(n, matcher)
	if !ok {
		return nil, fmt.Errorf("Cannot find th")
	}

	tableNode, ok := scrape.FindParent(thNode, func(n *html.Node) bool {
		if n.DataAtom == atom.Table {
			return true
		}
		return false
	})

	if !ok {
		return nil, fmt.Errorf("Cannot find parent")
	}

	return tableNode, nil
}

func parseTable(n *html.Node) (map[string][]string, int) {
	rowMatcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Tr && n.FirstChild != nil {
			return n.FirstChild.DataAtom == atom.Td
		}
		return false
	}

	resMap := make(map[string][]string)
	cols := 0

	rows := scrape.FindAll(n, rowMatcher)
	for _, row := range rows {
		cells := scrape.FindAll(row, scrape.ByTag(atom.Td))
		if cols <= 0 {
			cols = len(cells) - 1
		}
		vals := make([]string, cols)
		for i, cell := range cells[1:] {
			vals[i] = scrape.Text(cell)
		}
		resMap[scrape.Text(cells[0])] = vals
	}
	return resMap, cols
}

type SurfboardCM struct {
	baseURL url.URL
}

func NewSurfboardCM(hostname string) *SurfboardCM {
	return &SurfboardCM{
		baseURL: url.URL{
			Scheme: "http",
			Host:   hostname,
		},
	}
}

func mustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func mustParseFloat64(s string) float64 {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func (cm *SurfboardCM) SignalData() *SignalData {
	u := cm.baseURL
	u.Path = "cmSignalData.htm"
	doc, err := fetchDoc(u)

	usTable, err := findTable(doc, "Upstream")
	if err != nil {
		panic(err)
	}
	data, cols := parseTable(usTable)
	upChannels := make([]UpstreamChannel, cols)
	for key, stats := range data {
		for i, val := range stats {
			switch key {
			case "Channel ID":
				upChannels[i].ID = mustAtoi(val)
			case "Frequency":
				val = strings.Split(val, " ")[0]
				upChannels[i].Freq = mustAtoi(val)
			case "Ranging Service ID":
				upChannels[i].RangingServiceID = mustAtoi(val)
			case "Symbol Rate":
				val = strings.Split(val, " ")[0]
				upChannels[i].SymbolRate = mustParseFloat64(val)
			case "Power Level":
				val = strings.Split(val, " ")[0]
				upChannels[i].Power = mustAtoi(val)
			case "Upstream Modulation":
				upChannels[i].Modulation = val
			case "Ranging Status":
				upChannels[i].RangingStatus = val
			}
		}
	}

	dsTable, err := findTable(doc, "Downstream")
	if err != nil {
		panic(err)
	}
	data, cols = parseTable(dsTable)
	dsChannels := make([]DownstreamChannel, cols)
	for key, stats := range data {
		for i, val := range stats {
			switch {
			case key == "Channel ID":
				dsChannels[i].ID = mustAtoi(val)
			case key == "Frequency":
				val = strings.Split(val, " ")[0]
				dsChannels[i].Freq = mustAtoi(val)
			case key == "Signal to Noise Ratio":
				val = strings.Split(val, " ")[0]
				dsChannels[i].SNR = mustAtoi(val)
			case key == "Downstream Modulation":
				dsChannels[i].Modulation = val
			case strings.Contains(key, "Power Level"):
				val = strings.Split(val, " ")[0]
				dsChannels[i].Power = mustAtoi(val)

			}
		}
	}

	statsTable, err := findTable(doc, "Signal Stats (Codewords)")
	if err != nil {
		panic(err)
	}
	data, cols = parseTable(statsTable)
	for key, stats := range data {
		for i, val := range stats {
			switch {
			case key == "Total Unerrored Codewords":
				dsChannels[i].UnerroredCodewords = mustAtoi(val)
			case key == "Total Correctable Codewords":
				dsChannels[i].CorrectableCodewords = mustAtoi(val)
			case key == "Total Uncorrectable Codewords":
				dsChannels[i].UncorrectableCodewords = mustAtoi(val)
			}
		}
	}

	return &SignalData{
		Upstream:   upChannels,
		Downstream: dsChannels,
	}
}
