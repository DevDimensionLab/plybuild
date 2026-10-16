//go:build windows

package taskrun

// The public CLI is portable; the qualified same-owner provider runtime is not.
func ReadProductQuestion(Dependencies, string, string) (ProductQuestionReadback, error) {
	return ProductQuestionReadback{}, boundAnswerError("unsupported", "Product-answer continuation requires the qualified macOS worker adapter.")
}

func PreviewBoundAnswer(Dependencies, string, string, string, int, string, string) (BoundAnswerPreview, error) {
	return BoundAnswerPreview{}, boundAnswerError("unsupported", "Product-answer continuation requires the qualified macOS worker adapter.")
}

func ApplyBoundAnswer(Dependencies, BoundAnswerRequest) (BoundAnswerStatus, error) {
	return BoundAnswerStatus{}, boundAnswerError("unsupported", "Product-answer continuation requires the qualified macOS worker adapter.")
}

func ShowBoundAnswer(Dependencies, string, string) (BoundAnswerStatus, error) {
	return BoundAnswerStatus{}, boundAnswerError("unsupported", "Product-answer continuation requires the qualified macOS worker adapter.")
}
