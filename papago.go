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

func (p *Papago) sendRequest(ctx context.Context, req *request.Request, r interface{}) error {
	resp, err := req.Headers(map[string]string{
		"X-Naver-Client-Id":     p.clientID,
		"X-Naver-Client-Secret": p.clientSecret,
	}).Do(ctx)
	if err != nil {
		return err
	}

	if !resp.Success() {
		return fmt.Errorf("failed with %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	if err := resp.JSON(r); err != nil {
		return err
	}

	return nil
}

// Translate Translate
// please refer https://developers.naver.com/docs/papago/papago-nmt-api-reference.md
func (p *Papago) Translate(ctx context.Context, source, target string, text string) (string, error) {
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

	if err := p.sendRequest(ctx, request.Post("https://openapi.naver.com/v1/papago/n2mt").Forms(map[string]string{
		"source": source,
		"target": target,
		"text":   text,
	}), &r); err != nil {
		return "", err
	}

	return r.Message.Result.TranslatedText, nil
}

// DetectLangs Detect Languages
// https://developers.naver.com/docs/papago/papago-detectlangs-api-reference.md
func (p *Papago) DetectLangs(ctx context.Context, text string) (string, error) {
	var r struct {
		LangCode string `json:"langCode"`
	}

	if err := p.sendRequest(ctx, request.Post("https://openapi.naver.com/v1/papago/detectLangs").Form("query", text), &r); err != nil {
		return "", err
	}

	return r.LangCode, nil
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
	var r struct {
		Result []RomanizationResult `json:"aResult"`
	}

	if err := p.sendRequest(ctx, request.Get("https://openapi.naver.com/v1/krdict/romanization").Query("query", text), &r); err != nil {
		return nil, err
	}

	return r.Result, nil
}
