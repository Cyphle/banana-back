use chrono::Utc;
use crate::domain::profile::Profile;

#[derive(Debug, PartialEq)]
pub struct Account {
    pub id: i32,
    pub name: String,
    pub r#type: String,
    pub starting_amount: f64,
    pub creation_date: chrono::DateTime<Utc>,
}

pub struct CreateAccountCommand {
    pub name: String,
    pub r#type: String,
    pub starting_amount: f64,
    pub profile_id: i32
}