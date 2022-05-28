//go:build linux || freebsd || netbsd || openbsd
// +build linux freebsd netbsd openbsd

package beeep

// Alert displays a desktop notification and plays a beep.
func Alert(options ...Option) error {
	if err := Notify(options...); err != nil {
		return err
	}

	return Beep(options...)
}
