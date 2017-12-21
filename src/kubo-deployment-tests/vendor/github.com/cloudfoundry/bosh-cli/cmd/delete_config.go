package cmd

import (
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type DeleteConfigCmd struct {
	ui       boshui.UI
	director boshdir.Director
}

func NewDeleteConfigCmd(ui boshui.UI, director boshdir.Director) DeleteConfigCmd {
	return DeleteConfigCmd{ui: ui, director: director}
}

func (c DeleteConfigCmd) Run(opts DeleteConfigOpts) error {
	err := c.ui.AskForConfirmation()
	if err != nil {
		return err
	}

	deleted, err := c.director.DeleteConfig(opts.Args.Type, opts.Name)

	if !deleted {
		c.ui.PrintLinef("No configs to delete: no matches for type '%s' and name '%s' found.", opts.Args.Type, opts.Name)
	}

	if err != nil {
		return err
	}

	return nil
}
