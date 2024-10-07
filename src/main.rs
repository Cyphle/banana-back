mod config;

#[actix_web::main]
async fn main() {
    config::logger::config();

    log::info!("Starting the application");

    let db = config::database::connect().await.unwrap();
    let static_db = Box::leak(Box::new(db));


    log::info!("Application is now closed");
    // TODO faudra trouver un moyen de close la connexion.
    // db.close().unwrap()
}
