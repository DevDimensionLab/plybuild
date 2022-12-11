package logger

import "github.com/sirupsen/logrus"

type Collector struct {
	entries []*logrus.Entry
}

func (collector *Collector) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel, logrus.WarnLevel}
}

func (collector *Collector) Fire(entry *logrus.Entry) error {
	collector.entries = append(collector.entries, entry)
	return nil
}
