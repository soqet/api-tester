package jsonreader

type MainSchema struct {
	// api endpoint
	Endpoint string `json:"endpoint"`
	// groups of requests
	Groups []GroupSchema `json:"groups"`
	// requests without groups
	Requests []RequestSchema `json:"requests"`
}

type GroupSchema struct {
	// group name
	Name string `json:"name"`
	// endpoint overwrite for group
	Endpoint string `json:"endpoint"`
	// requests of group
	Requests []RequestSchema `json:"requests"`
}

type RequestSchema struct {
	// http method, may be capitalized
	Method string `json:"method"`
	// endpoint overwrite for request
	Endpoint string `json:"endpoint"`
	// requested resourse.
	// final url will be endpoint+resourse
	Resourse string `json:"resourse"`
	// request body
	Body string `json:"body"`
	// file containing body
	BodyFile string `json:"body-file"`
	// response code
	Code int `json:"code"`
	// response body
	Response string `json:"response"`
	// file containing response
	ResponseFile string `json:"response-file"`
}
