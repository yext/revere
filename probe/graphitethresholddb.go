package probe

// GraphiteThresholdDBModel defines the JSON serialization format for saving
// Graphite threshold probes' settings in the database.
type GraphiteThresholdDBModel struct {
	SourceID   int64
	Expression string

	Thresholds GraphiteThresholdThresholdsDBModel
	TriggerIf  string

	CheckPeriodMilli int64

	TimeToAuditMilli        int64
	RecentTimeToIgnoreMilli int64
	AuditFunction           string
}

// GraphiteThresholdThresholdsDBModel defines the JSON serialization format for
// saving Graphite threshold probes' threshold settings in the database.
type GraphiteThresholdThresholdsDBModel struct {
	Warning  *float64
	Error    *float64
	Critical *float64
}
