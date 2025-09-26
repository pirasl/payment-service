ALTER TABLE payment_intents 
ADD CONSTRAINT fk_payment_intents_customer_id 
FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL;
