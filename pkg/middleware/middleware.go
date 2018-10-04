package middleware

import "net/http"

type HandlerDecorator func (http.Handler) http.Handler

func Chain(method string, handler http.Handler, decorators ...HandlerDecorator) http.Handler {
	return http.HandlerFunc(func (res http.ResponseWriter, req *http.Request) {
		numberOfDecorators := len(decorators)

		if req.Method != method {
			return
		}

		if numberOfDecorators == 0 {
			handler.ServeHTTP(res, req)
		}

		current := handler
		for i := len(decorators) - 1; i >= 0; i-- {
			decorator := decorators[i]
			current = decorator(current)
		}

		current.ServeHTTP(res, req)
	})
}
