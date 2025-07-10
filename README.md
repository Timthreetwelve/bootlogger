<p align="center">
  <a target="_blank" rel="noopener noreferrer">
    <img width="128" src="winres/bl128.png" alt="bootlogger Logo">
  </a>
</p>
<h1 align="center">
  bootlogger
</h1>

<div align="center">

![GitHub License](https://img.shields.io/badge/license-MIT-green?style=plastic)
![Static Badge](https://img.shields.io/badge/built%20for-Windows-2196F3?style=plastic)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/Timthreetwelve/bootlogger?style=plastic)](https://github.com/Timthreetwelve/bootlogger/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/timthreetwelve/bootlogger?style=plastic&color=orange)](https://github.com/Timthreetwelve/bootlogger/releases/latest)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/timthreetwelve/bootlogger/latest?style=plastic)](https://github.com/Timthreetwelve/bootlogger/commits/main)
[![GitHub last commit](https://img.shields.io/github/last-commit/timthreetwelve/bootlogger?style=plastic)](https://github.com/Timthreetwelve/bootlogger/commits/main)

</div>

## bootlogger
bootlogger is a command-line tool for Windows that logs system restart events, commonly known as a [reboot](https://en.wikipedia.org/wiki/Reboot). It records the date and time of each reboot, along with the computer name, and optionally the Windows version and build info. The log file name and location can be specified and the boot time can be formatted in any of more than a dozen date/time formats. bootlogger can also add itself to the Windows registry to run automatically at startup, ensuring that all reboots are logged. BTW, the name is "bootlogger", all lowercase.

## Why bootlogger and Why Go?
Perhaps you've seen my [Windows Update Viewer](https://github.com/Timthreetwelve/WUView) app or maybe [Get My IP](https://github.com/Timthreetwelve/GetMyIP) or one of the others. I needed to take a short break from maintaining those and felt like trying something new. I'd read good things about the [Go language](https://go.dev) and decided to give it a go _(pun intended)_. I was using a PowerShell script to log reboots, so I decided that rewriting it in Go would be a good first project. The result is what you see here. I found it easy to get started with Go. The [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/cobra) packages made things a lot easier

## What it logs
#### Computer name
bootlogger tries to get the computer name from the host name as reported by the kernel. If that fails, it tries to get the name from the `COMPUTERNAME` environment variable. If that also fails, it will use "unknown" for the computer name.

#### Reboot date and time
bootlogger uses the `GetTickCount64` function to get the number of milliseconds that have elapsed since the system was started. That value is subtracted from the current date and time, giving the date and time of the last system start. That result is then formatted according to the `timeformat` configuration option discussed below.

#### Windows version and build information _(Optional)_
bootlogger uses the registry to determine the Windows version and build information. Specifically, information found in `HKLM\SOFTWARE\Microsoft\Windows NT\Current version\`.

 The version comes from the value of `ProductName`. Microsoft chose not to update this name from "Windows 10" when it released Windows 11. Therefore, if the Windows build number is equal to or greater than 22000, Windows 10 is changed to Windows 11 to accurately report the version of Windows.

The build number comes from combining the `CurrentBuildNumber` value with the `UBR` value giving the build number in the familiar xxxxx.yyyy format.

The computer name, reboot date and time, and the Windows version and build are written as a single record to the log file. As noted above, the logging of the version and build info is optional. That is controlled by the `no-buildinfo` option discussed below.

## Usage
Running bootlogger without any additional commands or flags will write a record to the log file using the default configuration or the updated configuration found in environment variables, a configuration file, and/or command-line flags. _Continue reading for examples_.

## Commands
Before discussing configuration, a few words about commands. Commands can be differentiated from flags by the fact that they don't begin with the dash, or two dashes, that flags use.

bootlogger has commands to perform some utility functions without writing a record to the log file. The first of these commands is `help`. Running `bootlogger.exe help` will show help for bootlogger.

`bootlogger.exe version` will show the version number.

`bootlogger.exe printconfig` will show the current configuration.

`bootlogger.exe printlog` will print the log file. Output can be piped through `more` for easier reading.

 A tool that logs reboots wouldn't help much if you had to _remember_ to run it after each reboot. `bootlogger.exe autostart` will allow you to enable or disable bootlogger from running at Windows startup,
or check its current status. When enabled, bootlogger will be added to the Windows startup registry key,
`HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, allowing it to run automatically when Windows starts. If this doesn't work for you, it can be added to Task Scheduler or the to `C:\Users\<user>\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\`.

Use `bootlogger.exe [command] --help` for more information about any of these commands.

#### Example command output:
```shell
bootlogger.exe printconfig
bootlogger configuration:
  logfile         = D:\Logs\bootlogger.test.log
  namewidth       = 12
  no-buildinfo    = true
  no-text         = false
  quiet           = false
  timeformat      = RFC822

bootlogger.exe autostart --enable
bootlogger was added to 'HKCU\Software\Microsoft\Windows\CurrentVersion\Run' and will run from: C:\PortableApps\bootlogger.exe when Windows starts.

bootlogger.exe version
bootlogger version: 0.1.0
```

## Configuration
bootlogger can be configured through command-line flags, environment variables, or a configuration file.

Configurable options are:

|Option|Type|Default value|Description|
|------|----|-------------|-----------|
|`dryrun`|Boolean|`false`|Write log line to console but _**not**_ to the log file. Useful for testing.|
|`logfile`|String|`bootlog.txt`|Log file filename with fully qualified ([absolute](https://www.computerhope.com/issues/ch001708.htm)) path. Use quotes if there are spaces in the path.|
|`namewidth`|Integer|`14`|Minimum computer name width. Longer names are _**not**_ truncated. Shorter names are padded with spaces.|
|`no-buildinfo`|Boolean|`false`|Version and build info will _**not**_ be included in the log entry.|
|`no-text`|Boolean|`false`|The "was rebooted on" text will _**not**_ be included in the log entry.|
|`quiet`|Boolean|`false`|Quiet operation. Do _**not**_ print non-error messages to the console.|
|`timeformat`|String|`12Hour`|Date/Time format. Use `12Hour` for 12 hour, `24Hour` for 24 hour or most of the Go pre-defined formats. _DateOnly and TimeOnly are _**not**_ included_. See https://pkg.go.dev/time#pkg-constants for details. |


### Configuration file
If used, the configuration file _**must**_ be named `config.yaml` and _**must**_ be in the same folder as bootlogger.exe. The configuration file _**must**_ be in YAML format. YAML was chosen for its simple syntax and readability. If you aren't familiar with YAML, this [Wikipedia article](https://en.wikipedia.org/wiki/YAML) may help.

#### Example configuration file:
```yaml
# Name of the log file
LogFile: "D:\Logs\bootlogger.test.log"

# Minimum width of the computer name field in the log
NameWidth: 12

# Time format. Either the custom 12-hour format or any of the Go pre-defined date/time formats.
TimeFormat: 24Hour
```

### Environment variables
All environment variables _**must**_ be prefixed with `BOOTLOG_` followed by the configuration option. _That's a single underscore character between BOOTLOG and the option._

#### Example environment variables:
```shell
setx BOOTLOG_TIMEFORMAT 12Hour

setx BOOTLOG_No-Buildinfo true
```

### Command-line flags
To configure bootlogger with command-line flags, prefix the configuration option with two consecutive dashes. So to configure the date/time format to use 24 hour time use `bootlogger.exe -timeformat 24hour`. Additional examples are found below.

Values for a flag can be separated by a space, separated by an equals sign, or follow the flag without a space.

Some command-line flags have a shorthand flag consisting of a single character. To use the shorthand flag, use only one dash before the character.

|flag|shorthand|
|-|-|
|`--help`|`-h`|
|`--logfile`|`-l`|
|`--namewidth`|`-w`|
|`--quiet `|`-q`|
|`--timeformat`|`-t`|

Note that `--dryrun, --no-text` and `--no-buildinfo` do not have shorthand flags.

Note that shorthand flags must be **lower case**.

#### Command-line examples:
```shell
bootlogger.exe --timeformat 12hour
Desktop-Computer was rebooted on 2025-06-19  7:43:19 AM [Windows 11 Pro build 26100.4061]

bootlogger.exe -t=24h
Desktop-Computer was rebooted on 2025-06-19 07:43:19 [Windows 11 Pro build 26100.4061]

bootlogger.exe --timeformat RFC822 --no-buildinfo
Desktop-Computer was rebooted on 19 Jun 25 07:43 CDT

bootlogger.exe --no-text
Desktop-Computer 2025-06-19  7:43:19 AM [Windows 11 Pro build 26100.4061]

bootlogger.exe --logfile d:\logs\bootlog.txt -q
#The entry will be made in the log file but nothing will be written to the console.

bootlogger.exe --dryrun
Dry run: the following line would be written to ./bootlogger.log
Desktop-Computer was rebooted on 2025-06-25 07:18:11 [Windows 11 Pro build 26100.4351]
```
## Installation
bootlogger is a "portable" app and as such doesn't require installation. Just copy the bootlogger.exe file to a folder where you have read and write permissions. Write permission is needed for the configuration file and for the log file itself.

## Acknowledgements
The following packages were used in the development of bootlogger:

- cobra - https://github.com/spf13/cobra
- color - https://github.com/fatih/color
- pflag - https://github.com/spf13/pflag
- viper - https://github.com/spf13/viper
