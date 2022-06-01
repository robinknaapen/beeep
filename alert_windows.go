//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js
// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	toast "github.com/go-toast/toast"
)

// Alert displays a desktop notification and plays a default system sound.
func Alert(options ...Option) error {
	if isWindows10 {
		opt := &Opt{
			title:   "app",
			message: "Something happened",
			icon:    "",
		}

		for _, o := range options {
			o(opt)
		}
	
		note := toastNotification(opt)
		note.Audio = toast.Default
		return note.Push()
	}

	if err := Notify(options...); err != nil {
		return err
	}
	return Beep(DefaultFreq, DefaultDuration)
}
