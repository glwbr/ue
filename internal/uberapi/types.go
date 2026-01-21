package uberapi

type UserResponse struct {
	Data struct {
		CurrentUser *CurrentUser `json:"currentUser"`
	} `json:"data"`
}

type CurrentUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type ActivitiesResponse struct {
	Data struct {
		Activities struct {
			Past struct {
				Activities    []Activity `json:"activities"`
				NextPageToken string     `json:"nextPageToken"`
				TypeName      string     `json:"__typename"`
			} `json:"past"`
			TypeName string `json:"__typename"`
		} `json:"activities"`
	} `json:"data"`
}

type Activity struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Subtitle    string `json:"subtitle"`
	TypeName    string `json:"__typename"`
}

type GetTripResponse struct {
	Data struct {
		GetTrip struct {
			Trip struct {
				BeginTripTime      string   `json:"beginTripTime"`
				DropoffTime        string   `json:"dropoffTime"`
				CityID             int      `json:"cityID"`
				CountryID          int      `json:"countryID"`
				Status             string   `json:"status"`
				Fare               string   `json:"fare"`
				Driver             string   `json:"driver"`
				UUID               string   `json:"uuid"`
				VehicleDisplayName string   `json:"vehicleDisplayName"`
				Waypoints          []string `json:"waypoints"`
				Marketplace        string   `json:"marketplace"`
				TypeName           string   `json:"__typename"`
			} `json:"trip"`
			MapURL   string  `json:"mapURL"`
			Receipt  Receipt `json:"receipt"`
			Rating   string  `json:"rating"`
			TypeName string  `json:"__typename"`
		} `json:"getTrip"`
	} `json:"data"`
}

type Receipt struct {
	Distance      string `json:"distance"`
	DistanceLabel string `json:"distanceLabel"`
	Duration      string `json:"duration"`
	VehicleType   string `json:"vehicleType"`
	TypeName      string `json:"__typename"`
}
