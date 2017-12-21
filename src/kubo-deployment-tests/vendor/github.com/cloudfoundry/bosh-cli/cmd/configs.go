package cmd

import (
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
)

type ConfigsCmd struct {
	ui       boshui.UI
	director boshdir.Director
}

func NewConfigsCmd(ui boshui.UI, director boshdir.Director) ConfigsCmd {
	return ConfigsCmd{ui: ui, director: director}
}

func (c ConfigsCmd) Run(opts ConfigsOpts) error {
	filter := boshdir.ConfigsFilter{
		Type: opts.Type,
		Name: opts.Name,
	}

	configs, err := c.director.ListConfigs(filter)
	if err != nil {
		return err
	}

	table := boshtbl.Table{
		Content: "configs",
		Header: []boshtbl.Header{
			boshtbl.NewHeader("Type"),
			boshtbl.NewHeader("Name"),
		},
	}

	for _, config := range configs {
		table.Rows = append(table.Rows, []boshtbl.Value{
			boshtbl.NewValueString(config.Type),
			boshtbl.NewValueString(config.Name),
		})
	}

	c.ui.PrintTable(table)
	return nil
}
