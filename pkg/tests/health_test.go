package tests_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"homework-temp-api/pkg/utils"
)

var _ = Describe("GET /api/health", Label("health", "service-management", "get"), func() {

	It("should return 200 OK", func() {
		resp, err := utils.Get("/api/health")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("should return status ok", func() {
		resp, err := utils.Get("/api/health")
		Expect(err).NotTo(HaveOccurred())

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		health, err := utils.ParseJSON[utils.HealthResponse](body)
		Expect(err).NotTo(HaveOccurred())
		Expect(health.Status).To(Equal("ok"))
	})
})
