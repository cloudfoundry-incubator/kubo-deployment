package work_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/bosh-cli/work"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

var _ = Describe("Pool", func() {
	It("runs the given tasks", func() {
		pool := work.Pool{
			Count: 2,
		}

		resultsChan := make(chan int, 3)

		err := pool.ParallelDo(
			func() error {
				resultsChan <- 1
				return nil
			},
			func() error {
				resultsChan <- 2
				return nil
			},
			func() error {
				resultsChan <- 3
				return nil
			},
		)
		Expect(err).ToNot(HaveOccurred())

		Expect(resultsChan).To(Receive(Equal(1)))
		Expect(resultsChan).To(Receive(Equal(2)))
		Expect(resultsChan).To(Receive(Equal(3)))
	})

	It("bubbles up any errors", func() {
		pool := work.Pool{
			Count: 2,
		}

		err := pool.ParallelDo(
			func() error {
				return nil
			},
			func() error {
				return bosherr.ComplexError{
					Err:   errors.New("fake-error"),
					Cause: errors.New("fake-cause"),
				}
			},
			func() error {
				return nil
			},
		)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-error"))
		Expect(err.Error()).To(ContainSubstring("fake-cause"))
	})

	It("stops working after the first error", func() {
		pool := work.Pool{
			Count: 1, // Force serial run
		}

		err := pool.ParallelDo(
			func() error {
				return nil
			},
			func() error {
				return bosherr.ComplexError{
					Err:   errors.New("fake-error"),
					Cause: errors.New("fake-cause"),
				}
			},
			func() error {
				Fail("Expected third test to not run")
				return nil
			},
			func() error {
				Fail("Expected fourth test to not run")
				return nil
			},
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-error"))
		Expect(err.Error()).To(ContainSubstring("fake-cause"))
	})
})
