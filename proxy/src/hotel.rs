use crate::util;
use async_trait::async_trait;
use clap::Parser;
use hyper::{Body, Request, Response};
use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};
use std::convert::Infallible;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::str::FromStr;
use tracing::info;

static FRONTEND_IP: OnceCell<String> = OnceCell::new();

static MAX_USER: OnceCell<usize> = OnceCell::new();

static TOTAL_HOTELS: OnceCell<usize> = OnceCell::new();
static TOTAL_CITIES: OnceCell<usize> = OnceCell::new();

static POPULARITIES: OnceCell<Vec<f64>> = OnceCell::new();
static CITIES: OnceCell<Vec<String>> = OnceCell::new();

#[derive(Parser)]
pub struct Hotel {
    #[clap(long, value_parser)]
    frontend: String,
}

// curl -X POST -H "Content-Type: application/json" -d '{"movie_id":"299534"}' $(kip page)/ro_read_page

#[async_trait]
impl util::Backend for Hotel {
    async fn prepare(&self) {
        FRONTEND_IP.set(self.frontend.clone()).unwrap();
        let username = whoami::username();
        let analysis_file =
            format!("/users/{username}/mucache/experiments/hotel/data/analysis.txt");
        info!("Reading analysis data from {}", analysis_file);
        read_analysis_file(&analysis_file);
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        service_hotel(_req).await
    }
}

pub fn read_analysis_file(path: &str) {
    let file = File::open(path).unwrap();
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().map(|l| l.unwrap()).collect();

    let max_user = usize::from_str(lines[0].split(':').nth(1).unwrap().trim()).unwrap();
    let total_hotels = usize::from_str(lines[1].split(':').nth(1).unwrap().trim()).unwrap();
    let total_cities = usize::from_str(lines[2].split(':').nth(1).unwrap().trim()).unwrap();

    let mut popularities = vec![0.0; total_cities];
    let mut cities = vec!["".to_string(); total_cities];

    let mut index = 0;
    for line in lines[3..].iter() {
        let popularity = f64::from_str(line.split('|').nth(1).unwrap().trim()).unwrap();
        let city = line.split('|').nth(0).unwrap().trim();
        popularities[index] = popularity;
        cities[index] = city.to_string();
        index += 1;
    }
    MAX_USER.set(max_user).unwrap();
    TOTAL_HOTELS.set(total_hotels).unwrap();
    TOTAL_CITIES.set(total_cities).unwrap();
    POPULARITIES.set(popularities).unwrap();
    CITIES.set(cities).unwrap();
}

pub fn random_user() -> usize {
    let max_user = MAX_USER.get().unwrap();
    rand::random::<usize>() % max_user
}

pub fn random_hotel() -> usize {
    let total_hotel = TOTAL_HOTELS.get().unwrap();
    rand::random::<usize>() % total_hotel
}

// TODO: We could popularities to query non-uniformly
pub fn random_location() -> usize {
    let max_loc = TOTAL_CITIES.get().unwrap();
    rand::random::<usize>() % max_loc
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Reservation {
    pub in_date: String,
    pub out_date: String,
    pub hotel_id: String,
    pub username: String,
    pub password: String,
    pub rooms: u32,
}

pub fn reservation() -> Reservation {
    let user_idx = random_user();
    let user_str = "username".to_owned() + &user_idx.to_string();
    let pass_str = "password".to_owned() + &user_idx.to_string();
    let hotel_idx = random_hotel();
    let hotel_id_str = hotel_idx.to_string();
    Reservation {
        in_date: "2023-04-17".to_string(),
        out_date: "2023-04-19".to_string(),
        hotel_id: hotel_id_str,
        username: user_str,
        password: pass_str,
        rooms: 1,
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Search {
    pub in_date: String,
    pub out_date: String,
    pub location: String,
}

pub fn search() -> Search {
    let location_idx = random_location();
    let city = &CITIES.get().unwrap()[location_idx];
    Search {
        in_date: "2023-04-17".to_string(),
        out_date: "2023-04-19".to_string(),
        location: city.to_string(),
    }
}

pub async fn service_reservation(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let c = reservation();
    // info!("{:?}", c);
    let req = serde_json::to_string(&c).unwrap();
    let _res = util::send_req(FRONTEND_IP.get().unwrap(), "reservation", req).await;
    // info!("{:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_search(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let c = search();
    let req = serde_json::to_string(&c).unwrap();
    // info!("{:?}", c);
    let _res = util::send_req(FRONTEND_IP.get().unwrap(), "ro_search_hotels", req).await;
    // info!("{:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_hotel(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let search_hotels_ratio = 80;
    let reservation_ratio = 20;
    assert_eq!(search_hotels_ratio + reservation_ratio, 100);
    let coin = rand::random::<usize>() % 100;
    if coin < search_hotels_ratio {
        service_search(_req).await
    } else {
        service_reservation(_req).await
    }
}
