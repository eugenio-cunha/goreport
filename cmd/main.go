package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"

	"net/http"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"github.com/gin-gonic/gin"
)

func getHtml(value string, data interface{}) string {
	tmpl := template.Must(template.New("html").Parse(value))

	var tpl bytes.Buffer
	tmpl.Execute(&tpl, data)
	return tpl.String()
}

func report(c *gin.Context) {
	var res map[string]interface{}
	if err := json.NewDecoder(c.Request.Body).Decode(&res); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	var (
		header = ""
		layout = ""
		footer = ""
		data   = new(map[string]interface{})
	)

	if res["data"] != nil && res["data"].(string) != "" {
		err := json.Unmarshal([]byte(res["data"].(string)), data)
		if err != nil {
			c.JSON(500, gin.H{
				"messagew": err.Error(),
				"data":     res["data"].(string),
			})
		}
	}

	if res["header"] != nil && res["header"].(string) != "" {
		header = getHtml(res["header"].(string), data)
	}

	if res["body"] != nil && res["body"].(string) != "" {
		layout = getHtml(res["body"].(string), data)
	}

	if res["footer"] != nil && res["footer"].(string) != "" {
		footer = getHtml(res["footer"].(string), data)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("data:text/html,%s", layout)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			b, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithScale(1).
				WithDisplayHeaderFooter(true).
				WithHeaderTemplate(header).
				WithFooterTemplate(footer).Do(ctx)
			if err != nil {
				c.JSON(500, gin.H{
					"message": err.Error(),
				})
				return nil
			}

			_, err = c.Writer.Write(b)
			if err != nil {
				c.JSON(500, gin.H{
					"message": err.Error(),
				})
				return nil
			}

			return nil
		}),
	); err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	}
}

func main() {
	router := gin.Default()

	router.StaticFS("/static", http.Dir("./playground/static"))
	router.LoadHTMLFiles("./playground/index.html")
	router.POST("/report", report)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Run()
}
