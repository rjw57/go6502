/*
	Package cli provides command line support for go6502.

	It parses CLI flags and exposes the resulting options.
*/
package cli

import (
	"flag"
	"fmt"
	"strings"
)

// Options stores the value of command line options after they're parsed.
type Options struct {
	ViaDumpBinary bool
	ViaDumpAscii  bool
	ViaSsd1306    bool
	Debug         bool
	DebugCmds     commandList
	Speedometer   bool
}

// ParseFlags uses the flag stdlib package to parse CLI options.
func ParseFlags() *Options {
	opt := &Options{}

	// Debug
	flag.BoolVar(&opt.Debug, "debug", false, "Run debugger")
	flag.Var(&opt.DebugCmds, "debug-commands", "Debugger commands to run, semicolon separated.")

	// VIA
	flag.BoolVar(&opt.ViaDumpBinary, "via-dump-binary", false, "VIA6522 dumps binary output")
	flag.BoolVar(&opt.ViaDumpAscii, "via-dump-ascii", false, "VIA6522 dumps ASCII output")
	flag.BoolVar(&opt.ViaSsd1306, "via-ssd1306", false, "SSD1306 OLED display on VIA6522 port B")

	// Speedometer
	flag.BoolVar(&opt.Speedometer, "speedometer", false, "Measure effective clock speed")

	flag.Parse()
	return opt
}

type commandList []string

func (cl *commandList) Set(value string) error {
	list := strings.Split(value, ";")
	for i, value := range list {
		list[i] = strings.TrimSpace(value)
	}
	*cl = list
	return nil
}

func (cl *commandList) String() string {
	return fmt.Sprint(*cl)
}