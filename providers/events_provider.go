package providers

type EventProvider interface {
	Events(lat float64, lon float64, rng int, sorting string) []DojoEvent
}

type DojoEvent struct {
	Title        string
	Description  string
	Logo         string
	Icon         string
	TicketUrl    string
	StartTime    int64
	EndTime      int64
	Capacity     int
	Participants int
	Location     DojoLocation
	Organizer    DojoOrganizer
	Free         bool
	Price        float32
}

type DojoLocation struct {
	Address    string
	City       string
	Country    string
	Name       string
	PostalCode string
	Latitude   float64
	Longitude  float64
}

type DojoOrganizer struct {
	Id       string
	Name     string
	Platform string
}
