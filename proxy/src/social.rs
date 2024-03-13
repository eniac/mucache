use crate::util;
use async_trait::async_trait;
use clap::Parser;
use hyper::{Body, Request, Response};
use once_cell::sync::OnceCell;
use rand::distributions::{Alphanumeric, DistString};
use rustc_hash::FxHasher;
use serde::{Deserialize, Serialize};
use std::convert::Infallible;
use std::fs::File;
use std::hash::Hasher;
use std::io::{BufRead, BufReader};
use std::str::FromStr;
use tracing::info;

static SHARD: OnceCell<usize> = OnceCell::new();
static COMPOSE_POST_IP: OnceCell<Vec<String>> = OnceCell::new();
static HOME_TIMELINE_IP: OnceCell<Vec<String>> = OnceCell::new();
static USER_TIMELINE_IP: OnceCell<Vec<String>> = OnceCell::new();

static FOLLOWERS: OnceCell<Vec<u32>> = OnceCell::new();
static MAX_USER: OnceCell<usize> = OnceCell::new();
static TOTAL_FOLLOWERS: OnceCell<usize> = OnceCell::new();
static STARTING_POST_SIZE: OnceCell<usize> = OnceCell::new();
static NUMBER_STARTING_POSTS_PER_USER: OnceCell<usize> = OnceCell::new();
static POST_SIZE: usize = 100;
static SAMPLE_POST: OnceCell<String> = OnceCell::new();

#[derive(Parser)]
pub struct Social {
    #[clap(long, value_parser, multiple = true)]
    compose_post: Vec<String>,
    #[clap(long, value_parser, multiple = true)]
    home_timeline: Vec<String>,
    #[clap(long, value_parser, multiple = true)]
    user_timeline: Vec<String>,
}

#[async_trait]
impl util::Backend for Social {
    async fn prepare(&self) {
        COMPOSE_POST_IP.set(self.compose_post.clone()).unwrap();
        USER_TIMELINE_IP.set(self.user_timeline.clone()).unwrap();
        HOME_TIMELINE_IP.set(self.home_timeline.clone()).unwrap();
        SHARD.set(self.compose_post.len()).unwrap();
        let username = whoami::username();
        let social_file =
            format!("/users/{username}/mucache/experiments/social/socfb/socfb-analysis.txt");
        info!("Reading social graph from {}", social_file);
        read_social_graph_analysis_file(&social_file);
        SAMPLE_POST
            .set(Alphanumeric.sample_string(&mut rand::thread_rng(), POST_SIZE))
            .unwrap();
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        service_social(_req).await
    }
}

pub fn read_social_graph_analysis_file(path: &str) {
    let file = File::open(path).unwrap();
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().map(|l| l.unwrap()).collect();

    let max_user = usize::from_str(lines[0].split(':').nth(1).unwrap().trim()).unwrap();
    let total_followers = usize::from_str(lines[1].split(':').nth(1).unwrap().trim()).unwrap();
    let starting_post_size = usize::from_str(lines[2].split(':').nth(1).unwrap().trim()).unwrap();
    let number_starting_posts_per_user =
        usize::from_str(lines[3].split(':').nth(1).unwrap().trim()).unwrap();

    let mut followers = vec![0; max_user + 1];

    for line in lines[5..].iter() {
        let (user_id, n_followers) = {
            let mut items = line.split_whitespace();
            (items.next().unwrap(), items.next().unwrap())
        };
        followers[usize::from_str(user_id).unwrap()] = u32::from_str(n_followers).unwrap();
    }
    MAX_USER.set(max_user).unwrap();
    TOTAL_FOLLOWERS.set(total_followers).unwrap();
    STARTING_POST_SIZE.set(starting_post_size).unwrap();
    NUMBER_STARTING_POSTS_PER_USER
        .set(number_starting_posts_per_user)
        .unwrap();
    FOLLOWERS.set(followers).unwrap();
}

