package models

type MetaInfo struct {
	Status 		string 		`json:"status"`
	Message 	string 		`json:"message"`
	Data 		interface{} `json:"data,omitempty"`
	Errors 		interface{} `json:"errors,omitempty"`
	Meta 		interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page 		int `json:"page"`
	Limit 		int `json:"limit"`
	TotalData 	int `json:"total_data"`
	TotalPages 	int `json:"total_pages"`
}