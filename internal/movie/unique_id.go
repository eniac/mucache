package movie

import (
	"context"
	"github.com/lithammer/shortuuid"
)

func GetUniqueId(ctx context.Context, reqId string) string {
	reviewId := shortuuid.New()
	return reviewId
}
