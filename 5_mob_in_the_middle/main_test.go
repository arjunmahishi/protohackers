package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_replaceToken(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{
			msg:  "send money on 7F1u3wSD5RbOHQmupo9nx4TnhQ",
			want: fmt.Sprintf("send money on %s", tonyToken),
		},
		{
			msg:  "7F1u3wSD5RbOHQmupo9nx4TnhQ is where you should send money",
			want: fmt.Sprintf("%s is where you should send money", tonyToken),
		},
		{
			msg:  "Token: 7F1u3wSD5RbOHQmupo9nx4TnhQ -- Send money",
			want: fmt.Sprintf("Token: %s -- Send money", tonyToken),
		},
		{
			msg:  "this is not a token: 7YWHMfk9JZe0LM0g1ZauHuiSxhIgCwlev0DzwehOnhYw-1234",
			want: "this is not a token: 7YWHMfk9JZe0LM0g1ZauHuiSxhIgCwlev0DzwehOnhYw-1234",
		},
		{
			msg:  "multiple tokens: 7F1u3wSD5RbOHQmupo9nx4TnhQ 7F1u3wSD5RbOHQmupo9nx4abcd",
			want: fmt.Sprintf("multiple tokens: %s %s", tonyToken, tonyToken),
		},
		{
			msg:  "76h86SQWrG7B1vlEBMcncmlBfxJQcE5h9d-n3UmdLqVykh4gU1kjCIaljv1fg-1234",
			want: "76h86SQWrG7B1vlEBMcncmlBfxJQcE5h9d-n3UmdLqVykh4gU1kjCIaljv1fg-1234",
		},
		{
			msg:  "n3UmdLqVykh4gU1kjCIaljv1fg-76h86SQWrG7B1vlEBMcncmlBfxJQcE5h9d",
			want: "n3UmdLqVykh4gU1kjCIaljv1fg-76h86SQWrG7B1vlEBMcncmlBfxJQcE5h9d",
		},
		{
			msg:  "no token",
			want: "no token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			assert.Equal(t, tt.want, replaceToken(tt.msg))
		})
	}
}
