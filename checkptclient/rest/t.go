package rest

// type RequestInit struct {
// 	r        *http.Request
// 	meth     HTTPMethod
// 	contType string
// 	handler  ErrorHandler
// 	url      string
// }
// type Builder struct {
// 	init RequestInit
// }

// func NewBuilder(url string) *Builder {
// 	return &Builder{
// 		init: RequestInit{
// 			url:      url,
// 			contType: "application/json",
// 			handler:  DefaultErrorHandler{},
// 		},
// 	}
// }
// func (b *Builder) Method(method HTTPMethod) *Builder {
// 	b.init.meth = method
// 	return b
// }
// func (b *Builder) ContentType(cont string) *Builder {
// 	b.init.contType = cont
// 	return b
// }
// func (b *Builder) ErrorHandler(handler ErrorHandler) *Builder {
// 	b.init.handler = handler
// 	return b
// }
// func (b *Builder) Build() *RequestInit {
// 	return &b.init
// }
