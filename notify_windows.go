//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js
// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"syscall"
	"time"

	toast "github.com/go-toast/toast"
	"github.com/tadvi/systray"
	"golang.org/x/sys/windows/registry"
)

var isWindows10 bool
var applicationID string

func init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return
	}

	isWindows10 = maj == 10

	if isWindows10 {
		applicationID = appID()
	}
}

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
	
	if isWindows10 {
		return toastNotify(opt)
	}

	err := baloonNotify(opt.title, opt.message, opt.icon, false)
	if err != nil {
		e := msgNotify(opt.title, opt.message)
		if e != nil {
			return errors.New("beeep: " + err.Error() + "; " + e.Error())
		}
	}

	return nil

}

func msgNotify(title, message string) error {
	msg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}
	cmd := exec.Command(msg, "*", "/TIME:3", title+"\n\n"+message)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func baloonNotify(title, message, appIcon string, bigIcon bool) error {
	tray, err := systray.New()
	if err != nil {
		return err
	}

	err = tray.ShowCustom(pathAbs(appIcon), title)
	if err != nil {
		return err
	}

	go func() {
		tray.Run()
		time.Sleep(3 * time.Second)
		tray.Stop()
	}()

	return tray.ShowMessage(title, message, bigIcon)
}

func toastNotify(opt *Opt) error {
	notification := toastNotification(opt)
	return notification.Push()
}

func toastNotification(opt *Opt) toast.Notification {
	t := toast.Notification{
		AppID:   applicationID,
		Title:   opt.title,
		Message: opt.message,
		Icon:    pathAbs(opt.icon),
	}
	
	if opt.durationName != "" {
		result, err := toast.Duration(opt.durationName)
		t.Duration = result
		if err != nil {
			t.Duration, _ = toast.Duration("short")
		}
	}
	
	return t
}

func appID() string {
	defID := "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command("powershell", "-NoProfile", "Get-StartApps")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return defID
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "powershell.exe") {
			sp := strings.Split(line, " ")
			if len(sp) > 0 {
				return sp[len(sp)-1]
			}
		}
	}

	return defID
}
