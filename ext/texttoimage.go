package ext

import "bytes"

type ImgGenratationSrv interface {
	GenerateImg(text string) (*bytes.Buffer, error)
}
