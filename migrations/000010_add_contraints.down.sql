DROP INDEX IF EXISTS idx_unique_default_payment_method;

-- Remove check constraints from refunds
ALTER TABLE refunds 
DROP CONSTRAINT IF EXISTS chk_refunds_amount_positive,
DROP CONSTRAINT IF EXISTS chk_refunds_reason_valid;

-- Remove check constraints from charges
ALTER TABLE charges 
DROP CONSTRAINT IF EXISTS chk_charges_amount_positive,
DROP CONSTRAINT IF EXISTS chk_charges_amount_captured_valid,
DROP CONSTRAINT IF EXISTS chk_charges_amount_refunded_valid;

-- Remove check constraints from payment_intents
ALTER TABLE payment_intents 
DROP CONSTRAINT IF EXISTS chk_payment_intents_amount_positive,
DROP CONSTRAINT IF EXISTS chk_payment_intents_currency_length,
DROP CONSTRAINT IF EXISTS chk_payment_intents_status_valid;