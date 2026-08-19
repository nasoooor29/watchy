package models_test

import (
	"backend/models"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want *models.EnvVars
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.GetEnv()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
