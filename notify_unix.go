//go:build (linux && !nodbus) || (freebsd && !nodbus) || (netbsd && !nodbus) || (openbsd && !nodbus)
// +build linux,!nodbus freebsd,!nodbus netbsd,!nodbus openbsd,!nodbus

package beeep

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/godbus/dbus/v5"
)

// Notify sends desktop notification.
//
// On Linux it tries to send notification via D-Bus and it will fallback to `notify-send` binary.
func Notify(options ...Option) error {
	opt := &Opt{
		title:   "app",
		message: "Something happened",
		icon:    "",
	}

	for _, o := range options {
		o(opt)
	}

	appIcon := pathAbs(opt.icon)

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

	conn, err := dbus.SessionBus()
	if err != nil {
		return cmd()
	}
	obj := conn.Object("org.freedesktop.Notifications", dbus.ObjectPath("/org/freedesktop/Notifications"))

	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0), appIcon, opt.title, opt.message, []string{}, map[string]dbus.Variant{}, int32(-1))
	if call.Err != nil {
		e := cmd()
		if e != nil {
			e := knotify()
			if e != nil {
				return errors.New("beeep: " + call.Err.Error() + "; " + e.Error())
			}
		}
	}

	return nil
}
