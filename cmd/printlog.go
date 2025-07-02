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
package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/timthreetwelve/bootlogger/utils"
)

// printlogCmd represents the printlog command
var printlogCmd = &cobra.Command{
	Use:     "printlog",
	Short:   "Prints the bootlogger log file to the console. Output can be piped through 'more' for easier reading.",
	Aliases: []string{"print-log", "pl"},
	Example: "  bootlogger printlog\n  bootlogger pl | more",
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.PrintLog()
		if err != nil {
			color.New(color.FgHiRed).Printf("Error reading the log file: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(printlogCmd)
	// No flags are needed for this command, it just prints the log file
	// Output may be large, but can be piped through 'more'
}
