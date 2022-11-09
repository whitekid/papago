package papago

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransalte(t *testing.T) {
	type args struct {
		source, target string
		text           string
	}

	tests := [...]struct {
		name           string
		args           args
		wantErr        bool
		wantTranslated string
	}{
		{"", args{"ko", "en", "만나서 반갑습니다."}, false, "Good to meet you."},
		{"", args{"en", "ko", "Nice to meet you."}, false, "만나서 반가워요."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			api := New(os.Getenv("NAVER_CLIENT_ID"), os.Getenv("NAVER_CLIENT_SECRET"))
			got, err := api.Translate(ctx, tt.args.source, tt.args.target, tt.args.text)
			if (err != nil) != tt.wantErr {
				require.Failf(t, "translated failed", "error = %v, wantErr = %v", err, tt.wantErr)
			}
			require.Equal(t, tt.wantTranslated, got)
		})
	}
}

func TestDetect(t *testing.T) {
	type args struct {
		text string
	}

	tests := [...]struct {
		name     string
		args     args
		wantErr  bool
		wantLang string
	}{
		{"", args{"만나서 반갑습니다."}, false, "ko"},
		{"", args{"Nice to meet you."}, false, "en"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			api := New(os.Getenv("NAVER_CLIENT_ID"), os.Getenv("NAVER_CLIENT_SECRET"))
			got, err := api.DetectLangs(ctx, tt.args.text)
			if (err != nil) != tt.wantErr {
				require.Failf(t, "DetectLangs failed", "error = %v, wantErr = %v", err, tt.wantErr)
			}
			require.Equal(t, tt.wantLang, got)
		})
	}
}

func TestRomanization(t *testing.T) {
	type args struct {
		text string
	}

	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		wantResult []RomanizationResult
	}{
		{"", args{"김정환"}, false, []RomanizationResult{{
			FirstName: "김",
			Items: []RomanizationResultItem{
				{Name: "Kim Junghwan", Score: "99"},
				{Name: "Kim Jeonghwan", Score: "70"},
				{Name: "Kim Jungwhan", Score: "39"},
				{Name: "Kim Jeongwhan", Score: "27"},
				{Name: "Kim Jenghwan", Score: "21"},
				{Name: "Kim Jengwhan", Score: "8"},
			},
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			api := New(os.Getenv("NAVER_CLIENT_ID"), os.Getenv("NAVER_CLIENT_SECRET"))
			got, err := api.Romanization(ctx, tt.args.text)
			if (err != nil) != tt.wantErr {
				require.Failf(t, "Romanization failed", "error = %v, wantErr = %v", err, tt.wantErr)
			}
			require.Equal(t, tt.wantResult, got)
		})
	}
}
