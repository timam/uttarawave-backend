The device management system allows for the following operations:

1. Create a new device
2. Retrieve a device by ID
3. Update device information
4. Delete a device
5. List all devices with pagination
6. Assign a device to a subscription or building
7. Unassign a device
8. Retrieve a device by its assignment (subscription or building)
9. Mark a device's status

-----

Test Cases:
1. Given an employee wants to add a new device to the system
   When they provide valid device details
   Then a new device record is created in the system with status "InStock"

2. Given an employee attempts to add a new device
   When they submit incomplete device details (missing brand)
   Then the system informs them that the brand is required

3. Given an employee wants to view details of a specific device
   When they provide a valid device ID
   Then the system returns the complete details of the requested device

4. Given an employee wants to update device information
   When they provide the device ID and valid updated details
   Then the device information is updated in the system

5. Given an employee wants to view all devices
   When they request the list of all devices with pagination parameters
   Then the system returns a paginated list of devices with their details

6. Given an employee wants to assign a device to a customer's subscription
   When they provide a valid device ID, assignment type "Subscription", and subscription ID
   Then the device is associated with the specified subscription and its status is set to "Assigned"

7. Given an employee wants to assign a device to a building
   When they provide a valid device ID, assignment type "Building", and building ID
   Then the device is associated with the specified building and its status is set to "Assigned"

8. Given an employee wants to unassign a device from its current association
   When they provide a valid device ID
   Then the device is unassigned from its current subscription or building and its status is set to "InStock"

9. Given an employee wants to retrieve a device by its assignment
   When they provide a valid assignment type and assignment ID
   Then the system returns the device associated with the specified assignment

10. Given an employee wants to mark a device with a specific status
    When they provide a valid device ID and a valid status
    Then the device's status is updated in the system

11. Given an employee wants to view all devices
    When they request the list of all devices without specifying pagination parameters
    Then the system returns the first page of devices with a default page size