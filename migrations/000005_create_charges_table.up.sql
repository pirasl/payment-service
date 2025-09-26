CREATE TABLE charges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_charge_id VARCHAR(255) UNIQUE NOT NULL,
    payment_intent_id UUID REFERENCES payment_intents(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    amount_captured BIGINT DEFAULT 0,
    amount_refunded BIGINT DEFAULT 0,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) NOT NULL,
    paid BOOLEAN DEFAULT FALSE,
    refunded BOOLEAN DEFAULT FALSE,
    captured BOOLEAN DEFAULT FALSE,
    disputed BOOLEAN DEFAULT FALSE,
    failure_code VARCHAR(100),
    failure_message TEXT,
    outcome JSONB,
    receipt_url VARCHAR(500),
    billing_details JSONB,
    payment_method_details JSONB,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for charges
CREATE INDEX idx_charges_stripe_id ON charges(stripe_charge_id);
CREATE INDEX idx_charges_payment_intent_id ON charges(payment_intent_id);
CREATE INDEX idx_charges_status ON charges(status);
CREATE INDEX idx_charges_created_at ON charges(created_at);
