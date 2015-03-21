package ratelimit

import (
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RateLimiter", func() {

	It("should accurately rate-limit at small rates", func() {
		n := 100000
		rl := New(n, time.Minute)
		for i := 0; i < n; i++ {
			Expect(rl.Limit()).To(BeFalse(), "on cycle %d", i)
		}
		Expect(rl.Limit()).To(BeTrue())
	})

	It("should accurately rate-limit at large rates", func() {
		n := 100000
		rl := New(n, time.Hour)
		for i := 0; i < n; i++ {
			Expect(rl.Limit()).To(BeFalse(), "on cycle %d", i)
		}
		Expect(rl.Limit()).To(BeTrue())
	})

	It("should correctly increase allowance", func() {
		n := 25
		rl := New(n, 50*time.Millisecond)
		for i := 0; i < n; i++ {
			Expect(rl.Limit()).To(BeFalse(), "on cycle %d", i)
		}
		Expect(rl.Limit()).To(BeTrue())
		Eventually(rl.Limit, "60ms", "10ms").Should(BeFalse())
	})

	It("should be thread-safe", func() {
		c := 10
		n := 10000
		wg := sync.WaitGroup{}
		rl := New(c*n, time.Minute)
		for i := 0; i < c; i++ {
			wg.Add(1)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				for j := 0; j < n; j++ {
					Expect(rl.Limit()).To(BeFalse(), "thread %d, cycle %d", i, j)
				}
			}()
		}
		wg.Wait()
		Expect(rl.Limit()).To(BeTrue())
	})

})

// --------------------------------------------------------------------

func TestGinkgoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "github.com/bsm/ratelimit")
}
