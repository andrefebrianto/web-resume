package errorx

type uncaughtPanic struct{ message string }

func (p uncaughtPanic) Error() string {
	return p.message
}
