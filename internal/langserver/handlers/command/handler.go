package command

import "context"

type CommandHandler interface {
	Handle(ctx context.Context, params CommandArgs) (interface{}, error)
}
