CREATE TABLE payment_intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_payment_intent_id VARCHAR(255) UNIQUE NOT NULL,
    amount BIGINT NOT NULL, -- Amount in smallest currency unit (cents)
    currency VARCHAR(3) NOT NULL DEFAULT 'usd',
    status VARCHAR(50) NOT NULL,
    client_secret VARCHAR(255),
    customer_id UUID,
    metadata JSONB DEFAULT '{}',
    description TEXT,
    receipt_email VARCHAR(255),
    shipping_address JSONB,
    billing_address JSONB,
    payment_method_id VARCHAR(255),
    payment_method_types TEXT[] DEFAULT ARRAY['card'],
    setup_future_usage VARCHAR(20),
    capture_method VARCHAR(20) DEFAULT 'automatic',
    confirmation_method VARCHAR(20) DEFAULT 'automatic',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    canceled_at TIMESTAMP WITH TIME ZONE,
    succeeded_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_payment_intents_stripe_id ON payment_intents(stripe_payment_intent_id);
CREATE INDEX idx_payment_intents_customer_id ON payment_intents(customer_id);
CREATE INDEX idx_payment_intents_status ON payment_intents(status);
CREATE INDEX idx_payment_intents_created_at ON payment_intents(created_at);