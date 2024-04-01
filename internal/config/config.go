package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func Configure() error {
	for _, cfg := range []struct {
		command string
		prepend []string
		calls   [][]string
	}{
		{
			command: "git",
			prepend: []string{"config", "--global"},
			calls: [][]string{
				{"delta.line-numbers-minus-style", "red"},
				{"delta.line-numbers-plus-style", "green"},
				{"delta.minus-style", "syntax 52"},
				{"delta.minus-emph-style", "syntax 88"},
				{"delta.plus-style", "syntax 22"},
				{"delta.plus-emph-style", "syntax 28"},
				{"tig.color.cursor", "15 235"},
				{"tig.color.title-focus", "white 234 bold"},
				{"tig.color.title-blur", "white 234 dim"},
			},
		},
	} {
		exe, err := exec.LookPath(cfg.command)
		if err != nil && errors.Is(err, exec.ErrNotFound) {
			fmt.Println(cfg.command, "not found, skipping relative configurations")
			continue
		}
		for _, call := range cfg.calls {
			call = append(cfg.prepend, call...)
			fmt.Println(exe, call)
			cmd := exec.Command(exe, call...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}
	// TODO: file based configs (vim/neovim colorscheme)
	return nil
}

func exists() {
}

// func cmd(name string, args ...string) {
// 	cmd := exec.Command(name, args...)
// 	cmd.Run()
// }
