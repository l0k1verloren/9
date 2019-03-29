package cmd

import (
	"fmt"
	"path/filepath"

	"git.parallelcoin.io/dev/9/pkg/util/cl"
)

func Parse(args []string) int {
	// parse commandline
	cmd, tokens, cmds := parseCLI(args)
	if cmd == nil {
		help := commands[HELP]
		cmd = &help
	}
	var datadir string
	// read configuration
	dd, ok := Config["app.datadir"]
	if ok {
		datadir = dd.Value.(string)
		if t, ok := tokens["datadir"]; ok {
			datadir = t.Value
			Config["app.datadir"].Value = datadir
		}
	}
	log <- cl.Debug{"loading config from:", datadir}
	configFile := filepath.Join(datadir, "config")
	fmt.Println("loading config from", configFile)

	fmt.Println(Config)
	// spew.Dump(Config)
	// run as configured
	_ = cmds
	cmd.Handler(
		args,
		tokens,
		cmds,
		commands)
	return 0
}
