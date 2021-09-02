package utils

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Options struct {
	config   string
	userPath string
	mspID    string
	userName string
}

func (o Options) Validate() error {
	if o.config == "" {
		return errors.New("--config is required")
	}
	if o.userPath == "" {
		return errors.New("--userPath is required")
	}
	if o.mspID == "" {
		return errors.New("--mspid is required")
	}
	return nil
}

type addUserCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   Options
}

func (c *addUserCmd) validate() error {
	return c.opts.Validate()
}

func (c *addUserCmd) run(args []string) error {
	configBytes, err := ioutil.ReadFile(c.opts.config)
	if err != nil {
		return err
	}
	networkConfigMap := map[string]interface{}{}
	err = yaml.Unmarshal(configBytes, networkConfigMap)
	if err != nil {
		return err
	}
	userMap := map[string]interface{}{}
	userBytes, err := ioutil.ReadFile(c.opts.userPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(userBytes, userMap)
	if err != nil {
		return err
	}
	orgsMap := networkConfigMap["organizations"].(map[string]interface{})
	orgMap := orgsMap[c.opts.mspID].(map[string]interface{})
	users := orgMap["users"].(map[string]interface{})

	users[c.opts.userName] = userMap
	configBytesNew, err :=  yaml.Marshal(networkConfigMap)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.opts.config, configBytesNew, 0777)
	if err != nil {
		return err
	}
	return nil
}

func newAddUserCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := addUserCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "adduser",
		Short: "Adds a previously enrolled user to the network config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.opts.userPath, "userPath", "", "the output of enrolling a user")
	f.StringVar(&c.opts.userName, "username", "", "the username")
	f.StringVar(&c.opts.config, "config", "", "networkconfig, you can use inspect to get the networkconfig")
	f.StringVar(&c.opts.mspID, "mspid", "", "the organization where the user will be added")
	return cmd
}
