package tests_test

import (
	"encoding/json"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"homework-temp-api/pkg/utils"
)

var _ = Describe("/api/chat", Label("chat", "inference", "requires-model"), func() {

	var testModel string

	BeforeEach(func() {
		testModel = utils.TestModel()
		if testModel == "" {
			Skip("TEST_MODEL is not set")
		}
	})

	Context("non-streaming", Label("non-streaming"), func() {

		It("should return 200 with valid response structure", func() {
			resp, err := utils.Post("/api/chat", utils.ChatRequest{
				Model: testModel,
				Messages: []utils.ChatMessage{
					{Role: "user", Content: "Say hello in one word"},
				},
				Stream: false,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := utils.ReadBody(resp)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).NotTo(BeEmpty())

			chatResp, err := utils.ParseJSON[utils.ChatResponse](body)
			Expect(err).NotTo(HaveOccurred())
			Expect(chatResp.Model).To(Equal(testModel))
			Expect(chatResp.Message).NotTo(BeNil())
			Expect(chatResp.Message.Role).To(Equal("assistant"))
			Expect(chatResp.Message.Content).NotTo(BeEmpty())
			Expect(chatResp.Done).To(BeTrue())
			Expect(chatResp.CreatedAt).NotTo(BeEmpty())
		})

		It("should support multi-turn conversation", func() {
			resp, err := utils.Post("/api/chat", utils.ChatRequest{
				Model: testModel,
				Messages: []utils.ChatMessage{
					{Role: "user", Content: "My name is Alice."},
					{Role: "assistant", Content: "Hello Alice! Nice to meet you."},
					{Role: "user", Content: "What is my name?"},
				},
				Stream: false,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := utils.ReadBody(resp)
			Expect(err).NotTo(HaveOccurred())

			chatResp, err := utils.ParseJSON[utils.ChatResponse](body)
			Expect(err).NotTo(HaveOccurred())
			Expect(chatResp.Message.Content).NotTo(BeEmpty())
		})

		It("should accept generation options", func() {
			resp, err := utils.Post("/api/chat", utils.ChatRequest{
				Model: testModel,
				Messages: []utils.ChatMessage{
					{Role: "user", Content: "Count from 1 to 3"},
				},
				Stream:  false,
				Options: map[string]any{"temperature": 0.1, "max_tokens": 50},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("streaming", Label("streaming"), func() {

		It("should return SSE events when stream is true (default)", func() {
			start := time.Now()
			resp, err := utils.Post("/api/chat", utils.ChatRequest{
				Model: testModel,
				Messages: []utils.ChatMessage{
					{Role: "user", Content: "Say hi in one word"},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(Or(
				Equal("text/event-stream"),
				ContainSubstring("text/event-stream"),
			))

			body, err := utils.ReadBody(resp)
			Expect(err).NotTo(HaveOccurred())

			events := utils.ParseSSEEvents(body)
			Expect(events).NotTo(BeEmpty())

			lastEvent := events[len(events)-1]
			var final struct {
				Model     string `json:"model"`
				Done      bool   `json:"done"`
				CreatedAt string `json:"created_at"`
			}
			Expect(json.Unmarshal([]byte(lastEvent), &final)).To(Succeed())
			Expect(final.Done).To(BeTrue())
			Expect(final.Model).NotTo(BeEmpty())

			GinkgoWriter.Printf("Chat streamed %d events in %v\n", len(events), time.Since(start))
		})
	})
})

var _ = Describe("/api/generate", Label("generate", "inference", "requires-model"), func() {

	var testModel string

	BeforeEach(func() {
		testModel = utils.TestModel()
		if testModel == "" {
			Skip("TEST_MODEL is not set")
		}
	})

	Context("non-streaming", Label("non-streaming"), func() {

		It("should return 200 with valid response structure", func() {
			resp, err := utils.Post("/api/generate", utils.GenerateRequest{
				Model:  testModel,
				Prompt: "Write a haiku about programming",
				Stream: false,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, err := utils.ReadBody(resp)
			Expect(err).NotTo(HaveOccurred())

			genResp, err := utils.ParseJSON[utils.GenerateResponse](body)
			Expect(err).NotTo(HaveOccurred())
			Expect(genResp.Model).To(Equal(testModel))
			Expect(genResp.Response).NotTo(BeEmpty())
			Expect(genResp.Done).To(BeTrue())
			Expect(genResp.CreatedAt).NotTo(BeEmpty())
		})

		It("should accept generation options", func() {
			resp, err := utils.Post("/api/generate", utils.GenerateRequest{
				Model:   testModel,
				Prompt:  "Repeat the word hello",
				Stream:  false,
				Options: map[string]any{"temperature": 0.1, "max_tokens": 20},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("streaming", Label("streaming"), func() {

		It("should return SSE events", func() {
			start := time.Now()
			resp, err := utils.Post("/api/generate", utils.GenerateRequest{
				Model:  testModel,
				Prompt: "Say hello in one word",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(Or(
				Equal("text/event-stream"),
				ContainSubstring("text/event-stream"),
			))

			body, err := utils.ReadBody(resp)
			Expect(err).NotTo(HaveOccurred())

			events := utils.ParseSSEEvents(body)
			Expect(events).NotTo(BeEmpty())

			lastEvent := events[len(events)-1]
			var final struct {
				Model     string `json:"model"`
				Done      bool   `json:"done"`
				CreatedAt string `json:"created_at"`
				Response  string `json:"response"`
			}
			Expect(json.Unmarshal([]byte(lastEvent), &final)).To(Succeed())
			Expect(final.Done).To(BeTrue())
			Expect(final.Model).NotTo(BeEmpty())

			GinkgoWriter.Printf("Generate streamed %d events in %v\n", len(events), time.Since(start))
		})
	})
})
