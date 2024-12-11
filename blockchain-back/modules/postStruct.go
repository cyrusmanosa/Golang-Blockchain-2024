package modules

type InputData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Message     string `json:"message"`
	Hash        string `json:"hash"`
	File        []byte `json:"file"`
	Status      string `json:"status"`
	SendTime    string `json:"send_time"`
}
