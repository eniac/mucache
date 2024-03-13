use crate::util;
use async_trait::async_trait;
use clap::Parser;
use hyper::{Body, Request, Response};
use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};
use std::convert::Infallible;
use tracing::info;


static FRONTEND_IPS: OnceCell<Vec<String>> = OnceCell::new();
static HIT_RATE: OnceCell<f32> = OnceCell::new();
static ENDPOINT: OnceCell<String> = OnceCell::new();
static CACHE: OnceCell<bool> = OnceCell::new();


#[derive(Parser)]
pub struct Fanin {
    #[clap(long, value_parser)]
    frontend1: String,
    #[clap(long, value_parser)]
    frontend2: String,
    #[clap(long, value_parser)]
    frontend3: String,
    #[clap(long, value_parser)]
    frontend4: String,
    #[clap(long, value_parser)]
    hitrate: f32,
    #[clap(long, value_parser)]
    endpoint: String,
}

#[async_trait]
impl util::Backend for Fanin {
    async fn prepare(&self) {
        let ips = vec![self.frontend1.clone(), 
                       self.frontend2.clone(), 
                       self.frontend3.clone(), 
                       self.frontend4.clone()];
        FRONTEND_IPS.set(ips).unwrap();
        assert!(self.hitrate <= 1.0);
        HIT_RATE.set(self.hitrate).unwrap();
        ENDPOINT.set(self.endpoint.clone()).unwrap();
        if self.endpoint.ends_with("hitormiss") {
            CACHE.set(true).unwrap();
        } else {
            CACHE.set(false).unwrap();
        }
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {

        let coin = rand::random::<usize>() % 4;

        if *CACHE.get().unwrap() {
            let r = ReadHitOrMissRequest { 
                    k: 1, 
                    hit_rate: *HIT_RATE.get().unwrap(),
            };
            let req = serde_json::to_string(&r).unwrap();
            let _res = util::send_req(
                &FRONTEND_IPS.get().unwrap()[coin],
                &ENDPOINT.get().unwrap(), req).await;
        } else {
            let r = ReadRequest { 
                k: 1, 
            };
            let req = serde_json::to_string(&r).unwrap();
            let _res = util::send_req(
                &FRONTEND_IPS.get().unwrap()[coin],
                &ENDPOINT.get().unwrap(), req).await;
            // info!("req {:?}", r);
            // info!("frontend {:?}", &FRONTEND_IPS.get().unwrap()[coin]);
            // info!("endpoint {:?}", &ENDPOINT.get().unwrap());
            // info!("read {:?}", _res);
        }
        Ok(Response::new(Body::from("OK!")))
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct ReadRequest {
    pub k: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct ReadHitOrMissRequest {
    pub k: i32,
    pub hit_rate: f32,
}

