1. Given an admin or employee wants to add a new building
   When they provide valid building details including CableTV connection or Internet connection
   Then a new building record is created in the system with the connection information
2. Given an admin wants to remove a building from the system
   When they provide the ID of an existing building
   Then the building is deleted from the system
3. Given an employee attempts to delete a building
   When they try to access the delete building endpoint
   Then the system informs them they don't have permission to perform this action
4. Given an admin or employee wants to view buildings sorted by area
   When they request the list of buildings with area sorting parameter
   Then the system returns a list of buildings sorted by area
5. Given an admin or employee wants to find buildings in a specific area
   When they provide an area name as a search parameter
   Then the system returns a list of buildings in that area

