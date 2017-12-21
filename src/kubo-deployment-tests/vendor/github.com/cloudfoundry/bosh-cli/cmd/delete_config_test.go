package cmd_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	fakedir "github.com/cloudfoundry/bosh-cli/director/directorfakes"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
)

var _ = Describe("DeleteConfigCmd", func() {
	var (
		ui       *fakeui.FakeUI
		director *fakedir.FakeDirector
		command  DeleteConfigCmd
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		director = &fakedir.FakeDirector{}
		command = NewDeleteConfigCmd(ui, director)
	})

	Describe("Run", func() {
		var (
			opts DeleteConfigOpts
		)

		BeforeEach(func() {
			opts = DeleteConfigOpts{
				Args: DeleteConfigArgs{
					Type: "my-type",
				},
				Name: "my-name",
			}
		})

		act := func() error { return command.Run(opts) }

		It("does not stop if confirmation is rejected", func() {
			ui.AskedConfirmationErr = errors.New("stop")

			err := act()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("stop"))

			Expect(director.DeleteConfigCallCount()).To(Equal(0))
		})

		Context("when type and name are given", func() {

			Context("when there is a matching config", func() {
				It("succeeds", func() {
					director.DeleteConfigReturns(true, nil)

					err := act()
					Expect(err).To(Not(HaveOccurred()))
				})
			})

			Context("when there is NO matching config", func() {
				It("succeeds with a message as hint", func() {
					director.DeleteConfigReturns(false, nil)

					err := act()
					Expect(err).To(Not(HaveOccurred()))
					Expect(ui.Said[0]).To(ContainSubstring("No configs to delete: no matches for type 'my-type' and name 'my-name' found."))
				})
			})

			It("fails", func() {
				director.DeleteConfigReturns(false, errors.New("fake-err"))

				err := act()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-err"))
			})
		})
	})
})
