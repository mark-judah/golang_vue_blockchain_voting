package models

type Voter struct {
	VoterId          string `json:"voterId"`
	PhoneNumber      string `json:"phoneNumber"`
	Timestamp        string `json:"timestamp"`
	VoterDetailsHash string `json:"voterDetailsHash"`
}
