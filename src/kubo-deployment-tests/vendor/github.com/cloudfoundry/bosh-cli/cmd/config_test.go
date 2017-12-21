package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	fakedir "github.com/cloudfoundry/bosh-cli/director/directorfakes"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
)

var _ = Describe("ConfigCmd", func() {
	var (
		ui       *fakeui.FakeUI
		director *fakedir.FakeDirector
		command  ConfigCmd
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		director = &fakedir.FakeDirector{}
		command = NewConfigCmd(ui, director)
	})

	Describe("Run", func() {
		var (
			opts ConfigOpts
		)

		BeforeEach(func() {
			opts = ConfigOpts{
				Args: ConfigArgs{Type: "my-type"},
				Name: "",
			}
		})

		act := func() error { return command.Run(opts) }

		It("shows config", func() {
			config := boshdir.Config{
				Content: "some-content",
			}

			director.LatestConfigReturns(config, nil)

			err := act()
			Expect(err).ToNot(HaveOccurred())
			Expect(ui.Blocks).To(Equal([]string{"some-content"}))
		})

		It("returns error if config cannot be retrieved", func() {
			director.LatestConfigReturns(boshdir.Config{}, errors.New("fake-err"))

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-err"))
		})
	})
})
