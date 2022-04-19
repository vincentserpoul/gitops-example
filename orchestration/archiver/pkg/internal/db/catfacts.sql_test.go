package db

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestQueries(t *testing.T) {
	id1 := uuid.New()

	tests := []struct {
		name    string
		insert  SaveHappycatFactParams
		id      uuid.UUID
		expect  HappycatFact
		wantErr bool
	}{
		{
			name: "happypath",
			insert: SaveHappycatFactParams{
				ID:   id1,
				Fact: "toto",
			},
			id: id1,
			expect: HappycatFact{
				ID:   id1,
				Fact: "toto",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := q.SaveHappycatFact(context.Background(), tt.insert); (err != nil) != tt.wantErr {
				t.Errorf("Queries.SaveHappycatFact() error = %v, wantErr %v", err, tt.wantErr)
			}

			cf, err := q.GetHappycatFact(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Queries.GetHappycatFact() error = %v, wantErr %v", err, tt.wantErr)
			}

			if reflect.DeepEqual(cf, &tt.expect) {
				t.Errorf("missing cf = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
