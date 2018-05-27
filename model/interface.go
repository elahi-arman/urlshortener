package model

//Linker is the interface to implement Link retrieval from a data store
type Linker interface {
	//CommitLink makes a change in the data store for the given link
	CommitLink(key string, l *Link) error

	//GetLink retrieves a link from the data store
	GetLink(scope string, user string, title string) (*Link, error)

	//GetLinksInScope retrieves all links in the given scope
	GetLinksInScope(scope string) ([]Link, error)

	//GetLinksByUser retrieves all links (from all scopes) for the given user
	GetLinksByUser(user string) ([]Link, error)

	//GetLinksInScopeByUser retrieves all links inthe given scope from the given user
	GetLinksInScopeByUser(scope string, user string) ([]Link, error)
}
