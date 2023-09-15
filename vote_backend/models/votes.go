package models

type Vote struct {
	Txid        string `json:"txid"`
	NodeId      string `json:"nodeId"`
	CandidateId string `json:"candidateId"`
	Timestamp   string `json:"timestamp"`
	VoteHash    string `json:"voteHash"`
}
