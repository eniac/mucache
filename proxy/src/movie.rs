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

// The write path/request
static FRONTEND_IP: OnceCell<String> = OnceCell::new();
// The read path/request
static PAGE_IP: OnceCell<String> = OnceCell::new();

static MAX_USER: OnceCell<usize> = OnceCell::new();

static TOTAL_MOVIES: OnceCell<usize> = OnceCell::new();

static POPULARITIES: OnceCell<Vec<f64>> = OnceCell::new();
static TITLES: OnceCell<Vec<String>> = OnceCell::new();

#[derive(Parser)]
pub struct Movie {
    #[clap(long, value_parser)]
    frontend: String,
    #[clap(long, value_parser)]
    page: String,
}

// curl -X POST -H "Content-Type: application/json" -d '{"movie_id":"299534"}' $(kip page)/ro_read_page

#[async_trait]
impl util::Backend for Movie {
    async fn prepare(&self) {
        FRONTEND_IP.set(self.frontend.clone()).unwrap();
        PAGE_IP.set(self.page.clone()).unwrap();
        let username = whoami::username();
        let analysis_file =
            format!("/users/{username}/mucache/experiments/movie/data/analysis.txt");
        info!("Reading movie analysis data from {}", analysis_file);
        read_analysis_file(&analysis_file);
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        service_movie(_req).await
    }
}

pub fn read_analysis_file(path: &str) {
    let file = File::open(path).unwrap();
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().map(|l| l.unwrap()).collect();

    let max_user = usize::from_str(lines[0].split(':').nth(1).unwrap().trim()).unwrap();
    let total_movies = usize::from_str(lines[1].split(':').nth(1).unwrap().trim()).unwrap();
    // let total_cast = usize::from_str(lines[2].split(':').nth(1).unwrap().trim()).unwrap();


    let mut popularities = vec![0.0; total_movies];
    let mut titles = vec!["".to_string(); total_movies];

    let mut index = 0;
    for line in lines[4..].iter() {
        let popularity = f64::from_str(line.split('|').nth(0).unwrap().trim()).unwrap();
        let title = line.split('|').nth(1).unwrap().trim();
        popularities[index] = popularity;
        titles[index] = title.to_string();
        index += 1;
    }
    assert_eq!(index, total_movies);
    MAX_USER.set(max_user).unwrap();
    TOTAL_MOVIES.set(total_movies).unwrap();
    POPULARITIES.set(popularities).unwrap();
    TITLES.set(titles).unwrap();
}

pub fn random_user() -> usize {
    let max_user = MAX_USER.get().unwrap();
    rand::random::<usize>() % max_user
}

// TODO: We could popularities to query non-uniformly
pub fn random_movie() -> usize {
    let total_movies = TOTAL_MOVIES.get().unwrap();
    rand::random::<usize>() % total_movies
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Compose {
    pub username: String,
    pub password: String,
    pub title: String,
    pub text: String,
}

pub fn compose_post() -> Compose {
    let user_idx = random_user();
    let movie_idx = random_movie();
    let user_str = "username".to_owned() + &user_idx.to_string();
    let pass_str = "password".to_owned() + &user_idx.to_string();
    let title = &TITLES.get().unwrap()[movie_idx];
    let text = uuid::Uuid::new_v4().to_string();
    Compose {
        username: user_str,
        password: pass_str,
        title: title.to_string(),
        text: text,
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Page {
    pub movie_id: String,
}

pub fn page() -> Page {
    let movie_idx = random_movie();
    let movie_id_str = movie_idx.to_string();
    Page {
        movie_id: movie_id_str,
    }
}

// TODO: Complete
pub async fn service_frontend(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let c = compose_post();
    let req = serde_json::to_string(&c).unwrap();
    let _res = util::send_req(FRONTEND_IP.get().unwrap(), "compose", req).await;
    // info!("{:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_page(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let c = page();
    let req = serde_json::to_string(&c).unwrap();
    let _res = util::send_req(PAGE_IP.get().unwrap(), "ro_read_page", req).await;
    // info!("{:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_movie(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let read_page_ratio = 90;
    let write_frontend_ratio = 10;
    assert_eq!(read_page_ratio + write_frontend_ratio, 100);
    let coin = rand::random::<usize>() % 100;
    if coin < read_page_ratio {
        service_page(_req).await
    } else {
        service_frontend(_req).await
    }
}
