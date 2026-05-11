package cmd

import (
	"github.com/spf13/cobra"
)

type root struct{}

func New() *cobra.Command {
	cmd := new(root)
	c := &cobra.Command{
		Short:   "nextunnel",
		Version: "v0.0.1",
		Args:    cobra.ExactArgs(0),
		Run:     cmd.run,
	}
	return c
}

func (*root) run(cmd *cobra.Command, _ []string) {
	_ = cmd.Help()
}
