package test

const (
	DefaultProbeJson = `{
		"url": "foo.bar",
		"expression": "test",
		"warningThreshold": 1000,
		"errorThreshold": 1200,
		"criticalThreshold": 1500,
		"auditFunction": "max",
		"checkPeriod": 5,
		"checkPeriodType": "hour",
		"triggerIf": ">",
		"auditPeriod": 10,
		"auditPeriodType": "minute"
	}`
	DefaultTargetJson = `{
		"emails": [
			{"emailTo":"test@ex.com", "replyTo":"test2@ex.com"}
		]
	}`
)

var (
	PeriodTypes = []string{
		"second",
		"minute",
		"hour",
		"day",
	}
)
