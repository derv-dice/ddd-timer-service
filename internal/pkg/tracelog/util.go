package tracelog

import "github.com/rs/zerolog"

func (t *TraceLogger) Trace(eventName string, attributes ...KeyValue) {
	t.AddEvent(eventName, zerolog.TraceLevel, attributes...)
}

func (t *TraceLogger) Debug(eventName string, attributes ...KeyValue) {
	t.AddEvent(eventName, zerolog.DebugLevel, attributes...)
}

func (t *TraceLogger) Info(eventName string, attributes ...KeyValue) {
	t.AddEvent(eventName, zerolog.InfoLevel, attributes...)
}

func (t *TraceLogger) InfoWithDuration(eventName string, attributes ...KeyValue) {
	if attributes == nil {
		attributes = []KeyValue{
			String(spanTimeElapsedLogKey, t.Duration().String()),
		}

		t.AddEvent(eventName, zerolog.InfoLevel, attributes...)
		return
	}

	attributes = append(attributes, String(spanTimeElapsedLogKey, t.Duration().String()))
	t.AddEvent(eventName, zerolog.InfoLevel, attributes...)
}

func (t *TraceLogger) Warn(eventName string, attributes ...KeyValue) {
	t.AddEvent(eventName, zerolog.WarnLevel, attributes...)
}

func (t *TraceLogger) Error(err error) {
	t.AddError(err)
}
