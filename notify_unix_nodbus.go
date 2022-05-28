//go:build (linux && nodbus) || (freebsd && nodbus) || (netbsd && nodbus) || (openbsd && nodbus)
// +build linux,nodbus freebsd,nodbus netbsd,nodbus openbsd,nodbus

package beeep

import (
	"errors"
	"os/exec"
)

// Notify sends desktop notification.
func Notify(options ...Option) error {
	opt := &Opt{
		title:   "app",
		message: "Something happened",
		icon:    "",
	}

	for _, o := range options {
		o(opt)
	}

	appIcon = pathAbs(opt.icon)

	cmd := func() error {
		send, err := exec.LookPath("sw-notify-send")
		if err != nil {
			send, err = exec.LookPath("notify-send")
			if err != nil {
				return err
			}
		}

		args := strings.Split(buildNotifySend(opt), " ")
		c := exec.Command(send, args...)
		return c.Run()
	}

	knotify := func() error {
		send, err := exec.LookPath("kdialog")
		if err != nil {
			return err
		}
		c := exec.Command(send, "--title", opt.title, "--passivepopup", opt.message, "10", "--icon", appIcon)
		return c.Run()
	}

	err := cmd()
	if err != nil {
		e := knotify()
		if e != nil {
			return errors.New("beeep: " + err.Error() + "; " + e.Error())
		}
	}

	return nil
}
