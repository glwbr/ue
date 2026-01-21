package trips

import "strings"

type TripStatus int

const (
	StatusUnknown TripStatus = iota
	StatusCompleted
	StatusCanceled
)

var statusStrings = map[TripStatus]string{
	StatusUnknown:   "UNKNOWN",
	StatusCompleted: "COMPLETED",
	StatusCanceled:  "CANCELED",
}

var stringToStatus = map[string]TripStatus{
	"UNKNOWN":   StatusUnknown,
	"COMPLETED": StatusCompleted,
	"CANCELED":  StatusCanceled,
}

func ParseTripStatus(status string) TripStatus {
	normalized := strings.ToUpper(strings.TrimSpace(status))
	if s, ok := stringToStatus[normalized]; ok {
		return s
	}
	return StatusUnknown
}

func (s TripStatus) String() string {
	if str, ok := statusStrings[s]; ok {
		return str
	}
	return statusStrings[StatusUnknown]
}

func (s TripStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *TripStatus) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	*s = ParseTripStatus(str)
	return nil
}
