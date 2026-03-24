package functionurl

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func WrapBuffered(handler http.Handler) func(context.Context, events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	return func(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
		var body io.Reader = strings.NewReader(request.Body)
		if request.IsBase64Encoded {
			body = base64.NewDecoder(base64.StdEncoding, body)
		}

		url := "https://" + request.RequestContext.DomainName + request.RawPath
		if request.RawQueryString != "" {
			url += "?" + request.RawQueryString
		}

		httpRequest, err := http.NewRequestWithContext(ctx, request.RequestContext.HTTP.Method, url, body)
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}
		httpRequest.RemoteAddr = request.RequestContext.HTTP.SourceIP

		for key, value := range request.Headers {
			httpRequest.Header.Add(key, value)
		}

		if len(request.Cookies) > 0 && httpRequest.Header.Get("Cookie") == "" {
			httpRequest.Header.Set("Cookie", strings.Join(request.Cookies, "; "))
		}

		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httpRequest)

		result := recorder.Result()
		defer result.Body.Close()

		responseBody, err := io.ReadAll(result.Body)
		if err != nil {
			return events.LambdaFunctionURLResponse{}, err
		}

		response := events.LambdaFunctionURLResponse{
			StatusCode: result.StatusCode,
			Headers:    make(map[string]string, len(result.Header)),
		}

		for key, values := range result.Header {
			if key == "Set-Cookie" {
				response.Cookies = append(response.Cookies, values...)
				continue
			}
			response.Headers[key] = strings.Join(values, ",")
		}

		if len(responseBody) == 0 {
			return response, nil
		}

		if utf8.Valid(responseBody) {
			response.Body = string(responseBody)
			return response, nil
		}

		response.Body = base64.StdEncoding.EncodeToString(responseBody)
		response.IsBase64Encoded = true
		return response, nil
	}
}

func StartBuffered(handler http.Handler, options ...lambda.Option) {
	lambda.StartHandlerFunc(WrapBuffered(handler), options...)
}
