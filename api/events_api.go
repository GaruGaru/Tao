package api

import (
	"github.com/GaruGaru/Tao/providers"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
)

type EventsApi struct {
	Provider      providers.EventProvider
	RedisProvider providers.RedisEventsProvider
	Statsd        statsd.Statter
}

type EventsRequest struct {
	latitude    float64
	longitude   float64
	searchRange int
	sorting     string
}

func createRequest(c gin.Context) (EventsRequest, error) {

	lat, e := strconv.ParseFloat(c.Query("lat"), 64)
	if e != nil {
		return EventsRequest{}, fmt.Errorf("invalid parameter lat")
	}

	lon, e := strconv.ParseFloat(c.Query("lon"), 64)
	if e != nil {
		return EventsRequest{}, fmt.Errorf("invalid parameter lon")
	}

	rng, e := strconv.Atoi(c.Query("range"))
	if e != nil {
		return EventsRequest{}, fmt.Errorf("invalid parameter range")
	}

	sorting := c.Query("sort_by")

	if sorting == "" {
		sorting = "distance"
	}

	return EventsRequest{
		latitude:    lat,
		longitude:   lon,
		searchRange: rng,
		sorting:     sorting,
	}, nil

}

func (api EventsApi) Run(port int) {

	r := gin.New()

	r.Use(cors.Default())
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/probe", probe)

	r.GET("/api/v1/events", api.eventsV1)

	r.GET("/api/v2/events", api.eventsV2)

	r.GET("/api/v3/events", api.eventsV3)

	r.Run(fmt.Sprintf(":%d", port))

}

func probe(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func (api EventsApi) eventsV1(c *gin.Context) {

	req, err := createRequest(*c)

	if err != nil {
		c.String(400, err.Error())
	}

	events, err := api.Provider.Events(req.latitude, req.longitude, req.searchRange, req.sorting)

	if err == nil {
		api.Statsd.Inc("request.eventsv1.ok", 1, 1)
		c.JSON(200, events)
	} else {
		api.Statsd.Inc("request.eventsv1.fail", 1, 1)
		c.Error(err)
	}

}

func (api EventsApi) eventsV2(c *gin.Context) {

	req, err := createRequest(*c)

	if err != nil {
		c.String(400, err.Error())
	}

	events, err := api.Provider.Events(req.latitude, req.longitude, req.searchRange, req.sorting)

	if err == nil {
		api.Statsd.Inc("request.eventsv2.ok", 1, 1)
		response := providers.DojoEventResponse{
			Count:  len(events),
			Events: events,
		}
		c.JSON(200, response)
	} else {
		api.Statsd.Inc("request.eventsv2.fail", 1, 1)
		c.Error(err)
	}

}

func (api EventsApi) eventsV3(c *gin.Context) {

	req, err := createRequest(*c)

	if err != nil {
		c.String(400, err.Error())
	}

	events, err := api.RedisProvider.Events(req.latitude, req.longitude, req.searchRange, req.sorting)

	if err == nil {
		api.Statsd.Inc("request.eventsv3.ok", 1, 1)
		response := providers.DojoEventResponse{
			Count:  len(events),
			Events: events,
		}
		c.JSON(200, response)
	} else {
		api.Statsd.Inc("request.eventsv3.fail", 1, 1)
		c.Error(err)
	}

}
