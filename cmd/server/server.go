package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/DiegoAmin/AmazonClone_PAP/internal/auth"
	"github.com/DiegoAmin/AmazonClone_PAP/internal/logger"
	"github.com/DiegoAmin/AmazonClone_PAP/internal/order"
	"github.com/DiegoAmin/AmazonClone_PAP/internal/product"
	"github.com/DiegoAmin/AmazonClone_PAP/internal/store"
)

func main() {
	// Initialize the store
	store, err := store.NewStore()
	if err != nil {
		panic(err)
	}

	//Initialize the auth store
	authStore := auth.NewAuthStore()

	// Add some basic products to the store
	p1, _ := product.NewProduct(1, "Laptop", 1000, 10)
	store.AddProduct(*p1)
	p2, _ := product.NewProduct(2, "Desktop", 500, 30)
	store.AddProduct(*p2)
	p3, _ := product.NewProduct(3, "TV32in", 3000, 40)
	store.AddProduct(*p3)

	// Open a port to listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:10000")
	if err != nil {
		panic(err)
	}

	// Print a message to indicate that the server is listening for incoming connections
	fmt.Println("TCP server listening on localhost:10000")

	// Accept incoming connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConn(conn, store, authStore)
	}
}

// handleConn handles the incoming connection from the client. It prompts the user to choose between login and register, and then calls the appropriate handler function based on the user's choice.
func handleConn(c net.Conn, store *store.Store, authStore *auth.AuthStore) {
	defer c.Close()

	choice := bufio.NewScanner(c)

	for {
		fmt.Fprintf(c, "Welcome to the Amazon Clone Server!\n")
		fmt.Fprintf(c, "1. Login\n")
		fmt.Fprintf(c, "2. Register\n")
		fmt.Fprintf(c, "Please enter your choice:\n")

		if choice.Scan() {
			switch choice.Text() {
			case "1":
				fmt.Fprintf(c, "Enter username:\n")
				choice.Scan()
				username := choice.Text()
				fmt.Fprintf(c, "Enter password:\n")
				choice.Scan()
				password := choice.Text()
				user, err := authStore.Login(username, password)
				if err != nil {
					fmt.Fprintf(c, "Login failed: %s\n", err.Error())
					continue
				}
				fmt.Fprintf(c, "Login successful! Welcome, %s!\n", user.Username)
				if user.Role == "admin" {
					handleAdmin(c, choice, store, user.Username, authStore)
				} else {
					handleClient(c, choice, store, user.Username)
				}
				return
			case "2":
				fmt.Fprintf(c, "Enter username:\n")
				choice.Scan()
				username := choice.Text()
				fmt.Fprintf(c, "Enter password:\n")
				choice.Scan()
				password := choice.Text()
				fmt.Fprintf(c, "Enter role (admin/customer):\n")
				choice.Scan()
				role := choice.Text()
				if role != "admin" && role != "customer" {
					fmt.Fprintf(c, "Invalid role. Please enter 'admin' or 'customer'.\n")
					continue
				}
				err := authStore.Register(username, password, role)
				if err != nil {
					fmt.Fprintf(c, "Registration failed: %s\n", err.Error())
					continue
				}
				fmt.Fprintf(c, "Registration successful!\n")
				fmt.Fprintf(c, "You can now login with your new credentials.\n")
			default:
				fmt.Fprintf(c, "Invalid choice. Please enter 1 or 2.\n")
			}
		}
	}
}

