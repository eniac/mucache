use crate::util;
use async_trait::async_trait;
use clap::Parser;
use hyper::{Body, Request, Response};
use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};
use std::convert::Infallible;
use tracing::info;

static MAX_KEY: OnceCell<u32> = OnceCell::new();
static MAX_NUMBER_OF_KEYS: OnceCell<usize> = OnceCell::new();
static MAX_CALLARGS: OnceCell<u32> = OnceCell::new();

static SERVER_IP: OnceCell<String> = OnceCell::new();

#[derive(Parser)]
pub struct CM {
    #[clap(long, value_parser)]
    ip: String,
}

#[async_trait]
impl util::Backend for CM {
    async fn prepare(&self) {
        SERVER_IP.set(self.ip.clone()).unwrap();
        set_globals(10, 3, 10);
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        service_cache_manager(_req).await
    }
}

fn set_globals(max_key: u32, max_number_of_keys: usize, max_callargs: u32) {
    MAX_KEY.set(max_key).unwrap();
    MAX_NUMBER_OF_KEYS.set(max_number_of_keys).unwrap();
    MAX_CALLARGS.set(max_callargs).unwrap();
}

fn random_key() -> u32 {
    let max_key = MAX_KEY.get().unwrap();
    rand::random::<u32>() % max_key + 1
}

fn random_number_of_keys() -> usize {
    let max_number_of_keys = MAX_NUMBER_OF_KEYS.get().unwrap();
    rand::random::<usize>() % max_number_of_keys + 1
}

fn random_callargs() -> u32 {
    let max_callargs = MAX_CALLARGS.get().unwrap();
    rand::random::<u32>() % max_callargs + 1
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StartRequest {
    pub callargs: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EndRequest {
    pub callargs: String,
    pub caller: String,
    pub deps: Vec<String>,
    pub returnval: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct InvalidateRequest {
    pub key: String,
}

pub fn start_request() -> StartRequest {
    let callargs_idx = random_callargs();
    let callargs_str = "req".to_owned() + &callargs_idx.to_string();
    StartRequest {
        callargs: callargs_str,
    }
}

pub fn end_request() -> EndRequest {
    let callargs_idx = random_callargs();
    let callargs_str = "req".to_owned() + &callargs_idx.to_string();
    let returnval_str = "ret".to_owned() + &callargs_idx.to_string();
    let deps_number = random_number_of_keys();
    let mut deps = Vec::with_capacity(deps_number);
    for _ in 0..deps_number {
        let key_idx = random_key();
        let key_str = "key".to_owned() + &key_idx.to_string();
        deps.push(key_str);
    }
    EndRequest {
        callargs: callargs_str,
        caller: "service".to_string(),
        deps,
        returnval: returnval_str,
    }
}

pub fn inv_request() -> InvalidateRequest {
    let key_idx = random_key();
    let key_str = "key".to_owned() + &key_idx.to_string();
    InvalidateRequest { key: key_str }
}

pub async fn service_cache_manager(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let request_type = rand::random::<u32>() % 3;
    let req = match request_type {
        0 => serde_json::to_string(&(start_request())).unwrap(),
        1 => serde_json::to_string(&(end_request())).unwrap(),
        2 => serde_json::to_string(&(inv_request())).unwrap(),
        _ => panic!("Unknown request type"),
    };
    let method = match request_type {
        0 => "start",
        1 => "end",
        2 => "inv",
        _ => panic!("Unknown request type"),
    };
    // info!("{:?}", req);
    let res = util::send_req(SERVER_IP.get().unwrap(), method, req).await;
    info!("{:?}", res);
    Ok(Response::new(Body::from("OK!")))
}
