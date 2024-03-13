use crate::util;
use async_trait::async_trait;
use clap::Parser;
use hyper::{Body, Request, Response};
use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};
use std::convert::Infallible;

static CALLER_IP: OnceCell<String> = OnceCell::new();
static CALLEE_IP: OnceCell<String> = OnceCell::new();
static TOTAL_KEYS: OnceCell<usize> = OnceCell::new();

#[derive(Parser)]
pub struct CMThroughput {
    #[clap(long, value_parser)]
    caller: String,
    #[clap(long, value_parser)]
    callee: String,
    #[clap(long, value_parser)]
    total_keys: usize,
}

#[async_trait]
impl util::Backend for CMThroughput {
    async fn prepare(&self) {
        CALLER_IP.set(self.caller.clone()).unwrap();
        CALLEE_IP.set(self.callee.clone()).unwrap();
        TOTAL_KEYS.set(self.total_keys.clone()).unwrap();

        for n in 0..*TOTAL_KEYS.get().unwrap() {
            // write k: 1, v: 1 to the database at callee
            let wr = WriteRequest { k: n, v: 1 };
            let req = serde_json::to_string(&wr).unwrap();
            let _res = util::send_req(CALLEE_IP.get().unwrap(), "write", req).await;
        }
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        let coin = rand::random::<usize>() % 101;

        if coin < 80 {
            let r = ReadRequest { k: random_key() };
            let req = serde_json::to_string(&r).unwrap();
            let _res = util::send_req(CALLER_IP.get().unwrap(), "read", req).await;
        } else {
            let wr = WriteRequest {
                k: random_key(),
                v: 1,
            };
            let req = serde_json::to_string(&wr).unwrap();
            let _res = util::send_req(CALLEE_IP.get().unwrap(), "write", req).await;
        }
        Ok(Response::new(Body::from("OK!")))
    }
}

pub fn random_key() -> usize {
    let max_key = TOTAL_KEYS.get().unwrap();
    rand::random::<usize>() % max_key
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct ReadRequest {
    pub k: usize,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct WriteRequest {
    pub k: usize,
    pub v: i32,
}
