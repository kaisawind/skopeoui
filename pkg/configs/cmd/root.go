package cmd

import (
	"github.com/spf13/cobra"
)

type RunECallback func(cmd *cobra.Command, args []string) error
type env struct {
	Usage   string
	Default string
	Value   *string
}
