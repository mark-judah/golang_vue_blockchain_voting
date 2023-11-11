package models

type Transaction struct {
	Txid            string `json:"txid"`
	NodeId          string `json:"nodeId"`
	CandidateId     string `json:"candidateId"`
	Timestamp       string `json:"timestamp"`
	TransactionHash string `json:"transactionHash"`
}
