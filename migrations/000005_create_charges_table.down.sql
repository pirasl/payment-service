DROP INDEX IF EXISTS idx_charges_created_at;
DROP INDEX IF EXISTS idx_charges_status;
DROP INDEX IF EXISTS idx_charges_payment_intent_id;
DROP INDEX IF EXISTS idx_charges_stripe_id;

-- Drop the charges table
DROP TABLE IF EXISTS charges;