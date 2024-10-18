Internet Packages' Test Cases:
1. Given an admin wants to create a new internet package
   When they provide valid package details(name, speed, price, connection type, bandwidth type)
   Then the system creates a new internet package and confirms its creation

2. Given an employee attempts to create a new internet package
   When they try to access the package creation feature
   Then the system informs them they don't have permission to perform this action
   
3. Given an admin attempts to create an internet package with invalid data
   When they submit the package details
   Then the system informs them of the specific validation errors
   
4. Given an admin wants to update an existing internet package
   When they provide valid updated details for the package
   Then the system updates the package and confirms the changes
   
5. Given an employee attempts to update an internet package
   When they try to access the package update feature
   Then the system informs them they don't have permission to perform this action
   
6. Given an admin wants to delete an existing internet package
   When they select a package for deletion
   Then the system removes the package and confirms its deletion
   
7. Given an employee attempts to delete an internet package
   When they try to access the package deletion feature
   Then the system informs them they don't have permission to perform this action
   
8. Given an admin or employee wants to view details of a specific internet package
   When they select a package to view
   Then the system displays the complete details of the selected package
   
9. Given an admin or employee wants to view all internet packages
   When they access the package listing feature
   Then the system displays a list of all available internet packages




Cable TV Packages' Test Cases:
1. Given an admin wants to create a new CableTV package
   When they provide valid package details
   Then the system creates a new CableTV package and confirms its creation

2. Given an employee attempts to create a new CableTV package
   When they try to access the package creation feature
   Then the system informs them they don't have permission to perform this action

3. Given an admin attempts to create a CableTV package with invalid data
   When they submit the package details
   Then the system informs them of the specific validation errors
   
4. Given an admin wants to update an existing CableTV package
   When they provide valid updated details for the package
   Then the system updates the package and confirms the changes
   
5. Given an employee attempts to update a CableTV package
   When they try to access the package update feature
   Then the system informs them they don't have permission to perform this action
   
6. Given an admin wants to delete an existing CableTV package
   When they select a package for deletion 
   Then the system removes the package and confirms its deletion
   
7. Given an employee attempts to delete a CableTV package
   When they try to access the package deletion feature
   Then the system informs them they don't have permission to perform this action
   
8. Given an admin or employee wants to view details of a specific CableTV package
   When they select a package to view
   Then the system displays the complete details of the selected package
   
9. Given an admin or employee wants to view all CableTV packages
   When they access the package listing feature
   Then the system displays a list of all available CableTV packages