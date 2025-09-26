DROP INDEX IF EXISTS idx_payment_intents_created_at;
DROP INDEX IF EXISTS idx_payment_intents_status;
DROP INDEX IF EXISTS idx_payment_intents_customer_id;
DROP INDEX IF EXISTS idx_payment_intents_stripe_id;

DROP TABLE IF EXISTS payment_intents;
