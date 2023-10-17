package ral

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/ecodeclub/notify-go/pkg/log"
)

func Test_ral(t *testing.T) {
	c := NewClient("getui")
	req := Request{}
	var resp Result

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})).With("logID", 12345678)
	_ = c.Ral(context.WithValue(context.TODO(), log.ContextLogKey{}, l), "Auth", req, &resp, &map[string]any{})
}

type Result struct {
	Msg  string
	Code int
	Data Data
}

type Data struct {
	ExpireTime string `json:"expire_time"`
	Token      string
}
