use crate::util;
use async_trait::async_trait;
use clap::Parser;
use hyper::{Body, Request, Response};
use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};
use std::convert::Infallible;

static CALLER_IP: OnceCell<String> = OnceCell::new();
static HIT_RATE: OnceCell<usize> = OnceCell::new();

#[derive(Parser)]
pub struct TwoServices {
    #[clap(long, value_parser)]
    caller: String,
    #[clap(long, value_parser)]
    hitrate: usize,
}

#[async_trait]
impl util::Backend for TwoServices {
    async fn prepare(&self) {
        CALLER_IP.set(self.caller.clone()).unwrap();
        assert!(self.hitrate <= 100);
        HIT_RATE.set(self.hitrate).unwrap();

        // write k: 1, v: 1 to the database at callee
        let wr = WriteRequest { k: 1, v: 1 };
        let req = serde_json::to_string(&wr).unwrap();
        let _res = util::send_req(CALLER_IP.get().unwrap(), "write", req).await;

        // read k: 1 from the caller once, warming up the cache
        let rr = ReadRequest { k: 1 };
        let req = serde_json::to_string(&rr).unwrap();
        let _res = util::send_req(CALLER_IP.get().unwrap(), "read", req).await;
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        // TODO: move coin to servers
        let coin = rand::random::<usize>() % 101;
        let r = ReadRequest { k: 1 };
        let req = serde_json::to_string(&r).unwrap();
        if coin < *HIT_RATE.get().unwrap() {
            let _res = util::send_req(CALLER_IP.get().unwrap(), "hit", req).await;
        } else {
            let _res = util::send_req(CALLER_IP.get().unwrap(), "miss", req).await;
        }
        Ok(Response::new(Body::from("OK!")))
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct ReadRequest {
    pub k: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct WriteRequest {
    pub k: i32,
    pub v: i32,
}
