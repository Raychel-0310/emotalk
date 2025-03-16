package main

import (
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"strings"
	"io"
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render template
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// emotion analyze function
func analyzeSentiment(text string) string {
	positiveWords := []string{"happy", "joyful", "great", "awesome", "excited"}
	negativeWords := []string{"sad", "angry", "terrible", "awful", "disappointed"}

	// compare with lower case
	text = strings.ToLower(text)

	// analyze emotion
	positiveScore := 0
	negativeScore := 0

	for _, word := range positiveWords {
		if strings.Contains(text, word) {
			positiveScore++
		}
	}

	for _, word := range negativeWords {
		if strings.Contains(text, word) {
			negativeScore++
		}
	}

	// judge positive or negative
	if positiveScore > negativeScore {
		return "Positive"
	} else if negativeScore > positiveScore {
		return "Negative"
	}
	return "Neutral"
}

// handler (main page)
func mainPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// handler (receive message and respond result for emotion analyzation)
func analyzeMessage(c echo.Context) error {
	message := c.FormValue("message")
	sentiment := analyzeSentiment(message)
	return c.JSON(http.StatusOK, map[string]string{
		"sentiment": sentiment,
	})
}

func main() {
	e := echo.New()

	// read template
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Static("/static", "static")

	e.GET("/", mainPage)
	e.POST("/analyze", analyzeMessage)

	e.Logger.Fatal(e.Start(":8080"))
}
