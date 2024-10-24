The system manages two specific types of packages: Internet Packages and Cable TV Packages. Each package has the following details:

Common fields for all packages:
- Unique ID
- Name
- Type (Internet or Cable TV)
- Price
- Connection Class
- Active Status (to indicate if the package is currently offered)

For Internet Packages:
- Bandwidth
- Bandwidth Type (Mbps)

For Cable TV Packages:
- Channel Count
- TV Count

-----
Common Test Cases for Both Package Types:

1. Given an admin wants to create a new package (Internet or Cable TV)
   When they provide valid package details (id, type, name, price, isActive, connectionClass, and type-specific fields)
   Then the system creates a new package and confirms its creation with the appropriate response model

2. Given an admin attempts to create a package with missing required fields
   When they submit the package details
   Then the system returns a 400 Bad Request error with specific validation messages

3. Given an admin wants to update an existing package
   When they provide valid updated details for the package
   Then the system updates the package and confirms the changes with the appropriate response model

4. Given an admin wants to delete an existing package
   When they select a package for deletion
   Then the system removes the package and confirms its deletion with a success message

5. Given a user wants to view details of a specific package
   When they request a package by ID
   Then the system displays the complete details of the selected package using the appropriate response model

6. Given a user wants to view all packages
   When they access the package listing feature
   Then the system displays a list of all available packages, using the appropriate response models

7. Given a user wants to view packages of a specific type
   When they request packages with a type parameter (Internet or Cable TV)
   Then the system displays a list of packages of the specified type, using the appropriate response models

Internet Package Specific Test Cases:

1. Given an admin wants to create a new Internet package
   When they provide valid package details including bandwidth and bandwidthType
   Then the system creates a new Internet package and confirms its creation with the InternetPackageResponse model

2. Given an admin attempts to create an Internet package without bandwidth or bandwidthType
   When they submit the package details
   Then the system returns a 400 Bad Request error with a message indicating missing required fields

Cable TV Package Specific Test Cases:

1. Given an admin wants to create a new Cable TV package
   When they provide valid package details including channelCount and tvCount
   Then the system creates a new Cable TV package and confirms its creation with the TVPackageResponse model

2. Given an admin attempts to create a Cable TV package without channelCount or tvCount
   When they submit the package details
   Then the system returns a 400 Bad Request error with a message indicating missing required fields

Error Handling Test Cases:

1. Given a user requests a non-existent package
   When they provide an invalid package ID
   Then the system returns a 404 Not Found error

2. Given an admin attempts to update a non-existent package
   When they provide an invalid package ID
   Then the system returns a 404 Not Found error

3. Given an admin attempts to delete a non-existent package
   When they provide an invalid package ID
   Then the system returns a 404 Not Found error


---
Create a Cable TV Package
```shell

```