package controllers

import (
	"log"
	"net/smtp"

	"blockchain-back/modules"
	"blockchain-back/util"
)

func SendRequest(req modules.InputData) {
	subject := req.CompanyName + `の` + req.Name + `さん`
	body := `<html>
		<body>
			<h1>` + req.Hash + `書類の取得をご希望の場合、確認してください。</h1>
				<h2>名前：` + req.Name + `</h2>
				<h2>メール：` + req.Email + `</h2>
				<h2>メッセージ：` + req.Message + `</h2>
				<h2>時間：` + req.SendTime + `</h2>
			<br>
			<a href="https://localhost:5173/Ok/` + req.Name + `"> CHECK BTN </a>
			<br>
			<br>
		</body>
	</html>`

	config, auth := EmailEnv()
	err := EmailProcess(subject, body, req.Email, config, auth)
	if err != nil {
		log.Printf("Err: %v", err)
	}
}

func SendRsp(req modules.InputData) {
	subject := req.CompanyName + `の` + req.Name + `さん`
	body := `<html>
		<body>
			<h1>` + req.Hash + `メッセージを確認しました。</h1>
				<h2>名前：` + req.Name + `</h2>
				<h2>メール：` + req.Email + `</h2>
				<h2>メッセージ：` + req.Message + `</h2>
				<h2>時間：` + req.SendTime + `</h2>
			<h3>ご不明な点がございましたら、お気軽にメールでお問い合わせください。</h3>
			<br>
			<a href="http://localhost:8080/pdf/` + req.Name + `"> 履歴書のリング </a>
			<br>
			<p>※このメールは送信専用です。返信はしないでください。</p>
		</body>
	</html>`

	config, auth := EmailEnv()
	err := EmailProcess(subject, body, req.Email, config, auth)
	if err != nil {
		log.Printf("Err: %v", err)
	}
}

func EmailEnv() (util.Config, smtp.Auth) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot connect to config file: ", err)
	}

	return config, smtp.PlainAuth("", config.Username, config.Password, config.Host)
}

func EmailProcess(subject, body, Email string, config util.Config, auth smtp.Auth) error {
	msg := "Subject: " + subject + "\n" + "Content-Type: text/html; charset=utf-8\n\n" + body

	return smtp.SendMail(config.Addr, auth, config.Username, []string{Email}, []byte(msg))
}
