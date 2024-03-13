use hyper::service::{make_service_fn, service_fn};
use hyper::Server;
use std::convert::Infallible;
use util::Backend;
mod cm;
mod cmthroughput;
mod hotel;
mod movie;
mod social;
mod boutique;
mod twoservices;
mod fanin;
mod util;
use clap::Parser;

// Don't forget to export inner structs
// They're used in the macro
use cm::CM;
use cmthroughput::CMThroughput;
use hotel::Hotel;
use movie::Movie;
use social::Social;
use boutique::Boutique;
use twoservices::TwoServices;
use fanin::Fanin;

#[derive(Parser)]
pub enum App {
    CM(CM),
    CMThroughput(CMThroughput),
    Hotel(Hotel),
    Social(Social),
    Boutique(Boutique),
    Movie(Movie),
    TwoServices(TwoServices),
    Fanin(Fanin),
}

impl_backend!(CM, CMThroughput, Hotel, Social, Boutique, Movie, TwoServices, Fanin);

#[tokio::main(worker_threads = 4)]
async fn main() {
    tracing_subscriber::fmt::init();
    let app = App::parse();
    app.prepare().await;
    app.run().await;
}
