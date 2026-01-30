package reporting

import (
	"bytes"
	"strings"
)

type testLogger interface {
	Logf(format string, args ...any)
}

type tlogReporter struct {
	test testLogger
}

func (self *tlogReporter) BeginStory(story *StoryReport) {
	if story == nil || story.Test == nil {
		self.test = nil
		return
	}
	if logger, ok := story.Test.(testLogger); ok {
		self.test = logger
	}
}

func (self *tlogReporter) Enter(scope *ScopeReport) {
	if self.test == nil || scope == nil {
		return
	}
	self.test.Logf("Convey: %s", scope.Title)
}

func (self *tlogReporter) Report(r *AssertionResult) {
	if self.test == nil || r == nil {
		return
	}
	if r.Error != nil {
		self.test.Logf("ERROR: %v (%s:%d)", r.Error, r.File, r.Line)
		if r.StackTrace != "" {
			self.test.Logf("%s", r.StackTrace)
		}
		return
	}
	if r.Failure != "" {
		self.test.Logf("FAIL: %s (%s:%d)", r.Failure, r.File, r.Line)
		if r.Expected != "" || r.Actual != "" {
			self.test.Logf("expected: %s", r.Expected)
			self.test.Logf("actual:   %s", r.Actual)
		}
		if r.StackTrace != "" {
			self.test.Logf("%s", r.StackTrace)
		}
		return
	}
	if r.Skipped {
		self.test.Logf("SKIP (%s:%d)", r.File, r.Line)
	}
}

func (self *tlogReporter) Exit() {}

func (self *tlogReporter) EndStory() {
	self.test = nil
}

func (self *tlogReporter) Write(content []byte) (written int, err error) {
	if self.test == nil {
		return len(content), nil
	}
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return len(content), nil
	}
	for _, line := range strings.Split(string(trimmed), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		self.test.Logf("%s", line)
	}
	return len(content), nil
}

func NewTLogReporter() *tlogReporter {
	return new(tlogReporter)
}
