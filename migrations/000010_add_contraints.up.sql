ALTER TABLE payment_intents 
ADD CONSTRAINT chk_payment_intents_amount_positive CHECK (amount > 0),
ADD CONSTRAINT chk_payment_intents_currency_length CHECK (length(currency) = 3),
ADD CONSTRAINT chk_payment_intents_status_valid CHECK (
    status IN (
        'requires_payment_method', 'requires_confirmation', 'requires_action',
        'processing', 'requires_capture', 'canceled', 'succeeded'
    )
);

ALTER TABLE charges 
ADD CONSTRAINT chk_charges_amount_positive CHECK (amount > 0),
ADD CONSTRAINT chk_charges_amount_captured_valid CHECK (amount_captured >= 0 AND amount_captured <= amount),
ADD CONSTRAINT chk_charges_amount_refunded_valid CHECK (amount_refunded >= 0 AND amount_refunded <= amount);

ALTER TABLE refunds 
ADD CONSTRAINT chk_refunds_amount_positive CHECK (amount > 0),
ADD CONSTRAINT chk_refunds_reason_valid CHECK (
    reason IS NULL OR reason IN ('duplicate', 'fraudulent', 'requested_by_customer')
);


CREATE UNIQUE INDEX idx_unique_default_payment_method 
ON payment_methods(customer_id) 
WHERE is_default = TRUE;