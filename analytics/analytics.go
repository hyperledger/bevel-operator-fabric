package analytics

type Analytics interface {
	TrackEvent(event string, properties map[string]interface{}) error
}

type amplitudeAnalytics struct {
}

func NewAmplitudeAnalytics(apiKey string) (*amplitudeAnalytics, error) {
	return &amplitudeAnalytics{}, nil
}

type dummyAnalytics struct{}

func NewDummyAnalytics() *dummyAnalytics {
	return &dummyAnalytics{}
}
