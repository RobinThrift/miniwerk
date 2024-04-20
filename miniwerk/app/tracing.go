package app

import "context"

type ctxReqIDKeyType string

const ctxReqIDKey = ctxReqIDKeyType("ctxReqIDKey")

func requestIDWithCtx(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxReqIDKey, id)
}

func requestIDFromCtx(ctx context.Context) (string, bool) {
	val := ctx.Value(ctxReqIDKey)
	id, ok := val.(string)
	if ok {
		return id, true
	}

	return "", false
}
