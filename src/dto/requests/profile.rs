#[derive(serde::Deserialize)]
pub struct CreateProfileRequest {
    pub username: String,
    pub email: String,
    pub first_name: String,
    pub last_name: String,
}