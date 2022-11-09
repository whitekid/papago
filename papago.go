package papago

import (
	"context"
	"fmt"

	"github.com/whitekid/goxp/request"
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

// Translate Translate
// please refer https://developers.naver.com/docs/papago/papago-nmt-api-reference.md
func (p *Papago) Translate(ctx context.Context, source, target string, text string) (string, error) {
	resp, err := request.Post("https://openapi.naver.com/v1/papago/n2mt").
		Headers(map[string]string{
			"X-Naver-Client-Id":     p.clientID,
			"X-Naver-Client-Secret": p.clientSecret,
		}).
		Forms(map[string]string{
			"source": source,
			"target": target,
			"text":   text,
		}).
		Do(ctx)
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", fmt.Errorf("error with %d", resp.StatusCode)
	}

	var r struct {
		Message struct {
			Type    string `json:"@type"`
			Service string `json:"@service"`
			Version string `json:"@version"`
			Result  struct {
				SrcLangType    string `json:"srcLangType"`
				TarLangType    string `json:"tarLangType"`
				TranslatedText string `json:"translatedText"`
			} `json:"result"`
		} `json:"message"`
	}
	defer resp.Body.Close()
	if err := resp.JSON(&r); err != nil {
		return "", err
	}

	return r.Message.Result.TranslatedText, nil
}

// DetectLangs Detect Languages
// https://developers.naver.com/docs/papago/papago-detectlangs-api-reference.md
func (p *Papago) DetectLangs(ctx context.Context, text string) (string, error) {
	resp, err := request.Post("https://openapi.naver.com/v1/papago/detectLangs").
		Headers(map[string]string{
			"X-Naver-Client-Id":     p.clientID,
			"X-Naver-Client-Secret": p.clientSecret,
		}).
		Form("query", text).
		Do(ctx)
	if err != nil {
		return "", err
	}

	if !resp.Success() {
		return "", fmt.Errorf("error with %d", resp.StatusCode)
	}

	var r struct {
		LanCode string `json:"langCode"`
	}
	defer resp.Body.Close()
	if err := resp.JSON(&r); err != nil {
		return "", err
	}

	return r.LanCode, nil
}

type RomanizationResult struct {
	FirstName string                   `json:"sFirstName"`
	Items     []RomanizationResultItem `json:"aItems"`
}

type RomanizationResultItem struct {
	Name  string `json:"name"`
	Score string `json:"score"`
}

// Romanization 한글인명 로마자 변환
// https://developers.naver.com/docs/papago/papago-romanization-api-reference.md
func (p *Papago) Romanization(ctx context.Context, text string) ([]RomanizationResult, error) {
	resp, err := request.Get("https://openapi.naver.com/v1/krdict/romanization").
		Headers(map[string]string{
			"X-Naver-Client-Id":     p.clientID,
			"X-Naver-Client-Secret": p.clientSecret,
		}).
		Query("query", text).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("error with status %d", resp.StatusCode)
	}

	var r struct {
		Result []RomanizationResult `json:"aResult"`
	}
	defer resp.Body.Close()
	if err := resp.JSON(&r); err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", r)

	return r.Result, nil
}
