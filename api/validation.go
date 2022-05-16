package api

// type OpenApi struct {
// 	doc    *openapi3.T
// 	router routers.Router
// }

// func NewOpenApi(filename string) OpenApi {
// 	doc, err := openapi3.NewLoader().LoadFromFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	router, err := legacyrouter.NewRouter(doc)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return OpenApi{doc, router}
// }

// func (oa *OpenApi) validateRequest(ctx context.Context, httpReq *http.Request) {
// 	_ = oa.doc.Validate(ctx)
// 	// Find route
// 	route, pathParams, _ := oa.router.FindRoute(httpReq)

// 	// Validate request
// 	requestValidationInput := &openapi3filter.RequestValidationInput{
// 		Request:    httpReq,
// 		PathParams: pathParams,
// 		Route:      route,
// 	}
// 	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
// 		panic(err)
// 	}

// }

// func (oa *OpenApi) validateResponse(ctx context.Context, httpResp *http.Response) {
// 	_ = oa.doc.Validate(ctx)
// 	// Find route
// 	route, pathParams, _ := oa.router.FindRoute(httpResp.Request)

// 	// Validate request
// 	requestValidationInput := &openapi3filter.RequestValidationInput{
// 		Request:    httpResp.Request,
// 		PathParams: pathParams,
// 		Route:      route,
// 	}

// 	responseValidationInput := &openapi3filter.ResponseValidationInput{
// 		RequestValidationInput: requestValidationInput,
// 		Status:                 httpResp.StatusCode,
// 		Header:                 httpResp.Header,
// 	}
// 	if httpResp.Body != nil {
// 		data, _ := json.Marshal(httpResp.Body)
// 		responseValidationInput.SetBodyBytes(data)
// 	}

// 	// Validate response.
// 	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
// 		panic(err)
// 	}
// }

// // "github.com/dgrijalva/jwt-go"
// // A Util function to generate jwt_token which can be used in the request header
// func GenToken(id uint) string {
// 	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
// 	// Set some claims
// 	jwt_token.Claims = jwt.MapClaims{
// 		"id":  id,
// 		"exp": time.Now().Add(time.Hour * 24).Unix(),
// 	}
// 	// Sign and get the complete encoded token as a string
// 	token, _ := jwt_token.SignedString([]byte(NBSecretPassword))
// 	return token
// }
