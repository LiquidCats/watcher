package log

type Logger[F any] interface {
	Debug(msg string, fields ...F)
	Info(msg string, fields ...F)
	Warn(msg string, fields ...F)
	Error(msg string, fields ...F)
	Fatal(msg string, fields ...F)
}

type NamedLogger[F any] interface {
	Named(name string) Logger[F]
}
