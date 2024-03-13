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
static CURRENCY_IP: OnceCell<String> = OnceCell::new();
static CART_IP: OnceCell<String> = OnceCell::new();

static MAX_USER: OnceCell<usize> = OnceCell::new();
static TOTAL_PRODUCTS: OnceCell<usize> = OnceCell::new();
static CATALOG_SIZE: OnceCell<usize> = OnceCell::new();

#[derive(Parser)]
pub struct Boutique {
    #[clap(long, value_parser, multiple = true)]
    frontend: String,
    #[clap(long, value_parser, multiple = true)]
    currency: String,
    #[clap(long, value_parser, multiple = true)]
    cart: String,
    #[clap(long, value_parser, multiple = true)]
    catalog_size: usize,
}

#[async_trait]
impl util::Backend for Boutique {
    async fn prepare(&self) {
        FRONTEND_IP.set(self.frontend.clone()).unwrap();
        CURRENCY_IP.set(self.currency.clone()).unwrap();
        CART_IP.set(self.cart.clone()).unwrap();
        CATALOG_SIZE.set(self.catalog_size.clone()).unwrap();
        let username = whoami::username();
        let analysis_file =
            format!("/users/{username}/mucache/experiments/boutique/data/analysis.txt");
        info!("Reading analysis from {}", analysis_file);
        read_analysis_file(&analysis_file);
    }

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
        service_boutique(_req).await
    }
}

pub fn read_analysis_file(path: &str) {
    let file = File::open(path).unwrap();
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader.lines().map(|l| l.unwrap()).collect();

    let max_user = usize::from_str(lines[0].split(':').nth(1).unwrap().trim()).unwrap();
    let total_products = usize::from_str(lines[1].split(':').nth(1).unwrap().trim()).unwrap();
    // We get this from a proxy argument
    // let _catalog_size = usize::from_str(lines[2].split(':').nth(1).unwrap().trim()).unwrap();

    MAX_USER.set(max_user).unwrap();
    TOTAL_PRODUCTS.set(total_products).unwrap();
}

// Read the following link to use weigthed indexing:
//   https://docs.rs/rand/latest/rand/distributions/weighted/struct.WeightedIndex.html
// Or even better the following (0(1) sampling):
//   https://docs.rs/rand_distr/0.4.3/rand_distr/weighted_alias/struct.WeightedAliasIndex.html
pub fn random_user() -> usize {
    let max_user = MAX_USER.get().unwrap();
    rand::random::<usize>() % max_user + 1
}

pub fn random_product() -> usize {
    let max_product = TOTAL_PRODUCTS.get().unwrap();
    rand::random::<usize>() % max_product
}

