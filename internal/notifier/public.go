package notifier

import (
	"context"
	"github.com/emortalmc/proto-specs/gen/go/model/privatemessage"
)

type Notifier interface {
	MessageSent(ctx context.Context, msg *privatemessage.PrivateMessage) error
}
