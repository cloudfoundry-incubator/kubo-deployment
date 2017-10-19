package integration_test

import (
	"net/http"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/cloudfoundry/bosh-cli/cmd"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	fakeui "github.com/cloudfoundry/bosh-cli/ui/fakes"
)

var _ = Describe("task command", func() {
	var (
		ui         *fakeui.FakeUI
		fs         boshsys.FileSystem
		deps       BasicDeps
		cmdFactory Factory
	)

	BeforeEach(func() {
		ui = &fakeui.FakeUI{}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		confUI := boshui.NewWrappingConfUI(ui, logger)

		fs = boshsys.NewOsFileSystem(logger)
		deps = NewBasicDepsWithFS(confUI, fs, logger)
		cmdFactory = NewFactory(deps)
	})

	execCmd := func(args []string) {
		cmd, err := cmdFactory.New(args)
		Expect(err).ToNot(HaveOccurred())

		err = cmd.Execute()
		Expect(err).ToNot(HaveOccurred())
	}

	It("streams task output", func() {
		directorCACert, director := BuildHTTPSServer()
		defer director.Close()

		processing := ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/tasks/123"),
			ghttp.RespondWith(http.StatusOK, `{"id":123, "state":"processing"}`),
		)

		director.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/info"),
				ghttp.RespondWith(http.StatusOK, `{"user_authentication":{"type":"basic","options":{}}}`),
			),
			processing,
			processing,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusRequestedRangeNotSatisfiable, "Byte range unsatisfiable\n"),
			),
			processing,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusOK, `{}`+"\n"),
			),
			processing,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusOK, `{"time":1503082451,"stage":"event-one`),
			),
			processing,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusOK, ""),
			),
			processing,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusOK, `","tags":[],"total":1,"task":"event-one-task","state":"started","progress":0}`+"\n{}\n"),
			),
			processing,
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusOK, ""),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123"),
				ghttp.RespondWith(http.StatusOK, `{"id":123, "state":"done"}`),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/tasks/123/output", "type=event"),
				ghttp.RespondWith(http.StatusOK, `{"time":1503082451,"stage":"event-two","tags":[],"total":1,"task":"event-two-task","index":1,"state":"started","progress":0}`+"\n"),
			),
		)

		execCmd([]string{"task", "123", "-e", director.URL(), "--ca-cert", directorCACert})

		output := strings.Join(ui.Blocks, "\n")
		Expect(output).To(ContainSubstring("event-one"))
		Expect(output).To(ContainSubstring("event-two"))
	})
})
