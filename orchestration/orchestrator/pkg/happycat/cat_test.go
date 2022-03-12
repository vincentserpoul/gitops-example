package happycat

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func Test_getCatFacts(t *testing.T) {
	t.Parallel()

	type args struct {
		timeout time.Duration
	}

	tests := []struct {
		name           string
		serverResponse string
		args           args
		want           CatFacts
		errIs          error
	}{
		{
			name: "happy path",
			serverResponse: `
			[
				{
					"status": {
						"verified": true,
						"feedback": "",
						"sentCount": 1
					},
					"_id": "5887e1d85c873e0011036889",
					"user": "5a9ac18c7478810ea6c06381",
					"text": "Cats make about 100 different sounds. Dogs make only about 10.",
					"__v": 0,
					"source": "user",
					"updatedAt": "2020-09-03T16:39:39.578Z",
					"type": "cat",
					"createdAt": "2018-01-15T21:20:00.003Z",
					"deleted": false,
					"used": true
				},
				{
					"status": {
						"verified": true,
						"sentCount": 1
					},
					"_id": "588e746706ac2b00110e59ff",
					"user": "588e6e8806ac2b00110e59c3",
					"text": "Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming.",
					"__v": 0,
					"source": "user",
					"updatedAt": "2020-08-26T20:20:02.359Z",
					"type": "cat",
					"createdAt": "2018-01-14T21:20:02.750Z",
					"deleted": false,
					"used": true
				}
			]`,
			args: args{
				timeout: 0,
			},
			want: CatFacts{
				"Cats make about 100 different sounds. Dogs make only about 10.",
				"Domestic cats spend about 70 percent of the day sleeping and 15 percent of the day grooming.",
			},
			errIs: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			trcr := trace.NewNoopTracerProvider().Tracer("test")

			s := testServer(tt.serverResponse, tt.args.timeout, 1)
			defer s.Close()

			got, err := getCatFacts(ctx, trcr, tt.args.timeout, s.URL)
			if !errors.Is(tt.errIs, err) {
				t.Errorf("getCatFacts() error = %v, wantErr %v", err, tt.errIs)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCatFacts() = %v, want %v", got, tt.want)
			}
		})
	}
}
