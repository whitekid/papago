package papago

import (
	"fmt"
	"net/http"

	"github.com/whitekid/go-utils/request"
)

type Lang string

const (
	Ko   Lang = "ko"
	En   Lang = "en"
	ZhCN Lang = "zh-CN"
	ZhTW Lang = "zh-TW"
	Es   Lang = "es"
	Fr   Lang = "fr"
	Vi   Lang = "vi"
	Id   Lang = "id"
)

func New(clientID, clientSecret string) *Papago {
	return &Papago{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

type Papago struct {
	clientID     string
	clientSecret string
}

func (p *Papago) Translate(source, target Lang, text string) (string, error) {
	resp, err := request.Post("https://openapi.naver.com/v1/papago/n2mt").
		Headers(map[string]string{
			"X-Naver-Client-Id":     p.clientID,
			"X-Naver-Client-Secret": p.clientSecret,
		}).
		Forms(map[string]string{
			"source": string(source),
			"target": string(target),
			"text":   text,
		}).
		Do()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error with %d", resp.StatusCode)
	}

	var r struct {
		Message struct {
			Type    string `json:"@type"`
			Service string `json:"@service"`
			Version string `json:"@version"`
			Result  struct {
				TranslatedText string `json:"translatedText"`
				SrcLangType    string `json:"srcLangType"`
			} `json:"result"`
		} `json:"message"`
	}
	if err := resp.JSON(&r); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return r.Message.Result.TranslatedText, nil
}
