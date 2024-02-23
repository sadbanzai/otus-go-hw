package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	BadValidationParams struct {
		Min int    `validate:"min:min"`
		Max int    `validate:"max:max"`
		Len string `validate:"len:len"`
		In  int    `validate:"in:x,y"`
	}
)

var user = User{
	ID:     "1234567890",
	Name:   "John Dow",
	Age:    52,
	Email:  "123",
	Role:   "admin",
	Phones: []string{"+1234567890", "+123456789"},
	meta:   []byte{},
}

var app = App{Version: "1.0.0-rc1"}

var response = Response{
	Code: 503,
	Body: "Service Unavailable",
}

var badvalparams = BadValidationParams{
	Min: 1,
	Max: 1,
	Len: "string",
	In:  1,
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			user,
			ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLen},
				ValidationError{Field: "Age", Err: ErrMax},
				ValidationError{Field: "Email", Err: ErrRegexp},
				ValidationError{Field: "Phones 1", Err: ErrLen},
			},
		},
		{
			app,
			ValidationErrors{
				ValidationError{Field: "Version", Err: ErrLen},
			},
		},
		{
			Token{},
			nil,
		},
		{
			response,
			ValidationErrors{
				ValidationError{Field: "Code", Err: ErrIn},
			},
		},
		{
			"some string",
			ErrNotStruct,
		},
		{
			badvalparams,
			ValidationErrors{
				ValidationError{Field: "Min", Err: ErrBadParam},
				ValidationError{Field: "Max", Err: ErrBadParam},
				ValidationError{Field: "Len", Err: ErrBadParam},
				ValidationError{Field: "In", Err: ErrBadParam},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tt.expectedErr, err)
			}
			_ = tt
		})
	}
}
