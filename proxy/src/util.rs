use async_trait::async_trait;
use bytes::Bytes;
use hyper::{Body, Request, Response};
use std::convert::Infallible;

thread_local! {
    pub static GCLIENT: hyper::Client<hyper::client::HttpConnector, hyper::Body> = hyper::Client::new();
}

#[async_trait]
pub trait Backend {
    async fn prepare(&self);

    async fn run(_req: Request<Body>) -> Result<Response<Body>, Infallible>;
}

pub async fn send_req(ip: &str, method: &str, req: String) -> Bytes {
    let url = "http://".to_owned() + ip + "/" + method;
    let r = hyper::Request::builder()
        .method(hyper::Method::POST)
        .uri(url)
        .header("content-type", "application/json")
        .header("Connection", "close")
        .body(hyper::Body::from(req))
        .unwrap();
    let client = GCLIENT.with(|c| c.clone());
    let resp = client.request(r).await.unwrap();
    hyper::body::to_bytes(resp.into_body()).await.unwrap()
}

// How to do it elegantly?
#[macro_export]
macro_rules! impl_backend {
    ($($app:ident),*) => {
        impl App {
            async fn prepare(&self) {
                match self {
                    $(
                        App::$app(inner) => inner.prepare().await,
                    )*
                }
            }

            async fn run(&self) {
                match self {
                    $(
                        App::$app(_inner) => {
                            let addr = std::net::SocketAddr::from(([127, 0, 0, 1], 3000));
                            let make_svc =
                                make_service_fn(|_conn| async { Ok::<_, Infallible>(service_fn($app::run)) });
                            let server = Server::bind(&addr).serve(make_svc);
                            server.await.unwrap();
                        }
                    )*
                }
            }
        }
    };
}
