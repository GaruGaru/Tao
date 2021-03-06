package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

type ZenDojo struct {
	ID                   string      `json:"id"`
	MysqlDojoID          interface{} `json:"mysql_dojo_id"`
	DojoLeadID           string      `json:"dojo_lead_id"`
	Name                 string      `json:"name"`
	Creator              string      `json:"creator"`
	Created              time.Time   `json:"created"`
	VerifiedAt           time.Time   `json:"verified_at"`
	VerifiedBy           string      `json:"verified_by"`
	Verified             int         `json:"verified"`
	NeedMentors          int         `json:"need_mentors"`
	Stage                int         `json:"stage"`
	MailingList          int         `json:"mailing_list"`
	AlternativeFrequency string      `json:"alternative_frequency"`
	Country              struct {
		CountryName   string      `json:"countryName"`
		CountryNumber interface{} `json:"countryNumber,omitempty"`
		Continent     string      `json:"continent"`
		Alpha2        string      `json:"alpha2"`
		Alpha3        string      `json:"alpha3"`
	} `json:"country"`
	County struct {
	} `json:"county"`
	State struct {
	} `json:"state"`
	City struct {
	} `json:"city"`
	Place struct {
		NameWithHierarchy string `json:"nameWithHierarchy"`
	} `json:"place"`
	Coordinates string `json:"coordinates"`
	GeoPoint    struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"geo_point"`
	Notes                      string        `json:"notes"`
	Email                      string        `json:"email"`
	Website                    string        `json:"website"`
	Twitter                    string        `json:"twitter"`
	GoogleGroup                string        `json:"google_group"`
	SupporterImage             string        `json:"supporter_image"`
	Deleted                    int           `json:"deleted"`
	DeletedBy                  interface{}   `json:"deleted_by"`
	DeletedAt                  interface{}   `json:"deleted_at"`
	Private                    int           `json:"private"`
	URLSlug                    string        `json:"url_slug"`
	Continent                  string        `json:"continent"`
	Alpha2                     string        `json:"alpha2"`
	Alpha3                     string        `json:"alpha3"`
	Address1                   string        `json:"address1"`
	Address2                   string        `json:"address2"`
	CountryNumber              interface{}   `json:"country_number,omitempty"`
	CountryName                string        `json:"country_name"`
	Admin1Code                 interface{}   `json:"admin1_code"`
	Admin1Name                 interface{}   `json:"admin1_name"`
	Admin2Code                 interface{}   `json:"admin2_code"`
	Admin2Name                 interface{}   `json:"admin2_name"`
	Admin3Code                 interface{}   `json:"admin3_code"`
	Admin3Name                 interface{}   `json:"admin3_name"`
	Admin4Code                 interface{}   `json:"admin4_code"`
	Admin4Name                 interface{}   `json:"admin4_name"`
	PlaceGeonameID             interface{}   `json:"place_geoname_id"`
	PlaceName                  string        `json:"place_name"`
	UserInvites                []interface{} `json:"user_invites"`
	CreatorEmail               string        `json:"creator_email"`
	TaoVerified                int           `json:"tao_verified"`
	ExpectedAttendees          int           `json:"expected_attendees"`
	Facebook                   interface{}   `json:"facebook"`
	Day                        interface{}   `json:"day"`
	StartTime                  interface{}   `json:"start_time"`
	EndTime                    interface{}   `json:"end_time"`
	Frequency                  string        `json:"frequency"`
	DistanceFromSearchLocation float64       `json:"distance_from_search_location"`
}

type ZenDojoEvent struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country struct {
		CountryName string `json:"countryName"`
		Alpha2      string `json:"alpha2"`
	} `json:"country"`
	City struct {
		NameWithHierarchy string `json:"nameWithHierarchy"`
	} `json:"city"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   string    `json:"createdBy"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	DojoID      string    `json:"dojoId"`
	Position    struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"position"`
	Public        bool   `json:"public"`
	Status        string `json:"status"`
	RecurringType string `json:"recurringType"`
	Dates         []struct {
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	} `json:"dates"`
	TicketApproval    bool        `json:"ticketApproval"`
	NotifyOnApplicant bool        `json:"notifyOnApplicant"`
	EventbriteID      interface{} `json:"eventbriteId"`
	EventbriteURL     interface{} `json:"eventbriteUrl"`
	UseDojoAddress    interface{} `json:"useDojoAddress"`
	StartTime         time.Time   `json:"startTime"`
	EndTime           time.Time   `json:"endTime"`
}

type ZenDojoEvents struct {
	Results []ZenDojoEvent `json:"results"`
	Total   int            `json:"total"`
}

type EventsSearchRequest struct {
	Query interface{} `json:"query"`
}
type BoundingBoxRequest struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Rng int     `json:"radius"`
}

