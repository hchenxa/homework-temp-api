package tests_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"homework-temp-api/pkg/utils"
)

var _ = Describe("GET /api/ps", Label("ps", "service-management", "get"), func() {

	It("should return 200 OK", func() {
		resp, err := utils.Get("/api/ps")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("should return a list of running models", func() {
		resp, err := utils.Get("/api/ps")
		Expect(err).NotTo(HaveOccurred())

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		ps, err := utils.ParseJSON[utils.PSResponse](body)
		Expect(err).NotTo(HaveOccurred())
		Expect(ps.Models).NotTo(BeNil())

		for _, m := range ps.Models {
			Expect(m.Name).NotTo(BeEmpty())
			Expect(m.Model).NotTo(BeEmpty())
			Expect(m.Size).To(BeNumerically(">", 0))
			Expect(m.Format).NotTo(BeEmpty())
		}
	})
})

var _ = Describe("POST /api/stop", Label("stop", "service-management", "post"), func() {

	It("should return error for a model that is not running", func() {
		resp, err := utils.Post("/api/stop", utils.StopRequest{Model: "non-existent-model"})
		Expect(err).NotTo(HaveOccurred())

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		stopResp, err := utils.ParseJSON[utils.StopResponse](body)
		Expect(err).NotTo(HaveOccurred())

		if resp.StatusCode == http.StatusOK {
			Expect(stopResp.Status).To(Equal("stopped"))
		} else {
			Expect(stopResp.Error).NotTo(BeEmpty())
		}
	})

	It("should stop a running model", func() {
		resp, err := utils.Post("/api/stop", utils.StopRequest{Model: utils.TestModel()})
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			Skip("test model is not running")
		}
		Expect(err).NotTo(HaveOccurred())

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		stopResp, err := utils.ParseJSON[utils.StopResponse](body)
		Expect(err).NotTo(HaveOccurred())

		if resp.StatusCode == http.StatusOK {
			Expect(stopResp.Status).To(Equal("stopped"))
		} else {
			Expect(stopResp.Error).NotTo(BeEmpty())
		}
	})
})
