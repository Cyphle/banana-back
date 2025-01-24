use serde::Deserialize;

#[derive(Deserialize)]
pub struct CreateAccountRequest {
    pub name: String,
    pub r#type: String,
    pub starting_amount: f64,
}