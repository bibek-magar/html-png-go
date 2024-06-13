package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/convert", convertHTMLToImage)
	router.Run(":8080")
}

func convertHTMLToImage(c *gin.Context) {
	var json struct {
		HTMLContent string `json:"html_content"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Allocate a buffer for the screenshot
	var buf []byte
	err := chromedp.Run(ctx, captureScreenshot(json.HTMLContent, &buf))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to capture screenshot"})
		return
	}

	// Return the image as a response
	c.Data(http.StatusOK, "image/png", buf)
}

func captureScreenshot(html string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("data:text/html," + url.PathEscape(html)),
		chromedp.EmulateViewport(1024, 768),
		chromedp.Sleep(2 * time.Second), // Wait for the content to load
		chromedp.CaptureScreenshot(res),
	}
}
