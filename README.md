# AmazonClone_PAP

## Project Description

This project is a client-server system built in Go, designed to practice distributed computing and client-server architecture. It simulates an e-commerce platform inspired by Amazon or MercadoLibre.

The system supports two types of users: **Administrators** and **Clients**. Each connects to a TCP server and interacts through a terminal-based menu.

**Administrators** can:
- Add new products to the inventory
- Update product stock and price
- View all orders placed in the system
- Complete (fulfill) orders
- List all registered users

**Clients** can:
- Browse available products
- Add and remove products from their cart
- View their cart
- Place orders
- View their order history
- Cancel pending orders

The system includes **user authentication** (login and registration) and a **save state feature** that persists products, orders, and users to JSON files, so data survives server restarts. All activity is logged to a file for traceability.

---

## Instructions

### Prerequisites
- Go 1.18 or higher installed. You can download it from [https://golang.org/dl/](https://golang.org/dl/)
- Git installed

### Download the project
```bash
git clone https://github.com/DiegoAmin/AmazonClone_PAP.git
cd AmazonClone_PAP
```

### Install dependencies
```bash
go mod tidy
```

### Run the server
Open a terminal and run:
```bash
go run cmd/server/server.go
```
The server will start listening on `localhost:10000`. On the first run, it will create `store.json` and `users.json` with default data.

### Connect a client
Open a **new terminal** and run:
```bash
go run cmd/client/client.go
```
You can open multiple terminals to simulate multiple clients connected at the same time.

### Default credentials
| Username | Password   | Role     |
|----------|------------|----------|
| admin    | admin123   | admin    |
| carlos   | carlos123  | customer |
| diego    | diego123   | customer |
| diego2   | diego2123  | customer |

---

## Examples of Administrative Functionalities

### Starting the server
```
$ go run cmd/server/server.go
TCP server listening on localhost:10000
```

### Connecting and logging in as admin
```
$ go run cmd/client/client.go
Connected to server!
Welcome to the Amazon Clone Server!
1. Login
2. Register
Please enter your choice:
1
Enter username:
admin
Enter password:
admin123
Login successful! Welcome, admin!
Welcome to Admin Mode!
1. Add Product
2. Update Stock
3. Update Price
4. Get orders report
5. Complete order
6. List users
7. Exit
Please enter your choice:
```

### 1. Add a product
```
Please enter your choice:
1
Enter product details:
ID:
4
Name:
Mouse
Price:
300
Stock:
25
Product added successfully!
```

### 2. Update stock
```
Please enter your choice:
2
Enter product ID:
4
Enter new stock quantity:
0
Stock updated successfully!
```

### 3. Update price
```
Please enter your choice:
3
Enter product ID:
4
Enter new price:
400
Price updated successfully!
```

### 4. Get orders report
```
Please enter your choice:
4
Orders Report:
Order ID: 1 | Client: carlos | Status: CREATED | Total: 2500.00
  -> Laptop x1 | Price at order: 1000.00 | Subtotal: 1000.00
  -> Desktop x3 | Price at order: 500.00 | Subtotal: 1500.00
```

### 5. Complete an order
```
Please enter your choice:
5
Enter order ID to complete:
1
Order completed successfully!
```
After completing, the order status changes to `COMPLETED`:
```
Orders Report:
Order ID: 1 | Client: carlos | Status: COMPLETED | Total: 2500.00
  -> Laptop x1 | Price at order: 1000.00 | Subtotal: 1000.00
  -> Desktop x3 | Price at order: 500.00 | Subtotal: 1500.00
```

### 6. List users
```
Please enter your choice:
6
Registered Users:
Username: carlos | Role: customer
Username: diego | Role: customer
Username: diego2 | Role: customer
Username: admin | Role: admin
```

### 7. Exit
```
Please enter your choice:
7
Exiting Admin Mode. Goodbye!
```

---

## Examples of Client Functionalities

### Connecting and logging in as client
```
$ go run cmd/client/client.go
Connected to server!
Welcome to the Amazon Clone Server!
1. Login
2. Register
Please enter your choice:
1
Enter username:
carlos
Enter password:
carlos123
Login successful! Welcome, carlos!
Welcome to Client Mode!
1. See list of products
2. Add product to cart
3. Remove product from cart
4. View products in cart
5. Place order
6. View all orders
7. Cancel order
8. Exit
Please enter your choice:
```

### 1. See list of products
```
Please enter your choice:
1
Available Products:
ID: 1 | Name: Laptop | Price: 1000.00 | Stock: 9
ID: 2 | Name: Desktop | Price: 500.00 | Stock: 27
ID: 3 | Name: TV32in | Price: 3000.00 | Stock: 40
ID: 4 | Name: Mouse | Price: 400.00 | Stock: 0
```

### 2. Add product to cart
```
Please enter your choice:
2
Enter product ID:
3
Enter quantity:
5
Product TV32in added to cart! Quantity: 5
```

### 3. Remove product from cart
```
Please enter your choice:
3
Enter product ID to remove:
3
Enter quantity to remove (current: 5):
1
Updated cart quantity: 4
```

### 4. View products in cart
```
Please enter your choice:
4
Your cart:
ID: 3 | Name: TV32in | Quantity: 4 | Subtotal: 12000.00
Cart Total: 12000.00
```

### 5. Place order
```
Please enter your choice:
5
Order placed successfully! Order ID: 2 | Total: 12000.00
```

### 6. View all orders
```
Please enter your choice:
6
Orders:
Order ID: 1 | Status: COMPLETED | Total: 2500.00
  -> Laptop x1 | Price at order: 1000.00 | Subtotal: 1000.00
  -> Desktop x3 | Price at order: 500.00 | Subtotal: 1500.00
Order ID: 2 | Status: CREATED | Total: 12000.00
  -> TV32in x4 | Price at order: 3000.00 | Subtotal: 12000.00
```

### 7. Cancel order
```
Please enter your choice:
7
Enter order ID to cancel:
2
Order cancelled successfully!
```
After cancelling, the order status changes to `CANCELLED`:
```
Orders:
Order ID: 2 | Status: CANCELLED | Total: 12000.00
  -> TV32in x4 | Price at order: 3000.00 | Subtotal: 12000.00
Order ID: 1 | Status: COMPLETED | Total: 2500.00
  -> Laptop x1 | Price at order: 1000.00 | Subtotal: 1000.00
  -> Desktop x3 | Price at order: 500.00 | Subtotal: 1500.00
```

### 8. Exit
```
Please enter your choice:
8
Exiting Client Mode. Goodbye!
```

---

### Registering a new user
```
Welcome to the Amazon Clone Server!
1. Login
2. Register
Please enter your choice:
2
Enter username:
newuser
Enter password:
newpass123
Enter role (admin/customer):
customer
Registration successful!
You can now login with your new credentials.
```

---

## Features for Future Work

1. **Graphical User Interface (GUI):** Replace the terminal-based interaction with a web or desktop UI, making the system more accessible and user-friendly for non-technical users.

2. **Product Blocking in Cart:** When a client adds a product to their cart, the item gets temporarily blocked so no other client can add it to their cart — preventing stock conflicts between concurrent users.

3. **Password Encryption and Token-based Authentication:** Currently passwords are stored in plain text. A future improvement would be to hash passwords using bcrypt and implement JWT tokens for session management, making the system more secure.