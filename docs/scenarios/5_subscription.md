1. Given an employee wants to add a new subscription for a customer
   When they provide valid subscription details (CustomerID, Type, PackageID)
   Then a new subscription record is created in the system
2. Given an employee attempts to add a new subscription
   When they submit incomplete subscription details (missing CustomerID)
   Then the system informs them that CustomerID is required
3. Given an employee attempts to add a new subscription
   When they provide an invalid PackageID
   Then the system informs them that the package ID is invalid
4. Given an employee attempts to add a new subscription
   When they provide an invalid subscription type
   Then the system informs them that the subscription type is invalid
5. Given an employee wants to view details of a specific subscription
   When they provide a valid subscription ID
   Then the system returns the complete details of the requested subscription
6. Given an employee wants to update subscription information
   When they provide the subscription ID and valid updated details
   Then the subscription information is updated in the system
7. Given an employee wants to remove a subscription from the system
   When they provide the ID of an existing subscription
   Then the subscription is deleted from the system
8. Given an employee wants to view all subscriptions
   When they request the list of all subscriptions
   Then the system returns a list of all subscriptions with their details
9. Given an employee wants to assign a device to a subscription
   When they provide a valid device ID and subscription ID
   Then the device is associated with the specified subscription