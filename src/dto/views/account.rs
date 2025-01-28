use chrono::Utc;
use serde::{Deserialize, Serialize};
use crate::domain::account::Account;

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct AccountView {
    pub id: i32,
    pub name: String,
    pub r#type: String,
    pub starting_amount: f64,
    pub creation_date: chrono::DateTime<Utc>,
}

impl AccountView {
    pub fn from(account: Account) -> AccountView {
        AccountView {
            id: account.id,
            name: account.name,
            r#type: account.r#type,
            starting_amount: account.starting_amount,
            creation_date: account.creation_date,
        }
    }
}