package main

import (
	"backend_spring_2024/internal/assert"
	"testing"
)

func TestIsRole(t *testing.T) {
	tests := []struct {
		name      string
		role      string
		secretKey []byte
		wantRole  string
		want      bool
	}{
		{
			name:      "correct admin token",
			role:      "admin",
			secretKey: []byte("correctSecretKey"),
			wantRole:  "admin",
			want:      true,
		},
		{
			name:      "correct user token",
			role:      "user",
			secretKey: []byte("correctSecretKey"),
			wantRole:  "user",
			want:      true,
		},
		{
			name:      "incorrect admin token",
			role:      "dog",
			secretKey: []byte("correctSecretKey"),
			wantRole:  "admin",
			want:      false,
		},
		{
			name:      "incorrect user token",
			role:      "sea",
			secretKey: []byte("correctSecretKey"),
			wantRole:  "user",
			want:      false,
		},
		{
			name:      "incorrect token",
			role:      "admin",
			secretKey: []byte("incorrectSecretKey"),
			wantRole:  "admin",
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := GenerateToken(tt.role, []byte("correctSecretKey"))
			assert.Equal(t, IsRole(token, tt.wantRole, tt.secretKey), tt.want)
		})
	}

}
