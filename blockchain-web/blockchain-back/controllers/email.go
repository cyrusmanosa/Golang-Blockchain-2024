package controllers

import (
	"blockchain-back/modules"
	"log"
	"net/smtp"
)

func SendRequest(req modules.InputData) {
	auth := smtp.PlainAuth("", "studiocmkc0110@gmail.com", "iodvpvmlyvadnhfb", "smtp.gmail.com")
	subject := `興味を待つ会社の方ー` + req.CompanyName + `の` + req.Name + `さん`
	body := `<html>
		<body>
			<h3>履歴書の取得をご希望の場合、確認してください。</h3>
				<h2>名前：` + req.Name + `</h2>
				<h2>会社名前：` + req.CompanyName + `</h2>
				<h2>メール：` + req.Email + `</h2>
				<h2>メッセージ：` + req.Message + `</h2>
				<h2>時間：` + req.SendTime + `</h2>
			<br>
			<a href="http://localhost:5173/Ok/` + req.Name + `"> CHECK BTN </a>
			<br>
			<br>
		</body>
	</html>`

	msg := "Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", auth, "studiocmkc0110@gmail.com", []string{"manhk123456@gmail.com"}, []byte(msg))
	if err != nil {
		log.Printf("Err: %v", err)
	}
}

func SendRsp(req modules.InputData) {
	auth := smtp.PlainAuth("", "studiocmkc0110@gmail.com", "iodvpvmlyvadnhfb", "smtp.gmail.com")
	subject := req.CompanyName + `の` + req.Name + `さん｜文家俊の書類リング`
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
				<h3>自己紹介</h3>
				<h4>
					私は文家俊と申します。今年30歳で、香港出身です。前職は動画編集者として働いていましたが、現在は日本のEccコンピュータ専門学校でIT開発研究コースに在籍しており、2025年3月に卒業予定です。
					将来は日本に定住し、フルスタックエンジニアとして活躍したいと考えています。私のキャリアはユニークで、香港の職場文化と日本の職場文化の違いを身をもって経験しました。
					香港では転職が一般的であり、私も事情があって短期間で2回の転職を経験しました。その度に昇進し、前職の社長や上司から理解を得て、推薦書を書いていただきました。
					現在学んでいるITの知識を活かし、フルスタックエンジニアとして多様なプロジェクトに携わりたいと考えています。
					日本での新しい挑戦に胸を膨らませながら、引き続き努力を重ねていく所存です。どうぞよろしくお願い申し上げます。
				</h4>
			<br>
			<p>※このメールは送信専用です。返信はしないでください。</p>
		</body>
	</html>`

	msg := "Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", auth, "studiocmkc0110@gmail.com", []string{req.Email}, []byte(msg))
	if err != nil {
		log.Printf("Err: %v", err)
	}
}
