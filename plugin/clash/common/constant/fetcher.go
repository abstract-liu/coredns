package constant

type FetcherType int

const (
	LOCAL_FILE FetcherType = iota
	REMOTE_HTTP
)

func (ft FetcherType) String() string {
	switch ft {
	case LOCAL_FILE:
		return "LOCAL_FILE"
	case REMOTE_HTTP:
		return "REMOTE_HTTP"
	default:
		return "UNKNOWN"
	}
}
