1. Given an employee wants to add a new customer to the system
   When they provide valid customer details (mobile, name, and address fields)
   Then a new customer record is created in the system
2. Given an employee wants to add a new customer to an existing building
   When they provide valid customer details, building ID, and flat number
   Then a new customer record is created with the building's address details
3. Given an employee attempts to add a new customer
   When they submit incomplete customer details (missing mobile or name)
   Then the system informs them that mobile and name are required fields
4. Given an employee wants to view details of a specific customer
   When they provide a valid mobile number
   Then the system returns the complete details of the requested customer
5. Given an employee wants to view a list of customers
   When they request the list with optional page and pageSize parameters
   Then the system returns a paginated list of customers with basic details
6. Given an employee wants to view a detailed list of all customers
   When they request the full details list with pagination
   Then the system returns a paginated list of customers with their subscriptions, devices, and building information
7. Given an employee wants to update a customer's information
   When they provide the customer ID and valid updated details
   Then the customer's information is updated in the system
8. Given an employee attempts to update a customer's mobile number
   When they provide a new mobile number in the update request
   Then the system informs them that the mobile number cannot be updated
9. Given an employee wants to remove a customer from the system
   When they provide the customer's mobile number
   Then the customer is deleted from the system
10.Given an employee attempts to delete a customer
   When they provide a mobile number that doesn't exist in the system
   Then the system informs them that the customer was not found
