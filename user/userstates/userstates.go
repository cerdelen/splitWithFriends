package userstates

type UserState int

const (
    None                            UserState = iota
    Start                           UserState = iota

    Configuration                   UserState = iota
    NewSplit                        UserState = iota
    DirectRequest                   UserState = iota

    // Configuration
    AddingContact                   UserState = iota
    RemovingContact                 UserState = iota

    // New Split
    AddContactsToSplit              UserState = iota

    // DirectRequest
    NewDirectRequest                UserState = iota
    AwaitingAmountDirectRequest     UserState = iota









    Awaiting_amount_to_split        UserState = iota
    Awaiting_new_contact_name       UserState = iota
    Awaiting_for_split_contacts     UserState = iota
    Awaiting_for_split_by_amount    UserState = iota
)


// String method for better printing
func (s UserState) String() string {
    return [...]string{"None", "Start", "Configuration", "NewSplit", "DirectRequest", "AddingContact", "RemovingContact", "AddContactsToSplit", "AwaitingAmountDirectRequest", "NewDirectRequest", "Awaiting_amount_to_split", "Awaiting_new_contact_name", "Awaiting_for_split_contacts", "Awaiting_for_split_by_amount"} [s]
}
