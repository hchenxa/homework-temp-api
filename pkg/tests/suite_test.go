package tests_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"homework-temp-api/pkg/utils"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CSGHub Lite API Suite")
}

var _ = BeforeSuite(func() {
	resp, err := utils.Get("/api/health")
	Expect(err).NotTo(HaveOccurred(), "API server at %s is unreachable", utils.BaseURL())
	defer func() { _ = resp.Body.Close() }()
	Expect(resp.StatusCode).To(Equal(200), "API server at %s returned unexpected status", utils.BaseURL())

	body, err := utils.ReadBody(resp)
	Expect(err).NotTo(HaveOccurred())
	health, err := utils.ParseJSON[utils.HealthResponse](body)
	Expect(err).NotTo(HaveOccurred())
	Expect(health.Status).To(Equal("ok"))

	GinkgoWriter.Printf("API server is reachable at %s\n", utils.BaseURL())
})
