package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func NewDiscardLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
