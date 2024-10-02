package profile

type Profile struct {
	ID        int64
	Username  string
	Email     string
	firstName string
	lastName  string
}

type CreateProfileCommand struct {
	Username  string
	Email     string
	firstName string
	lastName  string
}
