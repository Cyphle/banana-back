use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct ProfileView {
    pub id: i32,
    pub username: String,
    pub email: String,
    pub first_name: String,
    pub last_name: String,
}