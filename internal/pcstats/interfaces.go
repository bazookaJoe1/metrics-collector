package pcstats

// We have to make this interface general for all packages that use Metric for compatibility

// IMetric is the interface that wraps the logic of processing collected metrics.
type IMetric interface {
	GetName() string
	GetType() string
	GetGauge() (*float64, error)
	GetCounter() (*int64, error)
}
