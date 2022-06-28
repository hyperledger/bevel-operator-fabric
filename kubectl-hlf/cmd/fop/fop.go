package fop

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/fop/export"
	"github.com/spf13/cobra"
	"io"
)

func NewFOPCMD(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	fopCmd := &cobra.Command{
		Use: "fop",
	}
	fopCmd.AddCommand(
		export.NewExportCmd(stdOut, stdErr),
	)
	return fopCmd
}
