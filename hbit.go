package hbit

import "context"

// Report critical errors - currently a no-op
var ReportError = func(ctx context.Context, err error, args ...any) {}