// handleAdmin handles the admin mode of the server, allowing the admin to add products, update stock and price, complete orders, and get orders report.
func handleAdmin(c net.Conn, choice *bufio.Scanner, store *store.Store, adminUsername string, authStore *auth.AuthStore) {
	fmt.Fprintf(c, "Welcome to Admin Mode!\n")
	logger.Log(fmt.Sprintf("ADMIN: %s logged in", adminUsername))
	for {
		fmt.Fprintf(c, "1. Add Product\n")
		fmt.Fprintf(c, "2. Update Stock\n")
		fmt.Fprintf(c, "3. Update Price\n")
		fmt.Fprintf(c, "4. Get orders report\n")
		fmt.Fprintf(c, "5. Complete order\n")
		fmt.Fprintf(c, "6. List users\n")
		fmt.Fprintf(c, "7. Exit\n")
		fmt.Fprintf(c, "Please enter your choice:\n")

		if choice.Scan() {
			switch choice.Text() {
			case "1": // Add Product
				fmt.Fprintf(c, "Enter product details:\n")
				fmt.Fprintf(c, "ID:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				fmt.Fprintf(c, "Name:\n")
				choice.Scan()
				name := choice.Text()
				fmt.Fprintf(c, "Price:\n")
				choice.Scan()
				price, err := strconv.ParseFloat(choice.Text(), 64)
				if err != nil {
					fmt.Fprintf(c, "Invalid price. Please enter a valid number.\n")
					continue
				}
				fmt.Fprintf(c, "Stock:\n")
				choice.Scan()
				stock, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid stock. Please enter a valid integer.\n")
					continue
				}
				p, err := product.NewProduct(id, name, price, stock)
				if err != nil {
					fmt.Fprintf(c, "Error creating product: %s\n", err.Error())
					continue
				}
				err = store.AddProduct(*p)
				if err != nil {
					fmt.Fprintf(c, "Error adding product: %s\n", err.Error())
				} else {
					fmt.Fprintf(c, "Product added successfully!\n")
				}

			case "2": // Update Stock
				fmt.Fprintf(c, "Enter product ID:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				fmt.Fprintf(c, "Enter new stock quantity:\n")
				choice.Scan()
				stock, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid stock. Please enter a valid integer.\n")
					continue
				}
				err = store.UpdateStock(id, stock)
				if err != nil {
					fmt.Fprintf(c, "Error updating stock: %s\n", err.Error())
				} else {
					fmt.Fprintf(c, "Stock updated successfully!\n")
				}

			case "3": // Update Price
				fmt.Fprintf(c, "Enter product ID:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				fmt.Fprintf(c, "Enter new price:\n")
				choice.Scan()
				price, err := strconv.ParseFloat(choice.Text(), 64)
				if err != nil {
					fmt.Fprintf(c, "Invalid price. Please enter a valid number.\n")
					continue
				}
				err = store.UpdatePrice(id, price)
				if err != nil {
					fmt.Fprintf(c, "Error updating price: %s\n", err.Error())
				} else {
					fmt.Fprintf(c, "Price updated successfully!\n")
				}

			case "4": // Get orders report
				ordersReport := store.OrderHistory()
				if len(ordersReport) == 0 {
					fmt.Fprintf(c, "No orders found.\n")
					continue
				}
				fmt.Fprintf(c, "Orders Report:\n")
				for _, o := range ordersReport {
					fmt.Fprintf(c, "Order ID: %d | Client: %s | Status: %s | Total: %.2f\n", o.ID, o.Username, o.Status, o.Total)
					for _, item := range o.Items {
						// Use frozen price from OrderItem, not current product price
						// GetProduct is used only to display the name, price comes from the frozen item
						if p, exists := store.GetProduct(item.ProductID); exists {
							fmt.Fprintf(c, "  -> %s x%d | Price at order: %.2f | Subtotal: %.2f\n", p.Name, item.Quantity, item.Price, item.Price*float64(item.Quantity))
						}
					}
				}

			case "5": // Complete order
				fmt.Fprintf(c, "Enter order ID to complete:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				err = store.CompleteOrder(id)
				if err != nil {
					fmt.Fprintf(c, "Error completing order: %s\n", err.Error())
				} else {
					fmt.Fprintf(c, "Order completed successfully!\n")
				}
			case "6": // List users
				users := authStore.ListUsers()
				if len(users) == 0 {
					fmt.Fprintf(c, "No users found.\n")
					continue
				}
				fmt.Fprintf(c, "Registered Users:\n")
				for _, user := range users {
					fmt.Fprintf(c, "Username: %s | Role: %s\n", user.Username, user.Role)
				}

			case "7": // Exit
				fmt.Fprintf(c, "Exiting Admin Mode. Goodbye!\n")
				logger.Log(fmt.Sprintf("ADMIN: %s logged out", adminUsername))
				fmt.Fprintf(c, "EXIT\n")
				return

			default:
				fmt.Fprintf(c, "Invalid choice. Please enter a number between 1 and 7.\n")
				continue
			}
		}
	}
}

// handleClient handles the client mode of the server, allowing the client to view products, manage cart, place orders, view orders, and cancel orders.
func handleClient(c net.Conn, choice *bufio.Scanner, store *store.Store, username string) {
	// cart is a map where the key is the product ID and the value is the quantity.
	cart := make(map[int]int)
	fmt.Fprintf(c, "Welcome to Client Mode!\n")
	logger.Log(fmt.Sprintf("CLIENT: %s logged in", username))
	for {
		fmt.Fprintf(c, "1. See list of products\n")
		fmt.Fprintf(c, "2. Add product to cart\n")
		fmt.Fprintf(c, "3. Remove product from cart\n")
		fmt.Fprintf(c, "4. View products in cart\n")
		fmt.Fprintf(c, "5. Place order\n")
		fmt.Fprintf(c, "6. View all orders\n")
		fmt.Fprintf(c, "7. Cancel order\n")
		fmt.Fprintf(c, "8. Exit\n")
		fmt.Fprintf(c, "Please enter your choice:\n")

		if choice.Scan() {
			switch choice.Text() {
			case "1": // See list of products
				products := store.ListProducts()
				if len(products) == 0 {
					fmt.Fprintf(c, "No products available.\n")
					continue
				}
				fmt.Fprintf(c, "Available Products:\n")
				for _, p := range products {
					fmt.Fprintf(c, "ID: %d | Name: %s | Price: %.2f | Stock: %d\n", p.ID, p.Name, p.Price, p.Stock)
				}

			case "2": // Add product to cart
				fmt.Fprintf(c, "Enter product ID:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				p, exists := store.GetProduct(id)
				if !exists {
					fmt.Fprintf(c, "Product with ID %d not found.\n", id)
					continue
				}
				fmt.Fprintf(c, "Enter quantity:\n")
				choice.Scan()
				quantity, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid quantity. Please enter a valid integer.\n")
					continue
				}
				if quantity <= 0 {
					fmt.Fprintf(c, "Quantity must be greater than 0.\n")
					continue
				}
				if p.Stock < quantity {
					fmt.Fprintf(c, "Not enough stock. Available: %d\n", p.Stock)
					continue
				}
				cart[id] += quantity
				fmt.Fprintf(c, "Product %s added to cart! Quantity: %d\n", p.Name, cart[id])

			case "3": // Remove product from cart
				fmt.Fprintf(c, "Enter product ID to remove:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				if _, exists := cart[id]; !exists {
					fmt.Fprintf(c, "Product with ID %d is not in your cart.\n", id)
					continue
				}
				fmt.Fprintf(c, "Enter quantity to remove (current: %d):\n", cart[id])
				choice.Scan()
				quantity, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid quantity. Please enter a valid integer.\n")
					continue
				}
				if quantity <= 0 {
					fmt.Fprintf(c, "Quantity must be greater than 0.\n")
					continue
				}
				if quantity > cart[id] {
					fmt.Fprintf(c, "Cannot remove more than you have in cart. Current quantity: %d\n", cart[id])
					continue
				}
				cart[id] -= quantity
				if cart[id] == 0 {
					delete(cart, id)
					fmt.Fprintf(c, "Product removed from cart completely!\n")
				} else {
					fmt.Fprintf(c, "Updated cart quantity: %d\n", cart[id])
				}

			case "4": // View products in cart
				if len(cart) == 0 {
					fmt.Fprintf(c, "Your cart is empty.\n")
					continue
				}
				fmt.Fprintf(c, "Your cart:\n")
				total := 0.0
				for id, quantity := range cart {
					p, exists := store.GetProduct(id)
					if !exists {
						continue
					}
					subtotal := p.Price * float64(quantity)
					total += subtotal
					fmt.Fprintf(c, "ID: %d | Name: %s | Quantity: %d | Subtotal: %.2f\n", id, p.Name, quantity, subtotal)
				}
				fmt.Fprintf(c, "Cart Total: %.2f\n", total)

			case "5": // Place order
				if len(cart) == 0 {
					fmt.Fprintf(c, "Your cart is empty.\n")
					continue
				}
				items := []order.OrderItem{}
				for id, quantity := range cart {
					items = append(items, order.OrderItem{ProductID: id, Quantity: quantity})
				}
				newOrder, err := store.CreateOrder(items, username)
				if err != nil {
					fmt.Fprintf(c, "Error placing order: %s\n", err.Error())
					continue
				}
				cart = make(map[int]int) // clear cart after order
				fmt.Fprintf(c, "Order placed successfully! Order ID: %d | Total: %.2f\n", newOrder.ID, newOrder.Total)

			case "6": // View all orders
				ordersReport := store.OrderHistoryByUser(username)
				if len(ordersReport) == 0 {
					fmt.Fprintf(c, "No orders found.\n")
					continue
				}
				fmt.Fprintf(c, "Orders:\n")
				for _, o := range ordersReport {
					fmt.Fprintf(c, "Order ID: %d | Status: %s | Total: %.2f\n", o.ID, o.Status, o.Total)
					for _, item := range o.Items {
						// Use frozen price from OrderItem, not current product price
						// GetProduct is used only to display the name, price comes from the frozen item
						if p, exists := store.GetProduct(item.ProductID); exists {
							fmt.Fprintf(c, "  -> %s x%d | Price at order: %.2f | Subtotal: %.2f\n", p.Name, item.Quantity, item.Price, item.Price*float64(item.Quantity))
						}
					}
				}

			case "7": // Cancel order
				fmt.Fprintf(c, "Enter order ID to cancel:\n")
				choice.Scan()
				id, err := strconv.Atoi(choice.Text())
				if err != nil {
					fmt.Fprintf(c, "Invalid ID. Please enter a valid integer.\n")
					continue
				}
				err = store.CancelOrder(id, username)
				if err != nil {
					fmt.Fprintf(c, "Error cancelling order: %s\n", err.Error())
				} else {
					fmt.Fprintf(c, "Order cancelled successfully!\n")
				}

			case "8": // Exit
				fmt.Fprintf(c, "Exiting Client Mode. Goodbye!\n")
				logger.Log(fmt.Sprintf("CLIENT: %s logged out", username))
				fmt.Fprintf(c, "EXIT\n")
				return

			default:
				fmt.Fprintf(c, "Invalid choice. Please enter a number between 1 and 8.\n")
				continue
			}
		}
	}
}
