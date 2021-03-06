package libkb

import (
	"fmt"
	"github.com/keybase/client/go/logger"
	"golang.org/x/net/context"
	"strings"
)

type withLogTagKey string

func WithLogTag(ctx context.Context, k string) context.Context {
	ctx = logger.ConvertRPCTagsToLogTags(ctx)

	addLogTags := true
	tagKey := withLogTagKey(k)

	if tags, ok := logger.LogTagsFromContext(ctx); ok {
		if _, found := tags[tagKey]; found {
			addLogTags = false
		}
	}

	if addLogTags {
		newTags := make(logger.CtxLogTags)
		newTags[tagKey] = k
		ctx = logger.NewContextWithLogTags(ctx, newTags)
	}

	if _, found := ctx.Value(tagKey).(withLogTagKey); !found {
		tag := RandStringB64(3)
		ctx = context.WithValue(ctx, tagKey, tag)
	}
	return ctx
}

func LogTagsToString(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	tags, ok := logger.LogTagsFromContext(ctx)
	if !ok || len(tags) == 0 {
		return ""
	}
	var out []string
	for key, tag := range tags {
		if v := ctx.Value(key); v != nil {
			out = append(out, fmt.Sprintf("%s=%s", tag, v))
		}
	}
	return strings.Join(out, ",")
}

func CopyTagsToBackground(ctx context.Context) context.Context {
	ret := context.Background()
	if tags, ok := logger.LogTagsFromContext(ctx); ok {
		ret = logger.NewContextWithLogTags(ret, tags)
		for key := range tags {
			if ctxKey, ok := key.(withLogTagKey); ok {
				if val, ok := ctx.Value(ctxKey).(string); ok && len(val) > 0 {
					ret = context.WithValue(ret, ctxKey, val)
				}
			}
		}
	}
	return ret
}
