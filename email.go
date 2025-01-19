package main

import (
	"fmt"
	"net/smtp"
)

func sendEmail(to string, subject string, body string) error {
	// การตั้งค่าผู้ส่งและเซิร์ฟเวอร์
	from := "poschananbussaba@gmail.com" // อีเมลผู้ส่ง
	password := "ncyq wwfe rpob lrzw"    // App Password จาก Google
	smtpHost := "smtp.gmail.com"         // Gmail SMTP Host
	smtpPort := "587"                    // Gmail SMTP Port

	// การสร้างข้อความอีเมล
	message := []byte("Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// การตั้งค่า SMTP Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// ส่งอีเมล
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent Successfully!")
	return nil
}

func main() {
	err := sendEmail("linesunnyname@gmail.com", "Test Subject", "This is a test email.")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
