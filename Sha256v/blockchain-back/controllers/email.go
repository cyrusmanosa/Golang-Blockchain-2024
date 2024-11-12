package controllers

import (
	"Sha256v/modules"
	"Sha256v/util"
	"log"
	"net/smtp"
)

func SendRequest(req modules.InputData) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot connect to config file: ", err)
	}
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	subject := req.CompanyName + `の` + req.Name + `さん`
	body := `<html>
		<body>
			<h3>履歴書の取得をご希望の場合、確認してください。</h3>
				<h2>名前：` + req.Name + `</h2>
				<h2>会社名前：` + req.CompanyName + `</h2>
				<h2>メール：` + req.Email + `</h2>
				<h2>メッセージ：` + req.Message + `</h2>
				<h2>時間：` + req.SendTime + `</h2>
			<br>
			<a href="https://localhost:5173/Ok/` + req.Name + `"> CHECK BTN </a>
			<br>
			<br>
		</body>
	</html>`

	msg := "Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		body

	err = smtp.SendMail(config.Addr, auth, config.Username, []string{req.Email}, []byte(msg))
	if err != nil {
		log.Printf("Err: %v", err)
	}
}

func SendRsp(req modules.InputData) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot connect to config file: ", err)
	}
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	subject := req.CompanyName + `の` + req.Name + `さん`
	body := `<html>
		<body>
			<h3>メッセージを確認しました。ぜひ、機会をいただけますと幸いです。</h3>
				<h2>名前：` + req.Name + `</h2>
				<h2>会社名前：` + req.CompanyName + `</h2>
				<h2>メール：` + req.Email + `</h2>
				<h2>メッセージ：` + req.Message + `</h2>
				<h2>時間：` + req.SendTime + `</h2>
			<h3>ご不明な点がございましたら、お気軽にメールでお問い合わせください。</h3>
			<br>
			<a href="http://localhost:8080/svg/` + req.Name + `"> 履歴書のリング </a>
			<br>
			<p>※このメールは送信専用です。返信はしないでください。</p>
		</body>
	</html>`

	msg := "Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		body

	err = smtp.SendMail(config.Addr, auth, config.Username, []string{req.Email}, []byte(msg))
	if err != nil {
		log.Printf("Err: %v", err)
	}
}
