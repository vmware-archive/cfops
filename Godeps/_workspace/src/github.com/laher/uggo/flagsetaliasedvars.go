//extend FlagSet with support for flag aliasMap
package uggo

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)


// This FlagSet embeds Golang's FlagSet but adds some extra flavour 
// - support for 'aliased' arguments
// - support for 'gnuification' of short-form arguments
// - handling of --help and --version
// - some flexibility in 'usage' message
type FlagSetWithAliases struct {
	*flag.FlagSet
	AutoGnuify   bool //sets whether it interprets options gnuishly (e.g. -lah being equivalent to -l -a -h)
	name           string
	argUsage       string
	out            io.Writer
	aliasMap       map[string][]string
	isPrintUsage   *bool //optional usage (you can use your own, or none, instead)
	isPrintVersion *bool //optional (you can use your own, or none, instead)
	version        string
}


// factory which sets defaults and 
// NOTE: discards thet output of the embedded flag.FlagSet. This was necessary in order to override the 'usage' message
func NewFlagSet(desc string, errorHandling flag.ErrorHandling) FlagSetWithAliases {
	fs := flag.NewFlagSet(desc, errorHandling)
	fs.SetOutput(ioutil.Discard)
	return FlagSetWithAliases{fs, false, desc, "", os.Stderr, map[string][]string{}, nil, nil, "unknown"}
}

