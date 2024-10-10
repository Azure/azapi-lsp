package command

import (
	"context"
	"encoding/json"
)

type CommandHandler interface {
	Handle(ctx context.Context, params []json.RawMessage) (interface{}, error)
}
