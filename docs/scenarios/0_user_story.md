# System Overview

We have implemented a system that manages customers, subscriptions, invoices, payments, devices, buildings, and packages for a service provider (likely an internet or cable TV service). Here's a breakdown of the key components and their interactions:

## Key Components

### Customers
These are the end-users of the service. Customers can be associated with a building or have their own address details. Each customer has a unique ID, mobile number, name, and address information.

### Subscriptions
Customers can have one or more subscriptions to different services (Internet or Cable TV). Each subscription is associated with a specific package and has details like the package type, price, due amount, and renewal date.

### Invoices
Invoices are generated for the services provided to customers. Invoices can be created automatically when a subscription is due for renewal or manually when there's an update to the subscription that increases the due amount. Each invoice has a status (pending or paid), amount due, and due date.

### Payments
Customers make payments against the invoices. Payments can be full or partial, and they update the subscription's due amount and the invoice's status accordingly.

### Devices
Devices (e.g., cable boxes, modems) can be associated with a customer's subscription or a building. Each device has details like brand, model, serial number, type, usage, and status.

### Buildings
Buildings can have internet or cable TV connections. Buildings can have devices associated with them, and customers can be linked to a specific building and flat number.

### Packages
The system manages two types of packages: Internet Packages and Cable TV Packages. Each package has details like name, speed, price, connection type, bandwidth type, and channel lineup.

## Flow of Operations

### Customer Management
Employees can create, update, view, and delete customer records. Customers can be associated with a building or have their own address details.

### Subscription Management
Employees can create, update, view, and delete subscriptions for customers. Subscriptions are associated with a specific package and have details like the package type, price, due amount, and renewal date.

### Invoice Generation and Payment Processing
Invoices are generated automatically when a subscription is due for renewal or manually when there's an update to the subscription that increases the due amount. Customers make payments against the invoices, which update the subscription's due amount and the invoice's status accordingly.

### Device Management
Employees can create, update, view, and delete device records. Devices can be associated with a customer's subscription or a building.

### Building Management
Admins and employees can create, view, and delete building records. Buildings can have internet or cable TV connections and can have devices associated with them.

### Package Management
Admins can create, update, view, and delete internet and cable TV packages. Employees can only view package details.

### Handling Overdue Payments and Expired Subscriptions
If a subscription's renewal date passes without full payment, it may be marked as expired. For subscriptions with devices (e.g., onu, switch), the system can mark the device for collection if the subscription expires.

## Key Points to Note
- Invoices are generated both automatically (for renewals) and manually (for mid-cycle changes).
- Subscriptions can have pending due amounts if payments are partial or missed.
- The system allows for flexible payment handling, including partial payments.
- There's a mechanism to handle expired subscriptions and associated devices.
- Customers can be associated with buildings or have their own address details.
- Devices can be associated with subscriptions or buildings.
- Packages are managed separately for internet and cable TV services.
- Employees have different levels of access for various operations.