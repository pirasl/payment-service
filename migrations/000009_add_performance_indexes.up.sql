-- Composite indexes for common queries
CREATE INDEX idx_payment_intents_customer_status ON payment_intents(customer_id, status);
CREATE INDEX idx_payment_intents_status_created ON payment_intents(status, created_at);
CREATE INDEX idx_charges_payment_intent_status ON charges(payment_intent_id, status);
CREATE INDEX idx_webhook_events_type_processed ON webhook_events(event_type, processed);

-- Partial indexes for specific use cases
CREATE INDEX idx_unprocessed_webhooks ON webhook_events(created_at) 
WHERE processed = FALSE;

CREATE INDEX idx_failed_payment_intents ON payment_intents(created_at) 
WHERE status IN ('requires_payment_method', 'canceled');