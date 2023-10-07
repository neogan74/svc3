package web

type Midleware func(Handler) Handler

func wrapMiddleware(mw []Midleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one . replace the
	// handler with new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}