pub fn random_product_str() -> String {
    let product_idx = random_product();
    "p".to_owned() + &product_idx.to_string() 
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Home {
    pub user_id: String,
    pub catalog_size: usize,
}

pub fn home() -> Home {
    let user_idx = random_user();
    let user_str = &user_idx.to_string();
    Home {
        user_id: user_str.to_string(),
        catalog_size: *CATALOG_SIZE.get().unwrap(),
    }
}

pub async fn service_home(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("home");
    let c = home();
    let req = serde_json::to_string(&c).unwrap();
    let ip = FRONTEND_IP.get().unwrap();
    let _res = util::send_req(ip, "ro_home", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}


#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct BrowseProduct {
    pub product_id: String,
}

pub fn browse_product() -> BrowseProduct {
    let product_str = random_product_str();
    BrowseProduct { product_id: product_str }
}

pub async fn service_browse_product(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("home");
    let c = browse_product();
    let req = serde_json::to_string(&c).unwrap();
    // info!("browse product req {:?}", req);
    let ip = FRONTEND_IP.get().unwrap();
    let _res = util::send_req(ip, "ro_browse_product", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SetCurrency {
    pub currencyCode: String,
    pub rate: String,
}

// TODO: This might be wrong
pub fn set_currency() -> SetCurrency {
    SetCurrency { currencyCode: "EUR".to_string(), rate: "1.2345".to_string() }
}

pub async fn service_set_currency(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("home");
    let c = set_currency();
    let req = serde_json::to_string(&c).unwrap();
    let ip = CURRENCY_IP.get().unwrap();
    let _res = util::send_req(ip, "set_currency", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}



#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ViewCart {
    pub user_id: String,
}

pub fn view_cart() -> ViewCart {
    let user_idx = random_user();
    let user_str = &user_idx.to_string();
    ViewCart { user_id: user_str.to_string() }
}

pub async fn service_view_cart(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("home");
    let c = view_cart();
    let req = serde_json::to_string(&c).unwrap();
    let ip = FRONTEND_IP.get().unwrap();
    let _res = util::send_req(ip, "ro_view_cart", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AddToCart {
    pub user_id: String,
    pub product_id: String,
    pub quantity: i32,
}

pub fn add_to_cart() -> AddToCart {
    let user_idx = random_user();
    let user_str = &user_idx.to_string();
    let product_str = random_product_str();
    AddToCart { 
        user_id: user_str.to_string(),
        product_id: product_str,
        quantity: 1
    }
}

pub async fn service_add_to_cart(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("home");
    let c = add_to_cart();
    let req = serde_json::to_string(&c).unwrap();
    let ip = CART_IP.get().unwrap();
    let _res = util::send_req(ip, "add_item", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Address {
    pub street_address: String,
    pub city: String,
    pub state: String,
    pub country: String,
    pub zip_code: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreditCard {
    pub card_number: String,
    pub card_type: String,
    pub expiration_month: i32,
    pub expiration_year: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Checkout {
    pub user_id: String,
    pub user_currency: String,
    pub address: Address,
    pub email: String,
    pub credit_card: CreditCard,
}

pub fn address() -> Address {
    Address { 
        street_address: "1600 Amphitheatre Parkway".to_string(),
        zip_code: 94043,
        city: "Mountain View".to_string(),
        state: "CA".to_string(),
        country: "United States".to_string(),
    }
}

pub fn credit_card() -> CreditCard {
    CreditCard { 
        card_number: "4432-8015-6152-0454".to_string(),
        card_type: "visa".to_string(),
        expiration_month: 1,
        expiration_year: 2039,
    }
}

pub fn checkout() -> Checkout {
    let user_idx = random_user();
    let user_str = &user_idx.to_string();
    Checkout { 
        user_id: user_str.to_string(),
        user_currency: "EUR".to_string(),
        address: address(),
        email: "someone@example.com".to_string(),
        credit_card: credit_card(),
    }
}

pub async fn service_checkout(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    // info!("home");
    let c = checkout();
    let req = serde_json::to_string(&c).unwrap();
    let ip = FRONTEND_IP.get().unwrap();
    let _res = util::send_req(ip, "checkout", req).await;
    // info!("compose post {:?}", _res);
    Ok(Response::new(Body::from("OK!")))
}


// Frontend: home + browse_prod + view_cart + checkout
// Frontend: 10 + 50 + 15 + 5 = 80%
// Currency: 10%
// Cart: 10%

pub async fn service_boutique(_req: Request<Body>) -> Result<Response<Body>, Infallible> {
    let home_ratio = 10;
    let set_currency_ratio = 10;
    let browse_product_ratio = 50;
    let add_to_cart_ratio = 10;
    let view_cart_ratio = 15;
    let checkout_ratio = 5;
    assert_eq!(
        home_ratio + set_currency_ratio + browse_product_ratio + add_to_cart_ratio + view_cart_ratio + checkout_ratio,
        100
    );
    let coin = rand::random::<usize>() % 100;
    if coin < home_ratio {
        service_home(_req).await
    } else if coin < home_ratio + set_currency_ratio {
        service_set_currency(_req).await
    } else if coin < home_ratio + set_currency_ratio + browse_product_ratio {
        service_browse_product(_req).await
    } else if coin < home_ratio + set_currency_ratio + browse_product_ratio + add_to_cart_ratio {
        service_add_to_cart(_req).await
    } else if coin < home_ratio + set_currency_ratio + browse_product_ratio + add_to_cart_ratio + view_cart_ratio {
        service_view_cart(_req).await
    } else {
        service_checkout(_req).await
    }
}
