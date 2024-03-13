import argparse
import json

iterations = 20

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("--hotels_file",
                        help="the path that contains the hotel file",
                        default="./experiments/hotel/data/hotels.json")
    parser.add_argument("--cities_file",
                        help="the path that contains the analysis file",
                        default="./experiments/hotel/data/cities.json")
    parser.add_argument("--total_hotels",
                        help="The total number of hotels to create",
                        default=100000,
                        type=int)
    args = parser.parse_args()
    return args


def parse_cities_file(cities_file: str):
    with open(cities_file, "r") as f:
        raw_cities = json.load(f)

    cities = []
    total_population = 0
    for i in range(iterations):
        for raw_city in raw_cities:
            city = {}
            city["name"] = raw_city["slug"] + "-" + str(i)
            city["population"] = raw_city["pop2023"]
            total_population += city["population"]
            cities.append(city)
    
    for city in cities:
        city["normalized_population"] = float(city["population"]) / total_population
    
    return cities

## Create hotels proportionally to the city population
def create_hotels(cities, total_hotels):
    hotels = []
    i = 1
    for city in cities:
        city_hotels = int(city["normalized_population"] * total_hotels)
        for city_hotel_i in range(city_hotels):
            hotel = {}
            hotel["id"] = str(i)
            hotel["name"] = f'{city["name"]}-hotel-{city_hotel_i}'
            hotel["phoneNumber"] = f'{city["name"]}-phone-{city_hotel_i}'
            hotel["address"] = {
                "city": city["name"]
            }
            hotels.append(hotel)
            i += 1
    return hotels

def main(args):
    cities = parse_cities_file(args.cities_file)
    hotels = create_hotels(cities, args.total_hotels)

    hotels_str = json.dumps(hotels)
    with open(args.hotels_file, "w") as f:
        f.write(hotels_str)
        f.write("\n")

if __name__ == '__main__':
    args = parse_arguments()
    main(args)
