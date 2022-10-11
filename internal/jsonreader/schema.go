package jsonreader

type MainSchema struct {
	// api endpoint
	Endpoint string `json:"endpoint"`
	// headers which will be setted for all requests
	Headers []*HeaderSchema `json:"headers"`
	// groups of requests
	Groups []*GroupSchema `json:"groups"`
	// non group requests
	Requests []*RequestSchema `json:"requests"`
}

type GroupSchema struct {
	// endpoint overwrite for group
	Endpoint string `json:"endpoint"`
	// header overwrites
	// overwrite main headers if specified
	Headers []*HeaderSchema `json:"headers"`
	// requests of group
	Requests []*RequestSchema `json:"requests"`
}

type RequestSchema struct {
	// http method, may be capitalized
	Method string `json:"method"`
	// additional headers for this request
	// added to the main or group headers
	Headers []*HeaderSchema `json:"headers"`
	// endpoint overwrite for request
	Endpoint string `json:"endpoint"`
	// requested resourse.
	// final url will be endpoint+resourse
	Resourse string `json:"resourse"`
	// request body
	Body string `json:"body"`
	// text file containing body
	// igores Body if not empty
	BodyFile string `json:"body-file"`
	// response code
	Code int `json:"code"`
	// response body
	Response string `json:"response"`
	// text file containing response body
	// igores Response if not empty
	ResponseFile string `json:"response-file"`
}

type HeaderSchema struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
