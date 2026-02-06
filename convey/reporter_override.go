package convey

import (
	"github.com/jtolds/gls"
	"github.com/smartystreets/goconvey/convey/reporting"
)

type ReporterFactory func() reporting.Reporter

const reporterFactoryKey = "reporterFactory"

func WithReporterFactory(factory ReporterFactory, action func()) {
	if factory == nil {
		action()
		return
	}
	ctxMgr.SetValues(gls.Values{reporterFactoryKey: factory}, action)
}

func BuildQuietReporter() reporting.Reporter {
	return reporting.NewReporters(
		reporting.NewGoTestReporter(),
		newNilReporter(),
	)
}
