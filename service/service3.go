package service

import (
	"cc-service3/ext"
	"cc-service3/storage"
	"errors"
	"fmt"
	"log"
	"time"
)

const emailMessageFormat = `Hello,
The picture that you requested is ready to download!
You can download the picture from the link below:
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
	log.Println("starting ticker...")
	for range ticker {
		readyReqs, err := s.DataBase.GetAllReadies()
		if err != nil {
			return fmt.Errorf("error at reading database: %w", err)
		}

		for _, r := range readyReqs {
			log.Println("--- generating image for ", r.ReqId)
			imgfile, err := s.TextToImg.GenerateImg(r.ImageCaption)
			if err != nil {
				s.failureHandler(r.ReqId)
				return fmt.Errorf("error at image generation: %w", err)
			}

			log.Println("uploading image ", r.ReqId)
			imgurl, err := s.PicStore.Upload(imgfile, imgfile.Size(), r.ReqId)
			if err != nil {
				s.failureHandler(r.ReqId)
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
				Subject:       fmt.Sprintf("CCS3: Picture for Requested %s is Ready", r.ReqId),
				Message:       fmt.Sprintf(emailMessageFormat, imgurl),
			}
			err = s.MSender.Send(mail)
			if err != nil {
				log.Println("failed to send notifying email")
			}
		}
	}
	return nil
}

func (s Service3) failureHandler(requstId string) {
	err := s.DataBase.SetStatus(requstId, "failure")
	if err != nil {
		var notfound *storage.RequestNotFoundError
		if !errors.As(err, &notfound) {
			fmt.Printf("couldn't update request %s status to \"failed\": %v\n", requstId, err)
		}
		return
	}
	fmt.Printf("request %s status is set to 'failure'", requstId)
}
