package hotel

type HotelProfile struct {
	HotelId string `json:"hotel_id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Info    string `json:"info"`
	// Note: Normally there would be more fields here, images, etc
}

type Rate struct {
	HotelId string `json:"hotelid"`
	Price   int    `json:"price"`
}

type Reservation struct {
	CustomerName string `json:"customer_name"`
	InDate       string `json:"in_date"`
	OutDate      string `json:"out_date"`
	RoomNumber   int    `json:"room_number"`
}

type HotelAvailability struct {
	Capacity     int           `json:"capacity"`
	Reservations []Reservation `json:"reservations"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Password []byte `json:"password"`
	Salt     string `json:"salt"`
}

// Frontend
type SearchHotelsRequest struct {
	InDate   string `json:"in_date"`
	OutDate  string `json:"out_date"`
	Location string `json:"location"`
}

type SearchHotelsResponse struct {
	Profiles []HotelProfile `json:"profiles"`
}

type StoreHotelRequest struct {
	HotelId  string `json:"hotel_id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Rate     int    `json:"rate"`
	Capacity int    `json:"capacity"`
	Info     string `json:"info"`
}

type StoreHotelResponse struct {
	HotelId string `json:"hotel_id"`
}

type FrontendReservationRequest struct {
	HotelId  string `json:"hotel_id"`
	InDate   string `json:"in_date"`
	OutDate  string `json:"out_date"`
	Rooms    int    `json:"rooms"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FrontendReservationResponse struct {
	Success bool `json:"success"`
}

// Search
type NearbyRequest struct {
	InDate   string `json:"inDate"`
	OutDate  string `json:"outDate"`
	Location string `json:"location"`
}

type NearbyResponse struct {
	Rates []Rate `json:"rates"`
}

type StoreHotelLocationRequest struct {
	HotelId  string `json:"hotel_id"`
	Location string `json:"location"`
}

type StoreHotelLocationResponse struct {
	HotelId string `json:"hotel_id"`
}

// Rate
type GetRatesRequest struct {
	HotelIds []string `json:"hotel_ids"`
}

type GetRatesResponse struct {
	Rates []Rate `json:"rates"`
}

type StoreRateRequest struct {
	Rate Rate `json:"rate"`
}

type StoreRateResponse struct {
	HotelId string `json:"hotel_id"`
}

// Profile
type GetProfilesRequest struct {
	HotelIds []string `json:"hotel_ids"`
}

type GetProfilesResponse struct {
	Profiles []HotelProfile `json:"profiles"`
}

type StoreProfileRequest struct {
	Profile HotelProfile `json:"profile"`
}

type StoreProfileResponse struct {
	HotelId string `json:"hotel_id"`
}

// Reservation
type CheckAvailabilityRequest struct {
	CustomerName string   `json:"customer_name"`
	HotelIds     []string `json:"hotel_ids"`
	InDate       string   `json:"in_date"`
	OutDate      string   `json:"out_date"`
	RoomNumber   int      `json:"room_number"`
}

type CheckAvailabilityResponse struct {
	HotelIds []string `json:"hotel_ids"`
}

type MakeReservationRequest struct {
	CustomerName string `json:"customer_name"`
	HotelId      string `json:"hotel_id"`
	InDate       string `json:"in_date"`
	OutDate      string `json:"out_date"`
	RoomNumber   int    `json:"room_number"`
}

type MakeReservationResponse struct {
	Success bool `json:"success"`
}

type AddHotelAvailabilityRequest struct {
	HotelId  string `json:"hotel_id"`
	Capacity int    `json:"capacity"`
}

type AddHotelAvailabilityResponse struct {
	Hotelid string `json:"hotelid"`
}

// user
type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Ok bool `json:"ok"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
