package main

import (
	"github.com/GaruGaru/Tao/providers"
	"os"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {

	provider := providers.EventBrite{
		ApiKey: os.Getenv("EVENTBRITE_TOKEN"),
	}

	r := gin.Default()
	r.GET("/probe", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/api/v1/events", func(c *gin.Context) {
		lat, e := strconv.ParseFloat(c.Query("lat"), 64)
		if e != nil {
			c.String(400, "Invalid parameter lat")
			return
		}
		lon, e := strconv.ParseFloat(c.Query("lon"), 64)
		if e != nil {
			c.String(400, "Invalid parameter lon")
			return
		}
		rng, e := strconv.Atoi(c.Query("range"))
		if e != nil {
			c.String(400, "Invalid parameter range")
			return
		}
		sorting := c.Query("sort_by")
		if len(sorting) == 0 {
			sorting = "distance"
		}
		events, err := provider.Events(lat, lon, rng, sorting)

		if err == nil {
			c.JSON(200, events)
		} else {
			c.Error(err)
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080

}

