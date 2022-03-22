package http

import "context"

type NamedURLParamsGetter func(ctx context.Context, key string) string
