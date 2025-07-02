/*
Copyright © 2025 Tim Kennedy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

// getBootTime retrieves the last boot time of the system.
// It uses the GetTickCount64 function from kernel32.dll to get the uptime
// and calculates the boot time by subtracting the uptime from the current time.
func getBootTime() time.Time {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 := kernel32.NewProc("GetTickCount64")

	ticks, _, _ := getTickCount64.Call()
	uptime := time.Duration(ticks) * time.Millisecond
	bootTime := time.Now().Add(-uptime)
	return bootTime
}

// GetBootTimeFormatted returns the boot date and time formatted as strings.
// The time preference is determined by the "timeformat" configuration setting.
// Supported formats include "12", "24", and standard Go time formats.
// See https://pkg.go.dev/time#pkg-constants for details.
func GetBootTimeFormatted() string {
	bootTime := getBootTime()

	switch tf := strings.ToUpper(viper.GetString("timeformat")); tf {

	case "12", "12HOUR", "12H":
		// Custom format for 12-hour time with AM/PM.
		// Hour will be padded with a space for values less than 10.
		dateStr := bootTime.Format("2006-01-02")
		timeStr := bootTime.Format("3:04:05 PM")
		return fmt.Sprintf("%s %11s", dateStr, timeStr)
	case "24", "24HOUR", "24H", "DATETIME":
		return bootTime.Format(time.DateTime)
	case "LAYOUT":
		return bootTime.Format(time.Layout)
	case "ANSIC":
		return bootTime.Format(time.ANSIC)
	case "UNIXDATE":
		return bootTime.Format(time.UnixDate)
	case "RUBYDATE":
		return bootTime.Format(time.RubyDate)
	case "RFC822":
		return bootTime.Format(time.RFC822)
	case "RFC822Z":
		return bootTime.Format(time.RFC822Z)
	case "RFC850":
		return bootTime.Format(time.RFC850)
	case "RFC1123":
		return bootTime.Format(time.RFC1123)
	case "RFC1123Z":
		return bootTime.Format(time.RFC1123Z)
	case "RFC3339":
		return bootTime.Format(time.RFC3339)
	case "RFC3339NANO":
		return bootTime.Format(time.RFC3339Nano)
	case "KITCHEN":
		return bootTime.Format(time.Kitchen)
	case "STAMP":
		return bootTime.Format(time.Stamp)
	case "STAMPMILLI":
		return bootTime.Format(time.StampMilli)
	case "STAMPMICRO":
		return bootTime.Format(time.StampMicro)
	case "STAMPNANO":
		return bootTime.Format(time.StampNano)
	default:
		fmt.Printf("Invalid time format specified: %s. Using default string format.\n", viper.GetString("timeformat"))
		return bootTime.Format(time.DateTime)
	}
}
