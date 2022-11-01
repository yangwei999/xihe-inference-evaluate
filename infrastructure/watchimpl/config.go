package watchimpl

type Config struct {
	// RPC Endpoint
	InferenceEndpoint string `json:"inference_endpoint" required:"true"`
}
