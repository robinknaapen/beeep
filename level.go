package beeep

// As described in https://specifications.freedesktop.org/notification-spec/notification-spec-latest.html#urgency-levels
type Level byte

// Available levels of urgency
const (
	LevelNormal = iota
	LevelWarning
	LevelCritial
)
