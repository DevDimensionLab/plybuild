package taskjournal

import "time"

// elapsedSeconds avoids time.Duration's roughly 292-year limit. Subtract the
// integer components before converting to float64 so nearby timestamps retain
// their nanoseconds, including across a second boundary. Normalize the remainder
// as Duration.Seconds does to preserve ordinary duration values and digests.
// RFC3339's supported years fit comfortably in an int64 count of seconds.
func elapsedSeconds(start, finish time.Time) float64 {
	seconds := finish.Unix() - start.Unix()
	nanos := int64(finish.Nanosecond()) - int64(start.Nanosecond())
	if seconds > 0 && nanos < 0 {
		seconds--
		nanos += 1e9
	} else if seconds < 0 && nanos > 0 {
		seconds++
		nanos -= 1e9
	}
	return float64(seconds) + float64(nanos)/1e9
}
