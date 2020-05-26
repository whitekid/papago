package papago

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransalte(t *testing.T) {
	type args struct {
		source, target Lang
		text           string
	}

	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		translated string
	}{
		{"", args{Ko, En, "만나서 반갑습니다."}, false, "Nice to meet you."},
		{"", args{En, Ko, "Nice to meet you."}, false, "만나서 반가워"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := New(os.Getenv("NAVER_CLIENT_ID"), os.Getenv("NAVER_CLIENT_SECRET"))
			translated, err := api.Translate(tt.args.source, tt.args.target, tt.args.text)
			if (err != nil) != tt.wantErr {
				require.Failf(t, "translated failed", "error = %v, wantErr = %v", err, tt.wantErr)
			}
			require.Equal(t, tt.translated, translated)
		})
	}
}
