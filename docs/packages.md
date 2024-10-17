Internet Packages' Test Cases:
1. Given an admin wants to create a new internet package
   When they provide package details (name, speed, price)
   Then a new internet package is created in the system
2. Given an admin wants to update an existing internet package
   When they provide the package ID and updated details
   Then the internet package information is updated in the system
3. Given an admin wants to delete an internet package
   When they provide the package ID
   Then the internet package is removed from the system
4. Given an admin/employee wants to view all internet packages
   When they request the list of all internet packages
   Then the system returns a list of all internet packages with their details
5. Given an admin/employee wants to view details of a specific internet package
   When they provide the package ID
   Then the system returns the details of the requested internet package

Cable TV Packages' Test Cases:
1. Given an admin wants to create a new CableTV package
   When they provide package details (name, channels, price)
   Then a new CableTV package is created in the system
2. Given an admin wants to update an existing CableTV package
   When they provide the package ID and updated details
   Then the CableTV package information is updated in the system
3. Given an admin wants to delete a CableTV package
   When they provide the package ID
   Then the CableTV package is removed from the system
4. Given an admin/employee wants to view all CableTV packages
   When they request the list of all CableTV packages
   Then the system returns a list of all CableTV packages with their details
5. Given an admin/employee wants to view details of a specific CableTV package
   When they provide the package ID
   Then the system returns the details of the requested CableTV package