type ZenPlatformProvider struct {
	Client ApiClient
}

func NewZenPlatformProvider() ZenPlatformProvider {
	return ZenPlatformProvider{
		Client: NewHttpApiClient(20 * time.Second),
	}
}

func (z ZenPlatformProvider) Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error) {

	dojos, err := z.fetchDojos(lat, lon, rng)

	log.Infof("Got %d dojos from zen", len(dojos))

	if err != nil {
		return nil, err
	}

	events, err := z.fetchEventsFromDojos(dojos)

	if err != nil {
		return nil, err
	}

	log.Infof("Done zen platform provider, total events: %d", len(events))

	return events, nil
}

func (z ZenPlatformProvider) fetchEventsFromDojos(dojos []ZenDojo) ([]DojoEvent, error) {

	dojosCount := len(dojos)
	dojosChannel := make(chan []DojoEvent, dojosCount)
	var wg sync.WaitGroup

	for _, dojo := range dojos {
		wg.Add(1)
		z.fetchEventsFromZenDojo(dojo, dojosChannel, &wg)
	}

	wg.Wait()
	log.Info("Done waiting")
	close(dojosChannel)

	var dojoEvents []DojoEvent

	for result := range dojosChannel {
		for _, event := range result {
			dojoEvents = append(dojoEvents, event)
		}
	}

	return dojoEvents, nil

}

func (z ZenPlatformProvider) fetchEventsFromZenDojo(dojo ZenDojo, eventsChannel chan []DojoEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	log.WithField("coderdojo", dojo.ID).Debug("Fetching zen events for dojo")
	rawUrl := fmt.Sprintf("https://zen.coderdojo.com/api/3.0/dojos/%s/events?query[status]=published&query[public]=1&query[afterDate]=%d&query[utcOffset]=0", dojo.ID, 0)
	url, err := url.Parse(rawUrl)

	if err != nil {
		logrus.Errorf("Unable to parse url: %s", rawUrl)
		return
	}

	var events ZenDojoEvents

	err = z.Client.Get(url.String(), &events)

	if err != nil {
		log.Error(err.Error())
		return
	}

	zenEvents := make([]DojoEvent, 0)
	for _, event := range events.Results {
		if event.StartTime.After(time.Now()) {
			zenEvents = append(zenEvents, zenToDojoEvent(dojo, event))
		}
	}

	if len(zenEvents) != 0 {
		log.Infof("Got %d events from zen coderdojo %s", len(zenEvents), dojo.ID)
	}

	eventsChannel <- zenEvents

}

func zenTicketUrl(dojo ZenDojo, event ZenDojoEvent) string {
	eventbriteUrl, ok := event.EventbriteURL.(string)

	if ok {
		return eventbriteUrl
	}

	if dojo.URLSlug != "" {
		return fmt.Sprintf("https://zen.coderdojo.com/dojos/%s", dojo.URLSlug)
	}

	if dojo.Email != "" {
		return fmt.Sprintf("mailto://%s", dojo.Email)
	}

	return ""

}

func zenToDojoEvent(dojo ZenDojo, event ZenDojoEvent) DojoEvent {
	return DojoEvent{
		Title:        event.Name,
		Description:  event.Description,
		Logo:         dojo.SupporterImage,
		Icon:         "",
		TicketUrl:    zenTicketUrl(dojo, event),
		StartTime:    event.StartTime.Unix(),
		EndTime:      event.EndTime.Unix(),
		Capacity:     0,
		Participants: 0,
		Location: DojoLocation{
			Latitude:  event.Position.Lat,
			Longitude: event.Position.Lng,
			Address:   event.Address,
			City:      event.City.NameWithHierarchy,
			Country:   event.Country.CountryName,
			Distance:  dojo.DistanceFromSearchLocation,
		},
		Organizer: DojoOrganizer{
			Name:     dojo.Name,
			Id:       dojo.ID,
			Platform: "zen",
		},
		Free:  true,
		Price: 0,
	}
}

func (z ZenPlatformProvider) fetchDojos(lat float64, lon float64, rng int) ([]ZenDojo, error) {
	url := "https://zen.coderdojo.com/api/2.0/dojos/search-bounding-box"
	geoRequest := EventsSearchRequest{
		Query: BoundingBoxRequest{
			Lat: lat,
			Lon: lon,
			Rng: rng,
		},
	}

	jbytes, err := json.Marshal(geoRequest)

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"url":  url,
		"body": string(jbytes),
	}).Info("Fetching zen nearby dojos")

	var dojos []ZenDojo

	err = z.Client.Post(url, bytes.NewBuffer(jbytes), &dojos)

	return dojos, err
}
