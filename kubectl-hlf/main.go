package main

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

func main() {
	if err := cmd.NewCmdHLF(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}).Execute(); err != nil {
		os.Exit(1)
	}
}
