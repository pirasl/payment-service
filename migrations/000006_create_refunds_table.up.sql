CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stripe_refund_id VARCHAR(255) UNIQUE NOT NULL,
    charge_id UUID REFERENCES charges(id) ON DELETE CASCADE,
    payment_intent_id UUID REFERENCES payment_intents(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    reason VARCHAR(50), -- duplicate, fraudulent, requested_by_customer
    status VARCHAR(50) NOT NULL,
    failure_reason VARCHAR(100),
    receipt_number VARCHAR(100),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for refunds
CREATE INDEX idx_refunds_stripe_id ON refunds(stripe_refund_id);
CREATE INDEX idx_refunds_charge_id ON refunds(charge_id);
CREATE INDEX idx_refunds_payment_intent_id ON refunds(payment_intent_id);
CREATE INDEX idx_refunds_status ON refunds(status);
