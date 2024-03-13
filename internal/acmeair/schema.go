package acmeair

type FlightSegment struct {
	// The identifier of the flight
	FlightName string `json:"flight_name"`
	// The origin airport
	OriginPort string `json:"origin_port"`
	// The destination airport
	DestPort string `json:"dest_port"`
	Miles    int    `json:"miles"`
}

type Flight struct {
	FlightId    string `json:"flight_id"`
	SegmentName string `json:"segment_name"`
	// The segment is not stored in the DB (to avoid redundancy)
	Segment       *FlightSegment `json:"segment"`
	DepartureTime string         `json:"departure_time"`
	ArrivalTime   string         `json:"arrival_time"`
	// We store costs in ints (i.e., cents) for precision
	FirstClassBaseCost   int `json:"first_class_base_cost"`
	EconomyClassBaseCost int `json:"economy_class_base_cost"`
	// TODO: Maybe availability and costs need to be accessed from a different booking service
	FirstClassSeats   int    `json:"first_class_seats"`
	EconomyClassSeats int    `json:"economy_class_seats"`
	AirplaneType      string `json:"airplane_type"`
}

// flight

type GetFlightsRequest struct {
	FromAirport string `json:"from_airport"`
	ToAirport   string `json:"to_airport"`
	DeptDate    string `json:"dept_date"`
}

type GetFlightsResponse struct {
	Flights []Flight `json:"flights"`
}

type CreateSegmentRequest struct {
	OriginPort string `json:"origin_port"`
	DestPort   string `json:"dest_port"`
	Miles      int    `json:"miles"`
}

type CreateSegmentResponse struct {
	FlightName string `json:"flight_name"`
}

type CreateFlightRequest struct {
	SegmentId            string `json:"segment_id"`
	DepartureTime        string `json:"departure_time"`
	ArrivalTime          string `json:"arrival_time"`
	FirstClassBaseCost   int    `json:"first_class_base_cost"`
	EconomyClassBaseCost int    `json:"economy_class_base_cost"`
	FirstClassSeats      int    `json:"first_class_seats"`
	EconomyClassSeats    int    `json:"economy_class_seats"`
	AirplaneType         string `json:"airplane_type"`
}

type CreateFlightResponse struct {
	FlightKey string `json:"flight_key"`
}
