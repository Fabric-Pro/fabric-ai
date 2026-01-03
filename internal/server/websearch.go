package restapi

import (
	"net/http"

	"github.com/danielmiessler/fabric/internal/core"
	"github.com/gin-gonic/gin"
)

type WebSearchHandler struct {
	registry *core.PluginRegistry
}

// SearchRequest represents a request to search the web
type SearchRequest struct {
	Question string `json:"question" binding:"required" example:"What is the capital of France?"` // Search question (required)
}

// SearchResponse represents the search results response
type SearchResponse struct {
	Content  string `json:"content" example:"The capital of France is Paris..."` // Search results content
	Question string `json:"question" example:"What is the capital of France?"`   // Original search question
}

// ScrapeRequest represents a request to scrape a webpage
type ScrapeRequest struct {
	URL string `json:"url" binding:"required" example:"https://example.com"` // URL to scrape (required)
}

// ScrapeResponse represents the scraped content response
type ScrapeResponse struct {
	Content string `json:"content" example:"This is the main content of the page..."` // Scraped content
	URL     string `json:"url" example:"https://example.com"`                         // Original URL
}

func NewWebSearchHandler(r *gin.Engine, registry *core.PluginRegistry) *WebSearchHandler {
	handler := &WebSearchHandler{registry: registry}
	r.POST("/search", handler.Search)
	r.POST("/scrape", handler.Scrape)
	return handler
}

// Search godoc
// @Summary Search the web using Jina AI
// @Description Searches the web using Jina AI's search API and returns clean, LLM-friendly results
// @Tags websearch
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Search request with question"
// @Success 200 {object} SearchResponse "Successful search results"
// @Failure 400 {object} map[string]string "Bad request - invalid or missing question"
// @Failure 500 {object} map[string]string "Internal server error - search failed"
// @Security ApiKeyAuth
// @Router /search [post]
func (h *WebSearchHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question is required"})
		return
	}

	content, err := h.registry.Jina.ScrapeQuestion(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SearchResponse{
		Content:  content,
		Question: req.Question,
	})
}

// Scrape godoc
// @Summary Scrape a webpage using Jina AI
// @Description Scrapes a webpage using Jina AI's reader API and returns clean, LLM-friendly text
// @Tags websearch
// @Accept json
// @Produce json
// @Param request body ScrapeRequest true "Scrape request with URL"
// @Success 200 {object} ScrapeResponse "Successful scrape results"
// @Failure 400 {object} map[string]string "Bad request - invalid or missing URL"
// @Failure 500 {object} map[string]string "Internal server error - scrape failed"
// @Security ApiKeyAuth
// @Router /scrape [post]
func (h *WebSearchHandler) Scrape(c *gin.Context) {
	var req ScrapeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	content, err := h.registry.Jina.ScrapeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ScrapeResponse{
		Content: content,
		URL:     req.URL,
	})
}
