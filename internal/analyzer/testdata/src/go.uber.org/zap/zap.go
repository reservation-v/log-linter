package zap

type Field struct{}

type Logger struct{}

func (l *Logger) Debug(msg string, fields ...Field) {}
func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}

func String(key, value string) Field  { return Field{} }
func Any(key string, value any) Field { return Field{} }
