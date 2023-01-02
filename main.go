package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	disableConcurrency()

	settings, cmdArgs, err := parseArgs()
	if err != nil {
		err = fmt.Errorf("could not resolve flags: %w", err)
		panic(err)
	}

	configPath, err := homedir.Expand(settings.configPath)
	if err != nil {
		err = fmt.Errorf("could not resolve config file: %w", err)
		panic(err)
	}

	type envProfiles map[string]map[string]string
	var config envProfiles
	_, err = toml.DecodeFile(configPath, &config)
	if err != nil {
		err = fmt.Errorf("could not read config file: %w", err)
		panic(err)
	}

	profileEnvVars, ok := config[settings.profile]
	if !ok {
		err = fmt.Errorf("could not find env profile \"%s\" in %s", settings.profile, settings.configPath)
		panic(err)
	}
	environment := os.Environ()
	for key, value := range profileEnvVars {
		line := fmt.Sprintf("%s=%s", key, value)
		environment = append(environment, line)
	}

	if len(cmdArgs) < 1 {
		usage()
		panic("command not given")
	}

	// use higher-level os/exec package to resolve the executable path
	executable, err := exec.LookPath(cmdArgs[0])
	if err != nil {
		err = fmt.Errorf("could not find executable on $PATH: %w", err)
		panic(err)
	}

	if err := syscall.Exec(executable, cmdArgs, environment); err != nil {
		err = fmt.Errorf("could not run %s: %w", executable, err)
		panic(err)
	}
}

func disableConcurrency() {
	// making sure the program remains single-threaded
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func usage() {
	exec := os.Args[0]
	fmt.Printf("Usage of %s:\n", exec)
	fmt.Printf("$ %s -p <profile> <cmd with args...>\n", exec)
	fmt.Printf("runs <cmd with args...> using env vars defined for <profile>\n\n")
	fmt.Printf("Options:\n\n")
	fmt.Printf("  --profile string\n")
	fmt.Printf("  -p string\n")
	fmt.Printf("        Profile, i.e. a set of env vars to use for the command (default \"default\")\n")
	fmt.Printf("  --config string\n")
	fmt.Printf("  -f string\n")
	fmt.Printf("        File to read profiles from (default \"~/.envdo.toml\")\n")
	fmt.Printf(`        Example file:
            # This is an example configuration. Put something like this into ~/.envdo.toml
            
            [default]
            # if no profile is specified, envdo will use the default profile
            FOO = "yes"
            BAR = "for sure"
            
            [other]
            FOO = "correct"
            BAR = "yes"`)
}

// pflag and flag both try to parse into the executed command; we don't want this and roll our own CLI parser
type settings struct {
	profile    string
	configPath string
}

func parseArgs() (flags settings, args []string, err error) {
	err = nil
	flags = settings{
		profile:    "default",
		configPath: "~/.envdo.toml",
	}

	args = make([]string, len(os.Args)-1)
	copy(args, os.Args[1:])

	switch len(args) {
	case 0:
		err = fmt.Errorf("command not given")
		panic(err)
	case 1:
		if strings.HasPrefix(args[0], "-") {
			err = fmt.Errorf("unknown flag %s", args[0])
			panic(err)
		} else {
			return
		}
	}
	// more than just the command is specified; attempt to parse -f/-p and leave the rest for command
	for true {
		if len(args) >= 2 && (args[0] == "-f" || args[0] == "--config") {
			flags.configPath = args[1]
			args = args[2:]
			continue
		} else if len(args) >= 2 && (args[0] == "-p" || args[0] == "--profile") {
			flags.profile = args[1]
			args = args[2:]
			continue
		} else if len(args) > 0 && strings.HasPrefix(args[0], "-") {
			err = fmt.Errorf("unknown flag %s", args[0])
			panic(err)
		} else {
			break
		}
	}

	return
}
