package sqlx

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonColumn_Value(t *testing.T) {
	testcases := []struct {
		name    string
		valuer  driver.Valuer
		wantRes any
		wantErr error
	}{
		{
			name: "user",
			valuer: JsonColumn[User]{
				Valid: true,
				Val: User{
					Name: "test",
				},
			},
			wantRes: []byte(`{"Name":"test"}`),
		},
		{
			name:   "invalid",
			valuer: JsonColumn[User]{},
		},
		{
			name:   "nil",
			valuer: JsonColumn[*User]{},
		},
		{
			name:    "nil but valid",
			valuer:  JsonColumn[*User]{Valid: true},
			wantRes: []byte(`null`),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.valuer.Value()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, value)
		})
	}
}

func TestJsonColumn_Scan(t *testing.T) {
	testcases := []struct {
		name      string
		src       any
		wantErr   error
		wantValid bool
		wantVal   User
	}{
		{
			name:    "nil",
			wantVal: User{},
		},
		{
			name:      "string",
			src:       `{"Name":"test"}`,
			wantErr:   nil,
			wantValid: true,
			wantVal:   User{Name: "test"},
		},
		{
			name:      "bytes",
			src:       []byte(`{"Name":"test"}`),
			wantErr:   nil,
			wantValid: true,
			wantVal:   User{Name: "test"},
		},
		{
			name:      "int",
			src:       123,
			wantErr:   errors.New("ekit: json column can only scan string or []byte, but got int"),
			wantValid: false,
			wantVal:   User{},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			js := JsonColumn[User]{}
			err := js.Scan(tc.src)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantValid, js.Valid)
			if !js.Valid {
				return
			}
			assert.Equal(t, tc.wantVal, js.Val)
		})
	}
}

func TestJsonColumn_base(t *testing.T) {
	ExampleJsonColumn_Value()
	ExampleJsonColumn_Scan()
}

type User struct {
	Name string
}

func ExampleJsonColumn_Value() {
	js := JsonColumn[User]{
		Valid: true,
		Val: User{
			Name: "test",
		},
	}
	value, err := js.Value()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(value.([]byte)))
	// Output: {"Name":"test"}
}

func ExampleJsonColumn_Scan() {
	js := JsonColumn[User]{}
	err := js.Scan([]byte(`{"Name":"test"}`))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(js.Val)
	// Output: {test}
}
