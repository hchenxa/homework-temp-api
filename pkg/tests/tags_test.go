package tests_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"homework-temp-api/pkg/utils"
)

var _ = Describe("GET /api/tags", Label("tags", "model-management", "get"), func() {

	It("should return 200 OK", func() {
		resp, err := utils.Get("/api/tags")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("should return a list of models", func() {
		resp, err := utils.Get("/api/tags")
		Expect(err).NotTo(HaveOccurred())

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		tags, err := utils.ParseJSON[utils.TagsResponse](body)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags.Models).NotTo(BeNil())
	})

	It("each model should have required fields", func() {
		resp, err := utils.Get("/api/tags")
		Expect(err).NotTo(HaveOccurred())

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		tags, err := utils.ParseJSON[utils.TagsResponse](body)
		Expect(err).NotTo(HaveOccurred())

		for _, m := range tags.Models {
			Expect(m.Name).NotTo(BeEmpty())
			Expect(m.Model).NotTo(BeEmpty())
			Expect(m.Size).To(BeNumerically(">", 0))
			Expect(m.Format).NotTo(BeEmpty())
			Expect(m.ModifiedAt).NotTo(BeEmpty())
		}
	})
})
