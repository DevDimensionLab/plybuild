package currentstate

import "strings"

// Compatibility reads for records written before explicit observation origins.
// Never use these to encode new state or rewrite frozen historical records.
func LegacyProviderWait(w Worker) bool {
	return w.RetryNotBefore != "" && w.Reason == "usage_limit: same Run waits until "+w.RetryNotBefore
}

func LegacyWeekOrigins(s Week, e Entry) (Origin, Origin) {
	guidance, block := s.GuidanceOrigin, e.BlockOrigin
	if guidance == "" {
		guidance = OtherOrigin
		if s.NextAction == ProductAnswerWait || strings.HasPrefix(s.NextAction, NativeNextPrefix) {
			guidance = NativeRun
		}
	}
	if block == "" {
		block = OtherOrigin
		if e.State == "blocked" && s.Phase == "unknown" && strings.HasPrefix(s.NextAction, NativeNextPrefix) {
			block = AnswerBlock
		}
	}
	return guidance, block
}
