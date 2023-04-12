package beeep

import (
	"strconv"
	"strings"
	"time"
)

type Action struct {
	Key string
	Act func()
}

type Opt struct {
	title        string
	message      string
	icon         string
	level        Level
	actions      []Action
	duration     time.Duration
	durationName string
	freq         float64
	audio        string
}

type Option func(*Opt)

func AppOption(title string) Option {
	return func(o *Opt) {
		o.title = title
	}
}

func MessageOption(message string) Option {
	return func(o *Opt) {
		o.message = message
	}
}

func IconOption(icon string) Option {
	return func(o *Opt) {
		o.icon = icon
	}
}

func LevelOption(level Level) Option {
	return func(o *Opt) {
		o.level = level
	}
}

func ActionOption(action Action) Option {
	return func(o *Opt) {
		o.actions = append(o.actions, action)
	}
}

func DurationOption(duration time.Duration) Option {
	return func(o *Opt) {
		o.duration = duration
	}
}

// DurationNameOption can only be either short or long for
// windows 10 notifications
func DurationNameOption(duration string) Option {
	return func(o *Opt) {
		o.durationName = duration
	}
}

func FrequencyOption(freq float64) Option {
	return func(o *Opt) {
		o.freq = freq
	}
}

func buildNotifySend(o *Opt) string {
	var result strings.Builder

	if o.title != "" {
		result.WriteString(" -a " + o.title)
	}

	if o.icon != "" {
		result.WriteString(" -i " + o.icon)
	}

	if o.duration != 0 {
		result.WriteString(" -e " + strconv.Itoa(int(o.duration/1000)))
	}

	level := "normal"
	if o.level <= 0 {
		level = "low"
	} else if o.level >= 2 {
		level = "critical"
	}
	result.WriteString(" -u " + level)

	for i, m := range o.actions {
		result.WriteString(" -A=" + strconv.Itoa(i) + "=" + m.Key)
	}

	return result.String() + " " + o.message
}
