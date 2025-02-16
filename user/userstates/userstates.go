package userstates

type UserState int

const (
    None                        UserState = iota
    Start                       UserState = iota
    awaiting_amount_to_split    UserState = iota
    awaiting_new_contact_name   UserState = iota
    waiting_for_split_contacts  UserState = iota
)


// String method for better printing
func (s UserState) String() string {
	return [...]string{"None", "Start", "awaiting_amount_to_split", "awaiting_new_contact_name", "waiting_for_split_contacts"} [s]
}
