package userstates

type UserState int

const (
    None                            UserState = iota
    Start                           UserState = iota
    Configuration                   UserState = iota

    AddingContact                   UserState = iota
    RemovingContact                 UserState = iota




    RequestFromSingleContact        UserState = iota

    Awaiting_amount_to_split        UserState = iota
    Awaiting_new_contact_name       UserState = iota
    Awaiting_for_split_contacts     UserState = iota
    Awaiting_for_split_by_amount    UserState = iota
)


// String method for better printing
func (s UserState) String() string {
	return [...]string{"None", "Start", "awaiting_amount_to_split", "awaiting_new_contact_name", "waiting_for_split_contacts"} [s]
}
