package http

import (
	"net/http"
	"testing"
)

func TestMissingParamError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "happy path",
			fields: fields{Name: "v"},
			want:   "named URL param `v` is missing",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := MissingParamError{
				Name: tt.fields.Name,
			}

			if got := e.Error(); got != tt.want {
				t.Errorf("MissingParamError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMissingParamError_ToErrorResponse(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name string
	}

	tests := []struct {
		name   string
		fields fields
		want   *ErrorResponse
	}{
		{
			name:   "happy path",
			fields: fields{Name: "v"},
			want: &ErrorResponse{
				Error:          MissingParamError{Name: "v"},
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "missing_param_error",
				ErrorMsg:       "named URL param `v` is missing",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := MissingParamError{
				Name: tt.fields.Name,
			}

			if got := e.ToErrorResponse(); !got.IsEqual(tt.want) {
				t.Errorf("MissingParamError.ToErrorResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsingParamError_Error(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name  string
		Value string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "happy path",
			fields: fields{Name: "v", Value: "xxx"},
			want:   "can not parse named URL param `v`: `xxx` is invalid",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ParsingParamError{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("MissingParamError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsingParamError_ToErrorResponse(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name  string
		Value string
	}

	tests := []struct {
		name   string
		fields fields
		want   *ErrorResponse
	}{
		{
			name:   "happy path",
			fields: fields{Name: "v", Value: "xxx"},
			want: &ErrorResponse{
				Error:          ParsingParamError{Name: "v", Value: "xxx"},
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "parsing_param_error",
				ErrorMsg:       "can not parse named URL param `v`: `xxx` is invalid",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := ParsingParamError{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}

			if got := e.ToErrorResponse(); !got.IsEqual(tt.want) {
				t.Errorf("MissingParamError.ToErrorResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
