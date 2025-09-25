package main

type stripeConfig struct {
	webhookSecret string
	apiKey        string
}

func newStripeConfig() (*stripeConfig, error) {
	webhookSecret, err := getRequiredStringEnv("STRIPE_WEBHOOK_SECRET")
	if err != nil {
		return nil, err
	}

	apiKey, err := getRequiredStringEnv("STRIPE_API_KEY")
	if err != nil {
		return nil, err
	}

	return &stripeConfig{
		webhookSecret: *webhookSecret,
		apiKey:        *apiKey,
	}, nil
}