// Factory setting useful defaults
// Sets up --help and --version flags
func NewFlagSetDefault(name, argUsage, version string) FlagSetWithAliases {
	fs := flag.NewFlagSet(name+" "+argUsage, flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	// temp variables for storing defaults
	tmpPrintUsage := false
	tmpPrintVersion := false
	flagSet := FlagSetWithAliases{fs, true, name, argUsage, os.Stderr, map[string][]string{}, &tmpPrintUsage, &tmpPrintVersion, version}
	flagSet.BoolVar(flagSet.isPrintUsage, "help", false, "Show this help")
	flagSet.BoolVar(flagSet.isPrintVersion, "version", false, "Show version")
	flagSet.version = version
	return flagSet
}

// process built-in help and version flags.
//  Returns 'true' when one of these was set. (i.e. stop processing)
func (flagSet FlagSetWithAliases) ProcessHelpOrVersion() bool {
	if flagSet.IsHelp() {
		flagSet.Usage()
		return true
	} else if flagSet.IsVersion() {
		flagSet.PrintVersion()
		return true
	}
	return false
}

//convenience method for storing 'usage' behaviour
func (flagSet FlagSetWithAliases) IsHelp() bool {
	return *flagSet.isPrintUsage
}

//convenience method for storing 'get version' behaviour
func (flagSet FlagSetWithAliases) IsVersion() bool {
	return *flagSet.isPrintVersion
}

// PrintVersion
func (flagSet FlagSetWithAliases) PrintVersion() {
	fmt.Fprintf(flagSet.out, "`%s` version: '%s'\n", flagSet.name, flagSet.version)
}

// Print Usage message
func (flagSet FlagSetWithAliases) Usage() {
	fmt.Fprintf(flagSet.out, "Usage: `%s %s`\n", flagSet.name, flagSet.argUsage)
	flagSet.PrintDefaults()
}

// Set writer for displaying help and usage messages.
func (flagSet FlagSetWithAliases) SetOutput(out io.Writer) {
	flagSet.out = out
}

// Set up multiple names for a bool flag
func (flagSet FlagSetWithAliases) AliasedBoolVar(p *bool, items []string, def bool, description string) {
	flagSet.RecordAliases(items, "bool")
	for _, item := range items {
		flagSet.BoolVar(p, item, def, description)
	}
}

// Set up multiple names for a time.Duration flag
func (flagSet FlagSetWithAliases) AliasedDurationVar(p *time.Duration, items []string, def time.Duration, description string) {
	flagSet.RecordAliases(items, "duration")
	for _, item := range items {
		flagSet.DurationVar(p, item, def, description)
	}
}

// Set up multiple names for a float64 flag
func (flagSet FlagSetWithAliases) AliasedFloat64Var(p *float64, items []string, def float64, description string) {
	flagSet.RecordAliases(items, "float64")
	for _, item := range items {
		flagSet.Float64Var(p, item, def, description)
	}
}

// Set up multiple names for an int flag
func (flagSet FlagSetWithAliases) AliasedIntVar(p *int, items []string, def int, description string) {
	flagSet.RecordAliases(items, "int")
	for _, item := range items {
		flagSet.IntVar(p, item, def, description)
	}
}

// Set up multiple names for an int64 flag
func (flagSet FlagSetWithAliases) AliasedInt64Var(p *int64, items []string, def int64, description string) {
	flagSet.RecordAliases(items, "int64")
	for _, item := range items {
		flagSet.Int64Var(p, item, def, description)
	}
}

// Set up multiple names for a string flag
func (flagSet FlagSetWithAliases) AliasedStringVar(p *string, items []string, def string, description string) {
	flagSet.RecordAliases(items, "string")
	for _, item := range items {
		flagSet.StringVar(p, item, def, description)
	}
}

// returns true if the given flag name is the 'main' name or a subsequent name
func (flagSet FlagSetWithAliases) isAlternative(name string) bool {
	for _, altSlice := range flagSet.aliasMap {
		for _, alt := range altSlice {
			if alt == name {
				return true
			}
		}
	}
	return false
}

// keep track of aliases to a given flag
func (flagSet FlagSetWithAliases) RecordAliases(items []string, typ string) {
	var key string
	for i, item := range items {
		if i == 0 {
			key = item
			if _, ok := flagSet.aliasMap[key]; !ok {
				flagSet.aliasMap[key] = []string{}
			}
		} else {
			//key is same as before
			flagSet.aliasMap[key] = append(flagSet.aliasMap[key], item)
		}
	}
}

// parse flags from a given set of argv type flags
func (flagSet FlagSetWithAliases) Parse(call []string) error {
	if flagSet.AutoGnuify {
		call = Gnuify(call)
	}
	return flagSet.FlagSet.Parse(call)
}


func (flagSet FlagSetWithAliases) ParsePlus(call []string) (error, int) {
	err := flagSet.Parse(call)
	if err != nil {
		fmt.Fprintf(flagSet.out, "Flag error: %v\n\n", err.Error())
		flagSet.Usage()
		return err, 0
	}
	if flagSet.ProcessHelpOrVersion() {
		return EXIT_OK, 0
	}
	return nil, 0
}

// print defaults to the default output writer
func (flagSet FlagSetWithAliases) PrintDefaults() {
	flagSet.PrintDefaultsTo(flagSet.out)
}

// print defaults to a given writer. 
// Output distinguishes aliases
func (flagSet FlagSetWithAliases) PrintDefaultsTo(out io.Writer) {
	flagSet.FlagSet.VisitAll(func(fl *flag.Flag) {
		l := 0
		alts, isAliased := flagSet.aliasMap[fl.Name]
		if isAliased {
			li, _ := fmt.Fprintf(out, "  ")
			l += li
			if len(fl.Name) > 1 {
				li, _ := fmt.Fprint(out, "-")
				l += li
			}
			li, _ = fmt.Fprintf(out, "-%s", fl.Name)
			l += li
			for _, alt := range alts {
				fmt.Fprint(out, " ")
				l += 1
				if len(alt) > 1 {
					li, _ := fmt.Fprint(out, "-")
					l += li
				}
				li, _ := fmt.Fprintf(out, "-%s", alt)
				l += li
			}
			//defaults:
			//no known straightforward way to test for boolean types
			if fl.DefValue == "false" {
			} else {
				li, _ = fmt.Fprintf(out, "=%s", fl.DefValue)
				l += li
			}
			fmt.Fprint(out, " ")
			l += 1
		} else if !flagSet.isAlternative(fl.Name) {
			li, _ := fmt.Fprint(out, "  ")
			l += li
			if len(fl.Name) > 1 {
				li, _ := fmt.Fprint(out, "-")
				l += li
			}
			if fl.DefValue == "false" {
				li, _ = fmt.Fprintf(out, "-%s", fl.Name)
				l += li
			} else {
				format := "-%s=%s"
				li, _ = fmt.Fprintf(out, format, fl.Name, fl.DefValue)
				l += li
			}
		} else {
			//fmt.Fprintf(out, "alias %s\n", fl.Name)
		}
		if !flagSet.isAlternative(fl.Name) {
			for l < 25 {
				l += 1
				fmt.Fprintf(out, " ")
			}
			fmt.Fprintf(out, ": %s\n", fl.Usage)
		}

	})
	fmt.Fprintln(out, "")
}

// function which can (open and) return a File at some later time
type FileOpener func() (*os.File, error)

// converts arguments to readable file references.
// An argument with filename "-" is treated as the 'standard input'
func (flagSet FlagSetWithAliases) ArgsAsReadables() []FileOpener {
	args := flagSet.Args()
	if len(args) > 0 {
		readers := []FileOpener{}
		for _, arg := range args {
			if arg == "-" {
				reader := func() (*os.File, error) {
					return os.Stdin, nil
				}
				readers = append(readers, reader)
			} else {
				reader := func() (*os.File, error) {
					return os.Open(arg)
				}
				readers = append(readers, reader)
			}
		} 
		return readers
	} else {
		reader := func() (*os.File, error) {
			return os.Stdin, nil
		}
		return []FileOpener{reader}
	}
}

// atomically open all files at once. 
// Only use this when you actually want all open at once (rather than sequentially)
// e.g. writing the same data to all at once as-in a 'tee' operation.
func OpenAll(openers []FileOpener) ([]*os.File, error) {
	files := []*os.File{}
	for _, opener := range openers {
		file, err := opener()
		if err != nil {
			//close all opened files
			for _, openedfile := range files {
				openedfile.Close()
			}
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// Convert arguments to File openers.
// An argument with filename "-" is treated as the 'standard output'
func ToWriteableOpeners(args []string, flag int, perm os.FileMode) []FileOpener {
	return ToPipeWriteableOpeners(args, flag, perm, os.Stdout)
}

// Convert arguments to File openers.
// An argument with filename "-" is treated as the 'standard output'
// Takes a writer 'outPipe' for handling this special case
func ToPipeWriteableOpeners(args []string, flag int, perm os.FileMode, outPipe *os.File) []FileOpener {
	if len(args) > 0 {
		writers := []FileOpener{}
		for _, arg := range args {
			if arg == "-" {
				writer := func() (*os.File, error) {
					return outPipe, nil
				}
				writers = append(writers, writer)
			} else {
				writer := func() (*os.File, error) {
					return os.OpenFile(arg, os.O_WRONLY|flag, perm)
				}
				writers = append(writers, writer)
			}
		} 
		return writers
	} else {
		writer := func() (*os.File, error) {
			return os.Stdout, nil
		}
		return []FileOpener{writer}
	}
}

// Convert arguments to File openers.
// An argument with filename "-" is treated as the 'standard output'
func (flagSet FlagSetWithAliases) ArgsAsWriteables(flag int, perm os.FileMode) []FileOpener {
	args := flagSet.Args()
	return ToPipeWriteableOpeners(args, flag, perm, os.Stdout)
}

// Convert arguments to File openers.
// An argument with filename "-" is treated as the 'standard output'
// Takes a writer 'outPipe' for handling this special case
func (flagSet FlagSetWithAliases) ArgsAsPipeWriteables(flag int, perm os.FileMode, outPipe *os.File) []FileOpener {
	args := flagSet.Args()
	return ToPipeWriteableOpeners(args, flag, perm, outPipe)
}

