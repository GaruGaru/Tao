package providers

import (
	"time"
	"net/url"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"sync"
	"fmt"
	"os"
)

type Address struct {
	Address1                         string      `json:"address_1"`
	Address2                         interface{} `json:"address_2"`
	City                             string      `json:"city"`
	Region                           string      `json:"region"`
	PostalCode                       string      `json:"postal_code"`
	Country                          string      `json:"country"`
	Latitude                         string      `json:"latitude"`
	Longitude                        string      `json:"longitude"`
	LocalizedAddressDisplay          string      `json:"localized_address_display"`
	LocalizedAreaDisplay             string      `json:"localized_area_display"`
	LocalizedMultiLineAddressDisplay []string    `json:"localized_multi_line_address_display"`
}

type Venue struct {
	Address     Address `json:"address"`
	ResourceURI string  `json:"resource_uri"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Latitude    string  `json:"latitude"`
	Longitude   string  `json:"longitude"`
}

type Event struct {
	Name struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	} `json:"name"`
	Description struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	} `json:"description"`
	ID  string `json:"id"`
	URL string `json:"url"`
	Start struct {
		Timezone string    `json:"timezone"`
		Local    string    `json:"local"`
		Utc      time.Time `json:"utc"`
	} `json:"start"`
	End struct {
		Timezone string    `json:"timezone"`
		Local    string    `json:"local"`
		Utc      time.Time `json:"utc"`
	} `json:"end"`
	Created           time.Time   `json:"created"`
	Changed           time.Time   `json:"changed"`
	Capacity          int         `json:"capacity"`
	CapacityIsCustom  bool        `json:"capacity_is_custom"`
	Status            string      `json:"status"`
	Currency          string      `json:"currency"`
	Listed            bool        `json:"listed"`
	Shareable         bool        `json:"shareable"`
	OnlineEvent       bool        `json:"online_event"`
	TxTimeLimit       int         `json:"tx_time_limit"`
	HideStartDate     bool        `json:"hide_start_date"`
	HideEndDate       bool        `json:"hide_end_date"`
	Locale            string      `json:"locale"`
	IsLocked          bool        `json:"is_locked"`
	PrivacySetting    string      `json:"privacy_setting"`
	IsSeries          bool        `json:"is_series"`
	IsSeriesParent    bool        `json:"is_series_parent"`
	IsReservedSeating bool        `json:"is_reserved_seating"`
	Source            string      `json:"source"`
	IsFree            bool        `json:"is_free"`
	Version           string      `json:"version"`
	LogoID            string      `json:"logo_id"`
	OrganizerID       string      `json:"organizer_id"`
	VenueID           string      `json:"venue_id"`
	CategoryID        string      `json:"category_id"`
	SubcategoryID     interface{} `json:"subcategory_id"`
	FormatID          string      `json:"format_id"`
	ResourceURI       string      `json:"resource_uri"`
	Logo struct {
		CropMask struct {
			TopLeft struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"top_left"`
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"crop_mask"`
		Original struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"original"`
		ID           string `json:"id"`
		URL          string `json:"url"`
		AspectRatio  string `json:"aspect_ratio"`
		EdgeColor    string `json:"edge_color"`
		EdgeColorSet bool   `json:"edge_color_set"`
	} `json:"logo"`
}
type EventbriteResponse struct {
	Pagination struct {
		ObjectCount  int  `json:"object_count"`
		PageNumber   int  `json:"page_number"`
		PageSize     int  `json:"page_size"`
		PageCount    int  `json:"page_count"`
		HasMoreItems bool `json:"has_more_items"`
	} `json:"pagination"`
	Events [] Event `json:"events"`
	Location struct {
		Latitude  string `json:"latitude"`
		Within    string `json:"within"`
		Longitude string `json:"longitude"`
	} `json:"location"`
}

func NewEventBriteProvider() EventBrite {
	return EventBrite{ApiKey: os.Getenv("EVENTBRITE_TOKEN")}
}

type EventBrite struct {
	ApiKey string
}

func (e EventBrite) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	events, err := e.eventsList(lat, lon, rng, sorting, 1)

	if err != nil {
		return nil, err
	}

	eventsCount := events.Pagination.ObjectCount
	fmt.Printf("Got %d events from eventbrite\n", eventsCount)
	eventsChannel := make(chan DojoEvent, eventsCount)
	var wg sync.WaitGroup
	wg.Add(eventsCount)

	if events.Pagination.HasMoreItems{
		for i := 1; i < events.Pagination.PageCount+1; i++ {
			go fetchAndProcessEvents(e, lat, lon, rng, sorting, i, eventsChannel, &wg)
		}
	}else{
		for _, event := range events.Events {
			go e.processEvent(lat, lon, event, eventsChannel, &wg)
		}
	}


	wg.Wait()
	close(eventsChannel)

	var dojoEvents []DojoEvent

	for result := range eventsChannel {
		dojoEvents = append(dojoEvents, result)
	}

	return dojoEvents, nil
}

