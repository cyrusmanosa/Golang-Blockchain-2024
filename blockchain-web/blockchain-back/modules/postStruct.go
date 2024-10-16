package modules

type InputData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Message     string `json:"message"`
	// Cv          string `json:"cv"`
	File        []byte `json:"cv"`
	Status      string `json:"status"`
	SendTime    string `json:"send_time"`
	ConfirmTime string `json:"confirm_time"`
}
