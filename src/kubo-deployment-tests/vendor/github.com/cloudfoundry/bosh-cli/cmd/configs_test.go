package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	fakedir "github.com/cloudfoundry/bosh-cli/director/directorfakes"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
	boshtbl "github.com/cloudfoundry/bosh-cli/ui/table"
)

var _ = Describe("ConfigsCmd", func() {
	var (
		ui       *fakeui.FakeUI
		director *fakedir.FakeDirector
		command  ConfigsCmd
		configs  []boshdir.ConfigListItem
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		director = &fakedir.FakeDirector{}
		command = NewConfigsCmd(ui, director)
	})

	Describe("Run", func() {
		var (
			opts ConfigsOpts
		)

		BeforeEach(func() {
			opts = ConfigsOpts{}
			configs = []boshdir.ConfigListItem{boshdir.ConfigListItem{Type: "my-type", Name: "some-name"}, boshdir.ConfigListItem{Type: "my-type", Name: "other-name"}}
		})

		act := func() error { return command.Run(opts) }

		It("lists configs", func() {
			director.ListConfigsReturns(configs, nil)

			err := act()
			Expect(err).ToNot(HaveOccurred())
			Expect(director.ListConfigsCallCount()).To(Equal(1))
			Expect(director.ListConfigsArgsForCall(0)).To(Equal(boshdir.ConfigsFilter{}))

			Expect(ui.Table).To(Equal(boshtbl.Table{
				Content: "configs",

				Header: []boshtbl.Header{
					boshtbl.NewHeader("Type"),
					boshtbl.NewHeader("Name"),
				},

				Rows: [][]boshtbl.Value{
					{
						boshtbl.NewValueString("my-type"),
						boshtbl.NewValueString("some-name"),
					},
					{
						boshtbl.NewValueString("my-type"),
						boshtbl.NewValueString("other-name"),
					},
				},
			}))
		})

		It("returns error if configs cannot be listed", func() {
			director.ListConfigsReturns([]boshdir.ConfigListItem{}, errors.New("fake-err"))

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})

		Context("When filtering for type", func() {
			BeforeEach(func() {
				opts = ConfigsOpts{
					Type: "my-type",
				}
				configs = []boshdir.ConfigListItem{boshdir.ConfigListItem{Type: "my-type", Name: "some-name"}}
			})

			It("applies filters for just type", func() {
				director.ListConfigsReturns(configs, nil)

				err := act()
				Expect(err).ToNot(HaveOccurred())
				Expect(director.ListConfigsCallCount()).To(Equal(1))
				Expect(director.ListConfigsArgsForCall(0)).To(Equal(boshdir.ConfigsFilter{Type: "my-type"}))

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "configs",

					Header: []boshtbl.Header{
						boshtbl.NewHeader("Type"),
						boshtbl.NewHeader("Name"),
					},

					Rows: [][]boshtbl.Value{
						{
							boshtbl.NewValueString("my-type"),
							boshtbl.NewValueString("some-name"),
						},
					},
				}))
			})
		})

		Context("When filtering for name", func() {
			BeforeEach(func() {
				opts = ConfigsOpts{
					Name: "some-name",
				}
				configs = []boshdir.ConfigListItem{boshdir.ConfigListItem{Type: "my-type", Name: "some-name"}}
			})

			It("applies filters for just name", func() {
				director.ListConfigsReturns(configs, nil)

				err := act()
				Expect(err).ToNot(HaveOccurred())
				Expect(director.ListConfigsCallCount()).To(Equal(1))
				Expect(director.ListConfigsArgsForCall(0)).To(Equal(boshdir.ConfigsFilter{Name: "some-name"}))

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "configs",

					Header: []boshtbl.Header{
						boshtbl.NewHeader("Type"),
						boshtbl.NewHeader("Name"),
					},

					Rows: [][]boshtbl.Value{
						{
							boshtbl.NewValueString("my-type"),
							boshtbl.NewValueString("some-name"),
						},
					},
				}))
			})
		})

		Context("When filtering for both, type and name", func() {
			BeforeEach(func() {
				opts = ConfigsOpts{
					Type: "my-type",
					Name: "some-name",
				}
				configs = []boshdir.ConfigListItem{boshdir.ConfigListItem{Type: "my-type", Name: "some-name"}}
			})

			It("applies filters for type and name", func() {
				director.ListConfigsReturns(configs, nil)

				err := act()
				Expect(err).ToNot(HaveOccurred())
				Expect(director.ListConfigsCallCount()).To(Equal(1))
				Expect(director.ListConfigsArgsForCall(0)).To(Equal(boshdir.ConfigsFilter{Type: "my-type", Name: "some-name"}))

				Expect(ui.Table).To(Equal(boshtbl.Table{
					Content: "configs",

					Header: []boshtbl.Header{
						boshtbl.NewHeader("Type"),
						boshtbl.NewHeader("Name"),
					},

					Rows: [][]boshtbl.Value{
						{
							boshtbl.NewValueString("my-type"),
							boshtbl.NewValueString("some-name"),
						},
					},
				}))
			})
		})
	})
})
