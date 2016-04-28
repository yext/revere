package target

// EmailDBModel defines the JSON serialization format for saving email targets'
// settings in the database.
type EmailDBModel struct {
	Addresses []struct {
		To      string
		ReplyTo string
	}
}
