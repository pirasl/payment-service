DROP INDEX IF EXISTS idx_payment_methods_type;
DROP INDEX IF EXISTS idx_payment_methods_customer_id;
DROP INDEX IF EXISTS idx_payment_methods_stripe_id;

-- Drop the payment_methods table
DROP TABLE IF EXISTS payment_methods;
