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
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// PrintConfig prints the current configuration settings.
func PrintConfig() {
	settings := viper.AllSettings()

	color.HiWhite("bootlogger configuration:")
	cyan := color.New(color.FgCyan).PrintfFunc()
	cyanHi := color.New(color.FgHiCyan).PrintfFunc()

	// Collect and sort keys
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Print settings in sorted order
	for _, key := range keys {
		value := settings[key]
		cyan("  %-15s", key)
		fmt.Printf(" = ")
		cyanHi("%v\n", value)
	}

	log.Println("bootlogger configuration: ")
	for _, key := range keys {
		log.Printf("  %-15s  = %v", key, settings[key])
	}
}
