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
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/timthreetwelve/bootlogger/utils"
)

// autostartCmd represents the autostart command
var autostartCmd = &cobra.Command{
	Use:     "autostart",
	Aliases: []string{"auto", "start"},
	Example: "  bootlogger autostart --enable\n  bootlogger auto --disable\n  bootlogger autostart -s",
	Short:   "Check status or manage the execution of bootlogger at Windows startup.",
	Long: `The autostart command allow you to enable or disable bootlogger from running at Windows startup,
or check its current status. When enabled, bootlogger will be added to the Windows startup registry key,
HKCU\Software\Microsoft\Windows\CurrentVersion\Run, allowing it to run automatically when Windows starts.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Check if the user has specified --enable or --disable
		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")
		//Don't need to check status flag since it is the default action

		switch {
		case enable:
			msg, err := utils.EnableAutostart()
			if err != nil {
				color.New(color.FgHiRed).Printf("Error enabling autostart: %v\n", err)
				log.Println("bootlogger autostart enable: ", err)
				return
			}
			color.New(color.FgHiWhite).Println(msg)
			log.Println("bootlogger autostart enable: ", msg)
		case disable:
			msg, err := utils.DisableAutostart()
			if err != nil {
				color.New(color.FgHiRed).Printf("Error disabling autostart: %v\n", err)
				log.Println("bootlogger autostart disable: ", err)
				return
			}
			color.New(color.FgHiWhite).Println(msg)
			log.Println("bootlogger autostart disable: ", msg)
		default:
			msg, err := utils.GetAutostartStatus()
			if err != nil {
				color.New(color.FgHiRed).Printf("%s: %v\n", msg, err)
				log.Println("bootlogger autostart status: ", err)
				return
			}
			color.New(color.FgHiWhite).Println(msg)
			log.Println("bootlogger autostart status: ", msg)
			return
		}
	},
}

// init function adds the autostart command to the root command
// and sets up its flags.
func init() {
	rootCmd.AddCommand(autostartCmd)
	autostartCmd.Flags().BoolP("enable", "e", false, "Enable bootlogger to run at startup")
	autostartCmd.Flags().BoolP("disable", "d", false, "Disable bootlogger from running at startup")
	autostartCmd.Flags().BoolP("status", "s", false, "check autostart status")
	autostartCmd.MarkFlagsMutuallyExclusive("enable", "disable", "status")
}
