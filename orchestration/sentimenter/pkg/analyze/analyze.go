package analyze

import (
	"fmt"

	"github.com/cdipaolo/sentiment"
)

func AnalyzeFunc() (func(string) uint8, error) {
	model, err := sentiment.Restore()
	if err != nil {
		return nil, fmt.Errorf("could not restore model: %w", err)
	}

	return func(s string) uint8 {
		analysis := model.SentimentAnalysis(s, sentiment.English)

		return analysis.Score
	}, nil
}
