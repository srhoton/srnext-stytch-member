package models

// ALBRequest represents an Application Load Balancer request event
type ALBRequest struct {
	RequestContext                  RequestContext      `json:"requestContext"`
	HTTPMethod                      string              `json:"httpMethod"`
	Path                            string              `json:"path"`
	QueryStringParameters           map[string]string   `json:"queryStringParameters,omitempty"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters,omitempty"`
	Headers                         map[string]string   `json:"headers,omitempty"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders,omitempty"`
	Body                            string              `json:"body,omitempty"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded"`
	PathParameters                  map[string]string   `json:"pathParameters,omitempty"`
}

// RequestContext contains the request context from ALB
type RequestContext struct {
	ELB ELBContext `json:"elb"`
}

// ELBContext contains ELB-specific request context
type ELBContext struct {
	TargetGroupArn string `json:"targetGroupArn"`
}

// ALBResponse represents the response structure for ALB
type ALBResponse struct {
	StatusCode        int                 `json:"statusCode"`
	StatusDescription string              `json:"statusDescription,omitempty"`
	Headers           map[string]string   `json:"headers,omitempty"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders,omitempty"`
	Body              string              `json:"body,omitempty"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
}
