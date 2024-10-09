package main

import (
	"cc-service3/ext"
	"cc-service3/service"
	"cc-service3/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	acs3_accessKey := os.Getenv("CCSERV3_ACS3_ACCESSKEY")
	acs3_secretKey := os.Getenv("CCSERV3_ACS3_SECRETKEY")
	imagestore, err := storage.NewArvanCloudS3("cc-practice-004-out", "ir-thr-at1",
		"https://s3.ir-thr-at1.arvanstorage.com",
		acs3_accessKey,
		acs3_secretKey)
	if err != nil {
		log.Fatalf("Fatal error at object storage: %v", err)
	}

	mySQL_username := os.Getenv("CCSERV3_MYSQL_USERNAME")
	mySQL_password := os.Getenv("CCSERV3_MYSQL_PASSWORD")
	mySQL_address := "mysql-container:3306"
	database, err := storage.NewMySQLDB(mySQL_username, mySQL_password, mySQL_address, "ccp1")
	if err != nil {
		log.Fatalf("Fatal error at database: %v", err)
	}

	texttoimage := ext.NewHuggingFace(os.Getenv("CCSERV3_HF_APIKEY"))

	mailersend_apikey := os.Getenv("CCSERV3_MS_APIKEY")
	msservice := ext.NewMailerSend(mailersend_apikey,
		"no-reply@trial-jy7zpl96xv5l5vx6.mlsender.net",
		"Pouria")

	srv := service.NewService3(*database, *imagestore, *texttoimage, *msservice)
	err = srv.Execute()

	if err != nil {
		log.Fatal("Could not start service1")
	}
}
