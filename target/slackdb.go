package target

// SlackDBModel defines the JSON serialization format for saving slack targets'
// settings in the database.
type SlackDBModel struct {
	Channel string
}