func fetchAndProcessEvents(e EventBrite, lat float64, lon float64, rng int, sorting string, i int, eventsChannel chan DojoEvent, wg *sync.WaitGroup) {
		currEvents, err := e.eventsList(lat, lon, rng, sorting, i)
		fmt.Println(fmt.Sprintf("Processing %d events, pagination %d", len(currEvents.Events), i))
		if err == nil {
			for _, event := range currEvents.Events {
				go e.processEvent(lat, lon, event, eventsChannel, wg)
			}
		} else {
			fmt.Printf("Unable to get events")
		}
}

func (e EventBrite) eventsList(lat float64, lon float64, rng int, sorting string, page int) (EventbriteResponse, error) {
	apiUrl := e.eventListUrl(lat, lon, rng, sorting, page)

	resp, err := http.Get(apiUrl.String())

	fmt.Println("Calling " + apiUrl.String())

	if err != nil {
		return EventbriteResponse{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var events EventbriteResponse

	err = json.Unmarshal(body, &events)

	if err != nil {
		return EventbriteResponse{}, err
	}

	return events, nil
}

func (e EventBrite) processEvent(hLat float64, hLon float64, event Event, events chan DojoEvent, group *sync.WaitGroup) {

	defer group.Done()

	venue, err := e.venue(event.VenueID)

	if err == nil {
		events <- toDojoEvent(hLat, hLon, event, venue)
	} else {
		fmt.Printf("Unable to fetch venue for id %d: %s\n", event.VenueID, err)
	}

}

func toDojoEvent(hLat float64, hLon float64, event Event, venue Venue) (DojoEvent) {

	lat, _ := strconv.ParseFloat(venue.Address.Latitude, 64)
	lon, _ := strconv.ParseFloat(venue.Address.Longitude, 64)

	distance := Distance(hLat, hLon, lat, lon) / 1000

	return DojoEvent{
		Title:       event.Name.Text,
		Description: event.Description.Text,
		Logo:        event.Logo.URL,
		TicketUrl:   event.URL,
		Capacity:    event.Capacity,
		StartTime:   event.Start.Utc.Unix(),
		EndTime:     event.Start.Utc.Unix(),
		Location: DojoLocation{
			Address:    venue.Address.Address1,
			City:       venue.Address.City,
			Name:       venue.Name,
			Country:    venue.Address.Country,
			PostalCode: venue.Address.PostalCode,
			Latitude:   lat,
			Longitude:  lon,
			Distance:   distance,
		},
		Organizer: DojoOrganizer{Id: event.OrganizerID},
		Free:      event.IsFree,
	}
}

func (e EventBrite) venue(venueID string) (Venue, error) {
	venueApiUrl := e.venueUrl(venueID)
	resp, err := http.Get(venueApiUrl.String())

	if err != nil {
		return Venue{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var venue Venue

	err = json.Unmarshal(body, &venue)

	if err != nil {
		return Venue{}, err
	}

	return venue, nil

}

func (e EventBrite) eventListUrl(lat float64, lon float64, rng int, sorting string, page int) url.URL {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "www.eventbriteapi.com",
		Path:   "/v3/events/search/",
	}

	query := apiUrl.Query()
	query.Set("q", "coderdojo")
	query.Set("token", e.ApiKey)
	query.Set("location.latitude", strconv.FormatFloat(lat, 'f', 8, 64))
	query.Set("location.longitude", strconv.FormatFloat(lon, 'f', 8, 64))
	query.Set("location.within", strconv.Itoa(rng)+"km")
	query.Set("sort_by", sorting)
	query.Set("price", "free")
	query.Set("page", strconv.Itoa(page))

	apiUrl.RawQuery = query.Encode()

	return *apiUrl
}

func (e EventBrite) venueUrl(venueID string) url.URL {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "www.eventbriteapi.com",
		Path:   fmt.Sprintf("/v3/venues/%s/", venueID),
	}
	query := apiUrl.Query()
	query.Set("token", e.ApiKey)
	apiUrl.RawQuery = query.Encode()
	return *apiUrl
}
