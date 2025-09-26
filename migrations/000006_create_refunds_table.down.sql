DROP INDEX IF EXISTS idx_refunds_status;
DROP INDEX IF EXISTS idx_refunds_payment_intent_id;
DROP INDEX IF EXISTS idx_refunds_charge_id;
DROP INDEX IF EXISTS idx_refunds_stripe_id;

-- Drop the refunds table
DROP TABLE IF EXISTS refunds;
