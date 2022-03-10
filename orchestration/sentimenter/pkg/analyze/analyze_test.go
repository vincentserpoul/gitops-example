package analyze

import "testing"

func TestAnalyzeFunc(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want uint8
	}{
		{
			name: "tldr",
			args: args{s: "welcome to haiti. We are fun"},
			want: 1,
		},
	}

	analy, err := AnalyzeFunc()
	if err != nil {
		t.Errorf("analyze() error = %v", err)

		return
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := analy(tt.args.s)
			if got != tt.want {
				t.Errorf("tldrSentence() = %v, want %v", got, tt.want)
			}
		})
	}
}
