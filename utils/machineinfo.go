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
	"log"
	"os"

	"github.com/spf13/viper"
)

// GetComputerName retrieves the computer name from the Hostname function.
// If the Hostname function fails, it checks for the environment variable COMPUTERNAME.
// If the variable is not set, it returns "Unknown".
func GetComputerName() string {
	comp, err := os.Hostname()
	if err == nil && comp != "" {
		return comp
	}
	comp = os.Getenv("COMPUTERNAME")
	if comp != "" {
		return comp
	}
	log.Println("bootlogger could not determine computer name. using 'Unknown'")
	return "Unknown"
}

// GetFormattedComputerName returns the computer name formatted to the width specified in the configuration.
// The width can be used to maintain a consistent format of the log file if more than one computer is logged.
func GetFormattedComputerName() string {
	comp := GetComputerName()
	width := viper.GetInt("nameWidth")
	return fmt.Sprintf("%-*s", width, comp)
}
