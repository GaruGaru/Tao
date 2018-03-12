package providers

type EventProvider interface {
	Events(lat float64, lon float64, rng int, sorting string) ([]DojoEvent, error)
}

type DojoEvent struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Logo         string        `json:"logo"`
	Icon         string        `json:"icon"`
	TicketUrl    string        `json:"ticketurl"`
	StartTime    int64         `json:"starttime"`
	EndTime      int64         `json:"endtime"`
	Capacity     int           `json:"capacity"`
	Participants int           `json:"participants"`
	Location     DojoLocation  `json:"location"`
	Organizer    DojoOrganizer `json:"organizer"`
	Free         bool          `json:"free"`
	Price        float32       `json:"price"`
}

type DojoLocation struct {
	Address    string  `json:"address"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Name       string  `json:"name"`
	PostalCode string  `json:"postalcode"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Distance   float64 `json:"distance"`
}

type DojoOrganizer struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}
