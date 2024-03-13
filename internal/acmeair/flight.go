package acmeair

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/pkg/state"
	"strings"
)

func SegmentName(fromAirport string, toAirport string) string {
	return fmt.Sprintf("%s-%s", fromAirport, toAirport)
}

func GetSegment(ctx context.Context, fromAirport string, toAirport string) FlightSegment {
	segmentName := SegmentName(fromAirport, toAirport)
	// Note: In the original application this is cached (and this is updated very rarely)
	segment, err := state.GetState[FlightSegment](ctx, segmentName)
	if err != nil {
		panic(err)
	}
	return segment
}

func FlightKey(segmentName string, deptDate string) string {
	return fmt.Sprintf("%s-%s", segmentName, deptDate)
}

func GetFlightsBySegmentAndDate(ctx context.Context, segment FlightSegment, deptDate string) []Flight {
	segmentName := segment.FlightName
	// Note: Could we replace this with a get-if (the primary key is the flight id)
	// Note: In the original application this is cached
	flightKey := FlightKey(segmentName, deptDate)
	flightIds, err := state.GetState[[]string](ctx, flightKey)
	if err != nil {
		panic(err)
	}
	flights := GetFlightsByIds(ctx, flightIds)
	// TODO: Restore the segment in flights
	return flights
}

func GetFlightsByIds(ctx context.Context, flightIds []string) []Flight {
	flights := make([]Flight, len(flightIds))
	for i, flightId := range flightIds {
		flight, err := state.GetState[Flight](ctx, flightId)
		if err != nil {
			panic(err)
		}
		flights[i] = flight
	}
	//fmt.Printf("[ReviewStorage] Returning: %v\n", reviews)
	return flights
}

func GetFlightsByAirportsAndDepartureDate(ctx context.Context, fromAirport string, toAirport string, deptDate string) []Flight {
	// Get segment from KV store DB (from, to) -> FlightSegment
	segment := GetSegment(ctx, fromAirport, toAirport)

	// Get flight from DB using segment name and deptDate (segmentName, deptDate) -> []Flights
	flights := GetFlightsBySegmentAndDate(ctx, segment, deptDate)
	return flights
}

func CreateSegment(ctx context.Context, fromAirport string, toAirport string, miles int) string {
	name := SegmentName(fromAirport, toAirport)
	segment := FlightSegment{
		FlightName: name,
		OriginPort: fromAirport,
		DestPort:   toAirport,
		Miles:      miles,
	}
	state.SetState(ctx, name, segment)
	return name
}

func CreateFlight(ctx context.Context,
	segmentName string,
	depTime string,
	arrTime string,
	firstCost int,
	economyCost int,
	firstSeats int,
	economySeats int,
	airplaneType string,
) string {
	depDate := timeToDate(depTime)
	key := FlightKey(segmentName, depDate)
	flight := Flight{
		FlightId:             key,
		SegmentName:          segmentName,
		Segment:              nil,
		DepartureTime:        depTime,
		ArrivalTime:          arrTime,
		FirstClassBaseCost:   firstCost,
		EconomyClassBaseCost: economyCost,
		FirstClassSeats:      firstSeats,
		EconomyClassSeats:    economySeats,
		AirplaneType:         airplaneType,
	}
	return flight
}

// TODO: CreateFlight(segmentName, depTime, arrTime, ...):
//         date = timeToDate(depTime)
//         key = f'{segmentName}-{depDate}'
//         append(key, Flight(id, segmentName, depTime, arrTime, ...)
//         return key

// Assumes time includes date and is in the following form:
// 2022-05-29 15:50
func timeToDate(time string) string {
	// TODO: Extract date from time
	tokens := strings.Split(time, " ")
	return tokens[0]
}
