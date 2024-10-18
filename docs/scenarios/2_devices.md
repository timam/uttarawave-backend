1. Given an employee wants to add a new device to the system
   When they provide valid device details
   Then a new device record is created in the system
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
   When they request the list of all devices
   Then the system returns a list of all devices with their details
6. Given an employee wants to assign a device to a customer's subscription
   When they provide a valid device ID and subscription ID
   Then the device is associated with the specified subscription
7. Given an employee wants to assign a device to a building
   When they provide a valid device ID and building ID
   Then the device is associated with the specified building
8. Given an employee wants to unassign a device from its current association
   When they provide a valid device ID
   Then the device is unassigned from its current subscription or building