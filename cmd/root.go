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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/timthreetwelve/bootlogger/utils"
)

// Constants
const (
	configFileName     = "config.yaml"
	defaultDryRun      = false
	defaultFileName    = "bootlog.txt"
	defaultNameWidth   = 14
	defaultNoBuildInfo = false
	defaultNoText      = false
	defaultQuiet       = false
	defaultTimeFormat  = "12Hour"
	envPrefix          = "BOOTLOG"
)

// Version is set during the build process using the -ldflags option.
// If not set, it defaults to "unknown".
var Ver string = "unknown"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bootlogger",
	Short: "A Windows command line tool that records system restart events, aka reboots, in a configurable log file.",
	Long: `bootlogger is a command line tool for Windows that logs system restart events, also known as reboots.
It records the date and time of each reboot, along with the computer name, and optionally the Windows version and build info.
The log file name and location can be specified, and the boot time can be formatted in any of more than a dozen formats.
bootlogger can also add itself to the Windows registry to run automatically at startup, ensuring that all reboots are logged.

See https://github.com/Timthreetwelve/bootlogger for more information and documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.WriteLog()
		if err != nil {
			fmt.Printf("bootlogger error: %v\n", err)
			log.Printf("bootlogger ended with an error: %v", err)
			log.Println("")
			os.Exit(1)
		}
	},

	// Turn off the default completion command
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
		HiddenDefaultCmd:  true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("bootlogger error: %v\n", err)
		log.Printf("bootlogger ended with an error: %v", err)
		log.Println("")
		os.Exit(1)
	}
	log.Println("bootlogger executed successfully.")
	log.Println("")
}

// The init function is a special function used for initialization purposes.
// It is automatically executed when a package is initialized, before
// the execution of the main function or any other code in the package.
func init() {
	// Initialize application logging before anything else so that
	// any subsequent errors can be logged.
	InitLogging()

	// Start cobra and viper configuration
	// Allow case insensitive command names
	cobra.EnableCaseInsensitive = true

	// Disable Mousetrap
	cobra.MousetrapHelpText = ""

	// Initialize the configuration settings
	cobra.OnInitialize(initConfig)

	// Add a persistent flag for the timeformat option
	rootCmd.LocalFlags().StringP("timeformat", "t", defaultTimeFormat,
		"Time format, 12 or 24 hour, or most pre-defined Go time formats\nSee https://pkg.go.dev/time#pkg-constants for details")
	if err := viper.BindPFlag("timeformat", rootCmd.PersistentFlags().Lookup("timeformat")); err != nil {
		log.Printf("Error binding flag 'timeformat': %v", err)
	}

	// Add a persistent flag for the no-buildinfo option
	rootCmd.LocalFlags().Bool("no-buildinfo", defaultNoBuildInfo, "Do not include version and build information in the log")
	if err := viper.BindPFlag("no-buildinfo", rootCmd.PersistentFlags().Lookup("no-buildinfo")); err != nil {
		log.Printf("Error binding flag 'no-buildinfo': %v", err)
	}

	// Add a persistent flag for the no-text option
	rootCmd.LocalFlags().Bool("no-text", defaultNoBuildInfo, "Do not include the 'was rebooted on' text in the log")
	if err := viper.BindPFlag("no-text", rootCmd.PersistentFlags().Lookup("no-text")); err != nil {
		log.Printf("Error binding flag 'no-text': %v", err)
	}

	// Add a persistent flag for the logFile option
	rootCmd.LocalFlags().StringP("logfile", "l", defaultFileName, "Path to the log file, including the file name and extension")
	if err := viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile")); err != nil {
		log.Printf("Error binding flag 'logfile': %v", err)
	}

	// Add a persistent flag for the nameField option
	rootCmd.LocalFlags().IntP("namewidth", "w", defaultNameWidth, "Minimum width of the computer name field in the log")
	if err := viper.BindPFlag("namewidth", rootCmd.PersistentFlags().Lookup("namewidth")); err != nil {
		log.Printf("Error binding flag 'namewidth': %v", err)
	}

	// Add a persistent flag for the quiet option
	rootCmd.LocalFlags().BoolP("quiet", "q", defaultQuiet, "Do not print non-error messages to the console")
	if err := viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet")); err != nil {
		log.Printf("Error binding flag 'quiet': %v", err)
	}

	// Add a persistent flag for the quiet option
	rootCmd.LocalFlags().Bool("dryrun", defaultQuiet, "Print log entry to the console but do not write to the log file")
	if err := viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dryrun")); err != nil {
		log.Printf("Error binding flag 'dryrun': %v", err)
	}

	// Normalize flag names to lower case
	// This ensures that flags are case-insensitive and can be used in any case
	// For example, --LogFile, --logfile, and --LOGFILE will all be treated as the same flag
	// Note that this DOES NOT normalize the shorthand flags
	rootCmd.SetGlobalNormalizationFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// path to config file
	configPath := filepath.Join(utils.GetExecutableFolder(), configFileName)

	// Set the file name and path of the configuration file
	// A config file is not required, but if it exists, it will be used
	viper.SetConfigFile(configPath)

	// Set default values for the configuration
	viper.SetDefault("logfile", filepath.Join(execFolder, defaultFileName))
	viper.SetDefault("namewidth", defaultNameWidth)
	viper.SetDefault("no-buildinfo", defaultNoBuildInfo)
	viper.SetDefault("no-text", defaultNoText)
	viper.SetDefault("quiet", defaultQuiet)
	viper.SetDefault("timeformat", defaultTimeFormat)

	// Use env variables and set the prefix for viper
	// This allows viper to read environment variables with the prefix BOOTLOG_
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)

	// Read the config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			// If the config file does not exist, we can ignore the error
			log.Printf("bootlogger config file not found: %s", err)
			log.Println("bootlogger will use environment variables and/or command-line flags for configuration.")
			return
		} else if _, ok := err.(viper.ConfigParseError); ok {
			log.Printf("BootLogger error reading config file: %v", err)
			return
		}
	}

	logFile := viper.GetString("logfile")
	if !CheckFullyQualifiedPath(logFile) {
		fmt.Printf("Please use a fully qualified (absolute) path for the log file. '%s' is not fully qualified.\n", logFile)
		log.Printf("bootlogger log file name is not absolute. %s", logFile)
		os.Exit(1)
	}
}

// Initialize logging and write version and location to the log
func InitLogging() {
	appLogFile, err := utils.GetAppLogFile()
	if err != nil {
		log.Fatalf("bootLogger: error creating application log file: %v", err)
	}
	file, _ := os.OpenFile(appLogFile.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	log.SetOutput(file)

	execPath, err := utils.GetExecutablePath()
	if err != nil {
		log.Fatalf("bootLogger error: %v", err)
	}
	log.Printf("bootlogger version %s is starting from %s.", Ver, execPath)
}

// CheckFullyQualifiedPath ensures that the path id absolute aka fully qualified
func CheckFullyQualifiedPath(file string) bool {
	if filepath.IsAbs(file) {
		return true
	}
	return false
}
