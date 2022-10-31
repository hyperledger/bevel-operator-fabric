package follower

import (
	"context"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteOptions struct {
	Name string
}

func (o DeleteOptions) Validate() error {
	return nil
}

type deleteCmd struct {
	out         io.Writer
	errOut      io.Writer
	channelOpts DeleteOptions
}

func (c *deleteCmd) validate() error {
	return c.channelOpts.Validate()
}
func (c *deleteCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = oclient.HlfV1alpha1().
		FabricFollowerChannels().
		Delete(
			ctx,
			c.channelOpts.Name,
			v1.DeleteOptions{},
		)
	if err != nil {
		return err
	}
	log.Infof("Follower channel %s deleted", c.channelOpts.Name)
	return nil
}
func newDeleteFollowerChannelCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := deleteCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a follower channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.channelOpts.Name, "name", "", "Name of the Follower Channel to delete")
	return cmd
}
