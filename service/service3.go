package service

import (
	"cc-service3/ext"
	"cc-service3/storage"
	"fmt"
	"log"
	"time"
)

const emailMessageFormat = `Hello,
The picture that you requested is ready to download!
You can download the picture from the linke bellow:
%s`

type Service3 struct {
	DataBase  storage.RequestDB
	PicStore  storage.ImageStorage
	TextToImg ext.HuggingFace
	MSender   ext.MailSender
}

func NewService3(db storage.RequestDB, imgstore storage.ImageStorage, txtimg ext.HuggingFace, msender ext.MailSender) *Service3 {
	return &Service3{
		DataBase:  db,
		PicStore:  imgstore,
		TextToImg: txtimg,
		MSender:   msender,
	}
}

func (s Service3) Execute() error {
	ticker := time.Tick(5 * time.Second)
	for range ticker {
		log.Println("ticker started")
		readyReqs, err := s.DataBase.GetAllReadies()
		if err != nil {
			return fmt.Errorf("error at reading database: %w", err)
		}

		for _, r := range readyReqs {
			fmt.Println("----------")
			log.Println("generating image for ", r.ReqId)
			imgfile, err := s.TextToImg.GenerateImg(r.ImageCaption)
			if err != nil {
				return fmt.Errorf("error at image generation: %w", err)
			}

			// TODO: format as jpeg and upload as r.ReqId + ".jpg"
			log.Println("uploading image ", r.ReqId)
			imgurl, err := s.PicStore.Upload(imgfile, imgfile.Size(), r.ReqId)
			if err != nil {
				return fmt.Errorf("error at uploading image: %w", err)
			}

			log.Printf("setting url for %s: %s\n", r.ReqId, imgurl)
			err = s.DataBase.SetImageURL(r.ReqId, imgurl)
			if err != nil {
				return fmt.Errorf("error at database updating: %w", err)
			}

			log.Println("setting status final for ", r.ReqId)
			err = s.DataBase.SetStatus(r.ReqId, "done")
			if err != nil {
				return fmt.Errorf("error at database updating: %w", err)
			}

			log.Println("sending email")
			mail := ext.Email{
				RecipientMail: r.Email,
				Subject:       "CCS3: Requested Picture is Ready",
				Message:       fmt.Sprintf(emailMessageFormat, imgurl),
			}
			err = s.MSender.Send(mail)
			if err != nil {
				fmt.Println("failed to send notifying email")
			}
		}
	}
	return nil
}
