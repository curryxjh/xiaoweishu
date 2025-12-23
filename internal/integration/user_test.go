package integration

import (
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			
		})
	}
}
