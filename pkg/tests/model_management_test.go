package tests_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"homework-temp-api/pkg/utils"
)

var _ = Describe("POST /api/show", Label("show", "model-management", "post"), func() {

	It("should return 404 for non-existent model", func() {
		resp, err := utils.Post("/api/show", utils.ShowRequest{Model: "non-existent-model"})
		Expect(err).NotTo(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()
		Expect(resp.StatusCode).To(BeNumerically(">=", 400))
	})

	It("should return model details for an existing model", func() {
		resp, err := utils.Post("/api/show", utils.ShowRequest{Model: utils.TestModel()})
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			Skip("test model is not available on the server")
		}
		Expect(err).NotTo(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		show, err := utils.ParseJSON[utils.ShowResponse](body)
		Expect(err).NotTo(HaveOccurred())
		Expect(show.Details).NotTo(BeNil())
		Expect(show.Details.Name).NotTo(BeEmpty())
		Expect(show.Details.Model).To(Equal(utils.TestModel()))
		Expect(show.Details.Size).To(BeNumerically(">", 0))
		Expect(show.Details.Format).To(Equal("gguf"))
		Expect(show.Details.ModifiedAt).NotTo(BeEmpty())
	})
})

var _ = Describe("POST /api/pull", Label("pull", "model-management", "post", "requires-model"), func() {

	It("should return SSE progress events", func() {
		resp, err := utils.Post("/api/pull", utils.PullRequest{Model: utils.TestModel()})
		Expect(err).NotTo(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		events := utils.ParseSSEEvents(body)
		Expect(events).NotTo(BeEmpty())

		var first utils.PullEvent
		Expect(json.Unmarshal([]byte(events[0]), &first)).To(Succeed())
		Expect(first.Status).NotTo(BeEmpty())

		var last utils.PullEvent
		Expect(json.Unmarshal([]byte(events[len(events)-1]), &last)).To(Succeed())
		Expect(last.Status).NotTo(BeEmpty())
	})
})

var _ = Describe("DELETE /api/delete", Label("delete", "model-management", "delete", "requires-model"), func() {

	BeforeEach(func() {
		resp, err := utils.Post("/api/show", utils.ShowRequest{Model: utils.TestModel()})
		if err != nil {
			Skip("server is unreachable")
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode == http.StatusNotFound {
			Skip("test model is not available on the server, skipping delete test")
		}
	})

	It("should return deleted status", func() {
		resp, err := utils.Delete("/api/delete", utils.DeleteRequest{Model: utils.TestModel()})
		Expect(err).NotTo(HaveOccurred())
		defer func() { _ = resp.Body.Close() }()

		body, err := utils.ReadBody(resp)
		Expect(err).NotTo(HaveOccurred())

		delResp, err := utils.ParseJSON[utils.DeleteResponse](body)
		Expect(err).NotTo(HaveOccurred())
		Expect(delResp.Status).To(Equal("deleted"))
	})
})
