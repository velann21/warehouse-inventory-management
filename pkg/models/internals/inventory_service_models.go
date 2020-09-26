package internals

type SuccessfullyAddedArticle struct {
	Name     string
	Total    string
	ID       string
	Endpoint string
}

type FailedArticle struct {
	Name   string
	Total  string
	Reason string
}
