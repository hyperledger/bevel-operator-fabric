package org

import (
	"github.com/spf13/cobra"
	"io"
)

func NewOrgCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "org",
	}
	cmd.AddCommand(newOrgInspectCmd(out, errOut))
	return cmd
}
