1. Given an employee wants to process a cash payment for a subscription
   When they provide valid transaction details (SubscriptionID and Amount)
   Then a new transaction record is created and the subscription is updated accordingly
2. Given an employee attempts to process a cash payment
   When they provide an invalid SubscriptionID
   Then the system informs them that the subscription couldn't be found
3. Given an employee processes a cash payment for a subscription
   When the payment amount is equal to or greater than the total due amount
   Then the subscription's PaidUntil date is updated and DueAmount is set to 0
4. Given an employee processes a cash payment for a subscription
   When the payment amount is less than the total due amount
   Then the subscription's DueAmount is updated accordingly
5. Given an employee processes a cash payment for a subscription
   When the transaction is completed successfully
   Then the subscription's status is updated to "Active"
6. Given an employee wants to view transactions for a specific subscription
   When they provide a valid SubscriptionID
   Then the system returns a list of all transactions associated with that subscription
7. Given an employee needs to update a transaction's details
   When they provide valid updated information for an existing transaction
   Then the transaction record is updated in the system
