package auth

// Auther provides an interface to implement for AA an individual user
type Auther interface {
	//IsValid determines whether or not we have a valid user
	IsValid(user string, pass string) bool

	//IsAuthorizedForScope determines if the user is allowed to read / write for this scope
	IsUserAuthorizedForScope(user string, scope string) bool

	//Terminate is the handler to call when the program exist
	Terminate()
}
