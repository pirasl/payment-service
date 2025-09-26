DROP INDEX IF EXISTS idx_failed_payment_intents;
DROP INDEX IF EXISTS idx_unprocessed_webhooks;

-- Drop composite indexes
DROP INDEX IF EXISTS idx_webhook_events_type_processed;
DROP INDEX IF EXISTS idx_charges_payment_intent_status;
DROP INDEX IF EXISTS idx_payment_intents_status_created;
DROP INDEX IF EXISTS idx_payment_intents_customer_status;
