package models

type Region struct {
	Counties        []County
	Constituencies  []Constituency
	Wards           []Ward
	PollingStations []PollingStation
}