// Read the following link to use weigthed indexing:
//   https://docs.rs/rand/latest/rand/distributions/weighted/struct.WeightedIndex.html
// Or even better the following (0(1) sampling):
//   https://docs.rs/rand_distr/0.4.3/rand_distr/weighted_alias/struct.WeightedAliasIndex.html
pub fn random_user() -> usize {
    let max_user = MAX_USER.get().unwrap();
    rand::random::<usize>() % max_user + 1
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ComposePost {
    pub creator_id: String,
    pub text: String,
}

pub fn compose_post() -> ComposePost {
    let user_idx = random_user();
    let user_str = "User".to_owned() + &user_idx.to_string();
    // let text = uuid::Uuid::new_v4().to_string();
    // let text = Alphanumeric.sample_string(&mut rand::thread_rng(), POST_SIZE);
    ComposePost {
        creator_id: user_str,
        text: SAMPLE_POST.get().unwrap().clone(),
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ReadTimeline {
    pub user_id: String,
}

pub fn read_timeline() -> ReadTimeline {
    let user_idx = random_user();
    let user_str = "User".to_owned() + &user_idx.to_string();
    ReadTimeline { user_id: user_str }
}

pub async fn service_compose_post(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("composepost");
    let c = compose_post();
    let req = serde_json::to_string(&c).unwrap();
    let ip: &String;
    if *SHARD.get().unwrap() == 1 {
        ip = COMPOSE_POST_IP.get().unwrap().get(0).unwrap();
    } else {
        let mut hasher = FxHasher::default();
        hasher.write(req.as_bytes());
        let hash_value = hasher.finish();
        let shard_idx = hash_value % *SHARD.get().unwrap() as u64;
        ip = COMPOSE_POST_IP.get().unwrap().get(shard_idx as usize).unwrap();
    }
    let _res = util::send_req(ip, "compose_post", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_read_hometimeline(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("hometimeline");
    let c = read_timeline();
    let req = serde_json::to_string(&c).unwrap();
    let ip: &String;
    if *SHARD.get().unwrap() == 1 {
        ip = HOME_TIMELINE_IP.get().unwrap().get(0).unwrap();
    } else {
        let mut hasher = FxHasher::default();
        hasher.write(req.as_bytes());
        let hash_value = hasher.finish();
        let shard_idx = hash_value % *SHARD.get().unwrap() as u64;
        ip = HOME_TIMELINE_IP.get().unwrap().get(shard_idx as usize).unwrap();
    }
    let _res = util::send_req(ip, "ro_read_home_timeline", req).await;
    // info!("hometimeline {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_read_usertimeline(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("usertimeline");
    let c = read_timeline();
    let req = serde_json::to_string(&c).unwrap();
    let ip: &String;
    if *SHARD.get().unwrap() == 1 {
        ip = USER_TIMELINE_IP.get().unwrap().get(0).unwrap();
    } else {
        let mut hasher = FxHasher::default();
        hasher.write(req.as_bytes());
        let hash_value = hasher.finish();
        let shard_idx = hash_value % *SHARD.get().unwrap() as u64;
        ip = USER_TIMELINE_IP.get().unwrap().get(shard_idx as usize).unwrap();
    }
    let _res = util::send_req(ip, "ro_read_user_timeline", req).await;
    // info!("usertimeline {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

pub async fn service_social(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let read_home_timeline_ratio = 60;
    let read_user_timeline_ratio = 30;
    let compose_post_ratio = 10;
    assert_eq!(
        read_home_timeline_ratio + read_user_timeline_ratio + compose_post_ratio,
        100
    );
    let coin = rand::random::<usize>() % 100;
    if coin < read_home_timeline_ratio {
        service_read_hometimeline(_req).await
    } else if coin < read_home_timeline_ratio + read_user_timeline_ratio {
        service_read_usertimeline(_req).await
    } else {
        service_compose_post(_req).await
    }
}
