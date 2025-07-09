package entity

import "errors"

// EventType
// This struct defines the type of event that is being logged, such as login or signup.
type EventType string

const (
	LOGIN_SUCCESS           EventType = "LOGIN_SUCCESS"
	SIGNUP_SUCCESS          EventType = "SIGNUP_SUCCESS"
	LOGIN_FAILURE           EventType = "LOGIN_FAILURE"
	LOGOUT                  EventType = "LOGOUT"
	PASSWORD_RESET_REQUEST  EventType = "PASSWORD_RESET_REQUEST"
	PASSWORD_CHANGE_SUCCESS EventType = "PASSWORD_CHANGE_SUCCESS"
)

type FailureReason string

const (
	INVALID_CREDENTIAL FailureReason = "INVALID_CREDENTIAL"
	USER_NOT_FOUND     FailureReason = "USER_NOT_FOUND"
	ACCOUNT_LOCKED     FailureReason = "ACCOUNT_LOCKED"
)

// User
// This event is triggered when a user performs an action that requires authentication, such as logging in or signing up.
type User struct {
	Uid   string `json:"uid" binding:"required"`   // Unique identifier for the user in Firebase Authentication
	Email string `json:"email" binding:"required"` // Email address of the user
	Name  string `json:"name,omitempty"`           // Name of the user
	Plan  string `json:"plan,omitempty"`           // Subscription plan of the user
}

// Client
// This struct contains information about the client making the request, such as IP address and user agent.
type Client struct {
	IpAddress string `json:"ipAddress"` // IP address of the client
	UserAgent string `json:"userAgent"` // User agent string of the client
}

// IsValid
// This method validates the User struct to ensure that required fields are present.
func (u *User) IsValid() (err error) {

	if u == nil {
		err = errors.New("invalid user: user event cannot be nil")
		return err
	}

	if u.Uid == "" {
		err = errors.New("invalid user: uid is required")
		return err
	}
	if u.Email == "" {
		err = errors.New("invalid user: email is required")
		return err
	}
	return err
}

// IsValid
// This method validates the Client struct to ensure that required fields are present.
func (c *Client) IsValid() (err error) {

	if c == nil {
		err = errors.New("invalid client: client event cannot be nil")
		return err
	}

	if c.IpAddress == "" {
		err = errors.New("invalid client: ipAddress is required")
		return err
	}
	if c.UserAgent == "" {
		err = errors.New("invalid client: userAgent is required")
		return err
	}
	return err
}

func (f *FailureReason) IsValid() (err error) {

	if f != nil {
		if *f != INVALID_CREDENTIAL && *f != USER_NOT_FOUND && *f != ACCOUNT_LOCKED {
			err = errors.New("invalid failure reason: failure reason must be one of the predefined values")
		}
	}

	return err

}

// IsValid
// This method validates the EventType to ensure that it is one of the predefined event types.
func (e *EventType) IsValid() (err error) {
	switch *e {
	case LOGIN_SUCCESS, SIGNUP_SUCCESS, LOGIN_FAILURE, LOGOUT, PASSWORD_RESET_REQUEST, PASSWORD_CHANGE_SUCCESS:
		return nil
	default:
		err = errors.New("invalid event type")
	}
	return err
}

// GetEventType
// This method returns a human-readable string representation of the event type.
func (e EventType) GetEventType() string {
	switch e {
	case LOGIN_SUCCESS:
		return "Login Success"
	case SIGNUP_SUCCESS:
		return "Signup Success"
	case LOGIN_FAILURE:
		return "Login Failure"
	case LOGOUT:
		return "Logout"
	case PASSWORD_RESET_REQUEST:
		return "Password Reset Request"
	case PASSWORD_CHANGE_SUCCESS:
		return "Password Change Success"
	default:
		return "Unknown Event Type"
	}
}

// SetEventType
// This method sets the event type based on a string representation.
func (e *EventType) SetEventType(eventType string) (err error) {

	if e == nil {
		return errors.New("event type cannot be nil")
	}

	switch eventType {
	case "LOGIN_SUCCESS":
		*e = LOGIN_SUCCESS
	case "SIGNUP_SUCCESS":
		*e = SIGNUP_SUCCESS
	case "LOGIN_FAILURE":
		*e = LOGIN_FAILURE
	case "LOGOUT":
		*e = LOGOUT
	case "PASSWORD_RESET_REQUEST":
		*e = PASSWORD_RESET_REQUEST
	case "PASSWORD_CHANGE_SUCCESS":
		*e = PASSWORD_CHANGE_SUCCESS
	default:
		err = errors.New("invalid event type")
	}
	return err
}

// String returns a string representation of the EventType
func (e EventType) String() string {
	return string(e)
}

// IsLoginEvent checks if the event type is a login-related event
func (e EventType) IsLoginEvent() bool {
	return e == LOGIN_SUCCESS || e == LOGIN_FAILURE
}

// IsSuccessEvent checks if the event type represents a successful operation
func (e EventType) IsSuccessEvent() bool {
	return e == LOGIN_SUCCESS || e == SIGNUP_SUCCESS || e == PASSWORD_CHANGE_SUCCESS
}

// IsFailureEvent checks if the event type represents a failed operation
func (e EventType) IsFailureEvent() bool {
	return e == LOGIN_FAILURE
}
