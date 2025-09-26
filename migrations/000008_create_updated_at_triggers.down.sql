DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_charges_updated_at ON charges;
DROP TRIGGER IF EXISTS update_payment_methods_updated_at ON payment_methods;
DROP TRIGGER IF EXISTS update_customers_updated_at ON customers;
DROP TRIGGER IF EXISTS update_payment_intents_updated_at ON payment_intents;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();