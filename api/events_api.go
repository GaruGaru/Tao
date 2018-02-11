package api

import (
	"github.com/GaruGaru/Tao/providers"
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/cactus/go-statsd-client/statsd"
)

type EventsApi struct {
	Provider providers.EventProvider
	Statsd   statsd.Statter
}

func (api EventsApi) Run() {

	r := gin.Default()

	r.GET("/probe", probe)

	r.GET("/api/v1/events", api.eventsV1)

	r.Run() // listen and serve on 0.0.0.0:8080

}

func probe(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (api EventsApi) eventsV1(c *gin.Context) {

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

	events, err := api.Provider.Events(lat, lon, rng, sorting)

	if err == nil {
		//api.Statsd.Inc("request.eventsv1.ok", 1, 1)
		c.JSON(200, events)
	} else {
		//api.Statsd.Inc("request.eventsv1.fail", 1, 1)
		c.Error(err)
	}

}
