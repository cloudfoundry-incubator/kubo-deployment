package cmd

import (
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type ConfigCmd struct {
	ui       boshui.UI
	director boshdir.Director
}

func NewConfigCmd(ui boshui.UI, director boshdir.Director) ConfigCmd {
	return ConfigCmd{ui: ui, director: director}
}

func (c ConfigCmd) Run(opts ConfigOpts) error {
	config, err := c.director.LatestConfig(opts.Args.Type, opts.Name)
	if err != nil {
		return err
	}

	c.ui.PrintBlock([]byte(config.Content))
	return nil
}
