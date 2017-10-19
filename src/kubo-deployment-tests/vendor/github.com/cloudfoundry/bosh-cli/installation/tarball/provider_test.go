package tarball_test

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/cloudfoundry/bosh-cli/installation/tarball"
	fakebiui "github.com/cloudfoundry/bosh-cli/ui/fakes"
	"github.com/cloudfoundry/bosh-utils/httpclient"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Provider", func() {
	var (
		server    *ghttp.Server
		provider  Provider
		cache     Cache
		fs        *fakesys.FakeFileSystem
		source    *fakeSource
		fakeStage *fakebiui.FakeStage
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		fs = fakesys.NewFakeFileSystem()
		logger := boshlog.NewLogger(boshlog.LevelNone)
		cache = NewCache(filepath.Join("/", "fake-base-path"), fs, logger)
		httpClient := httpclient.NewHTTPClient(httpclient.DefaultClient, logger)
		provider = NewProvider(cache, fs, httpClient, 3, 0, logger)
		fakeStage = fakebiui.NewFakeStage()
	})

	Describe("Get", func() {
		Context("when URL starts with file://", func() {
			BeforeEach(func() {
				source = newFakeSource("file://fake-file", "fake-sha1", "fake-description")
				fs.WriteFileString("expanded-file-path", "")
				fs.ExpandPathExpanded = "expanded-file-path"
			})

			It("returns expanded path to file", func() {
				path, err := provider.Get(source, fakeStage)
				Expect(err).ToNot(HaveOccurred())
				Expect(path).To(Equal("expanded-file-path"))
			})
		})

		Context("when URL starts with http(s)://", func() {
			BeforeEach(func() {
				source = newFakeSource(server.URL(), "da39a3ee5e6b4b0d3255bfef95601890afd80709", "fake-description")
			})

			Context("when tarball is present in cache", func() {
				BeforeEach(func() {
					fs.WriteFileString("fake-source-path", "")
					cache.Save("fake-source-path", source)
				})

				It("returns cached tarball path", func() {
					path, err := provider.Get(source, fakeStage)
					Expect(err).ToNot(HaveOccurred())
					shaSum := sha1.Sum([]byte(source.GetURL()))
					expectedFileName := fmt.Sprintf("%x-da39a3ee5e6b4b0d3255bfef95601890afd80709", string(shaSum[:]))
					Expect(path).To(Equal(filepath.Join("/", "fake-base-path", expectedFileName)))
				})

				It("skips downloading stage", func() {
					_, err := provider.Get(source, fakeStage)
					Expect(err).ToNot(HaveOccurred())

					Expect(fakeStage.PerformCalls[0].Name).To(Equal("Downloading fake-description"))
					Expect(fakeStage.PerformCalls[0].SkipError.Error()).To(Equal("Found in local cache: Already downloaded"))
				})
			})

			Context("when tarball is not present in cache", func() {
				var (
					tempDownloadFilePath1 string
					tempDownloadFilePath2 string
					tempDownloadFilePath3 string
				)

				BeforeEach(func() {
					tempDownloadFile1, err := ioutil.TempFile("", "temp-download-file1")
					Expect(err).ToNot(HaveOccurred())
					tempDownloadFile2, err := ioutil.TempFile("", "temp-download-file2")
					Expect(err).ToNot(HaveOccurred())
					tempDownloadFile3, err := ioutil.TempFile("", "temp-download-file3")
					Expect(err).ToNot(HaveOccurred())
					fs.ReturnTempFiles = []boshsys.File{tempDownloadFile1, tempDownloadFile2, tempDownloadFile3}
					tempDownloadFilePath1 = tempDownloadFile1.Name()
					tempDownloadFilePath2 = tempDownloadFile2.Name()
					tempDownloadFilePath3 = tempDownloadFile3.Name()
				})

				AfterEach(func() {
					os.RemoveAll(tempDownloadFilePath1)
					os.RemoveAll(tempDownloadFilePath2)
					os.RemoveAll(tempDownloadFilePath3)
				})

				Context("when downloading succeds", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.RespondWith(200, "fake-body"),
								ghttp.VerifyRequest("GET", "/"),
							),
							ghttp.RespondWith(200, "fake-body"),
							ghttp.RespondWith(200, "fake-body"),
						)
					})

					It("downloads tarball from given URL and returns saved cache tarball path", func() {
						path, err := provider.Get(source, fakeStage)
						Expect(err).ToNot(HaveOccurred())
						shaSum := sha1.Sum([]byte(source.GetURL()))
						expectedFileName := fmt.Sprintf("%x-da39a3ee5e6b4b0d3255bfef95601890afd80709", string(shaSum[:]))
						Expect(path).To(Equal(filepath.Join("/", "fake-base-path", expectedFileName)))
						Expect(server.ReceivedRequests()).To(HaveLen(1))
					})

					It("logs downloading stage", func() {
						_, err := provider.Get(source, fakeStage)
						Expect(err).ToNot(HaveOccurred())

						Expect(fakeStage.PerformCalls).To(Equal([]*fakebiui.PerformCall{
							{Name: "Downloading fake-description"},
						}))
					})

					Context("when sha1 does not match", func() {
						BeforeEach(func() {
							source = newFakeSource(server.URL(), "expectedsha1", "fake-description")
						})

						It("returns an error", func() {
							_, err := provider.Get(source, fakeStage)
							Expect(err).To(HaveOccurred())
							Expect(err.Error()).To(ContainSubstring("Failed to download from '%s': Verifying digest for downloaded file: Expected stream to have digest 'expectedsha1' but was 'da39a3ee5e6b4b0d3255bfef95601890afd80709'", server.URL()))
						})

						It("retries downloading up to 3 times", func() {
							_, err := provider.Get(source, fakeStage)
							Expect(err).To(HaveOccurred())

							Expect(server.ReceivedRequests()).To(HaveLen(3))
						})

						It("removes the downloaded file", func() {
							_, err := provider.Get(source, fakeStage)
							Expect(err).To(HaveOccurred())
							Expect(fs.FileExists(tempDownloadFilePath1)).To(BeFalse())
							Expect(fs.FileExists(tempDownloadFilePath2)).To(BeFalse())
							Expect(fs.FileExists(tempDownloadFilePath3)).To(BeFalse())
						})
					})

					Context("when saving to cache fails", func() {
						BeforeEach(func() {
							// Creating cache base directory fails
							fs.MkdirAllError = errors.New("fake-mkdir-error")
						})

						It("returns an error", func() {
							_, err := provider.Get(source, fakeStage)
							Expect(err).To(HaveOccurred())
							Expect(err.Error()).To(ContainSubstring("fake-mkdir-error"))
						})

						It("removes the downloaded file", func() {
							_, err := provider.Get(source, fakeStage)
							Expect(err).To(HaveOccurred())
							Expect(fs.FileExists(tempDownloadFilePath1)).To(BeFalse())
							Expect(fs.FileExists(tempDownloadFilePath2)).To(BeFalse())
							Expect(fs.FileExists(tempDownloadFilePath3)).To(BeFalse())
						})
					})
				})

				Context("when downloading fails", func() {
					disconnectingRequestHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						conn, _, err := w.(http.Hijacker).Hijack()
						Expect(err).NotTo(HaveOccurred())

						conn.Close()
					})

					BeforeEach(func() {
						server.AppendHandlers(
							ghttp.CombineHandlers(
								ghttp.VerifyRequest("GET", "/"),
								disconnectingRequestHandler,
							),
							disconnectingRequestHandler,
							disconnectingRequestHandler,
						)
					})

					It("retries downloading up to 3 times", func() {
						_, err := provider.Get(source, fakeStage)
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(ContainSubstring("Get %s: EOF", server.URL())))

						Expect(server.ReceivedRequests()).To(HaveLen(3))
					})

					It("removes the downloaded file", func() {
						_, err := provider.Get(source, fakeStage)
						Expect(err).To(HaveOccurred())
						Expect(fs.FileExists(tempDownloadFilePath1)).To(BeFalse())
						Expect(fs.FileExists(tempDownloadFilePath2)).To(BeFalse())
						Expect(fs.FileExists(tempDownloadFilePath3)).To(BeFalse())
					})
				})
			})
		})

		Context("when URL does not start with either file:// or http(s)://", func() {
			BeforeEach(func() {
				source = newFakeSource("invalid-url", "fake-sha1", "fake-description")
			})

			It("returns an error", func() {
				_, err := provider.Get(source, fakeStage)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid source URL: 'invalid-url'"))
			})
		})
	})
})

type fakeSource struct {
	url         string
	sha1        string
	description string
}

func newFakeSource(url, sha1, description string) *fakeSource {
	return &fakeSource{url, sha1, description}
}

func (s *fakeSource) GetURL() string      { return s.url }
func (s *fakeSource) GetSHA1() string     { return s.sha1 }
func (s *fakeSource) Description() string { return s.description }
