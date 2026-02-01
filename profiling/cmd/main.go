package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	tomatoSauce    = "tomato sauce"
	mozzarella     = "mozzarella"
	freshBasil     = "fresh basil"
	pepperoni      = "pepperoni"
	mushrooms      = "mushrooms"
	bellPeppers    = "bell peppers"
	onions         = "onions"
	ham            = "ham"
	pineapple      = "pineapple"
	sausage        = "sausage"
	bbqSauce       = "bbq sauce"
	grilledChicken = "grilled chicken"
	redOnions      = "red onions"
	bacon          = "bacon"
	cheddar        = "cheddar"
	parmesan       = "parmesan"
	gorgonzola     = "gorgonzola"
	oliveOil       = "olive oil"
	mixedMushrooms = "mixed mushrooms"
	truffleOil     = "truffle oil"
)

type Pizza struct {
	Name     string
	Price    float64
	Toppings []string
}

type OrderItem struct {
	Pizza    Pizza
	Quantity int
}

type Cart struct {
	Items []OrderItem
}

type HistoricalOrder struct {
	Items []string
}

var availablePizzas = []Pizza{
	{"Margherita", 9.99, []string{tomatoSauce, mozzarella, freshBasil}},
	{"Pepperoni", 11.49, []string{tomatoSauce, mozzarella, pepperoni}},
	{"Vegetarian", 10.99, []string{tomatoSauce, mozzarella, mushrooms, bellPeppers, onions}},
	{"Hawaiian", 11.99, []string{tomatoSauce, mozzarella, ham, pineapple}},
	{"Supreme", 13.49, []string{tomatoSauce, mozzarella, pepperoni, sausage, mushrooms, bellPeppers}},
	{"BBQ Chicken", 12.99, []string{bbqSauce, mozzarella, grilledChicken, redOnions}},
	{"Meat Lovers", 13.99, []string{tomatoSauce, mozzarella, pepperoni, sausage, bacon, ham}},
	{"Four Cheese", 11.49, []string{mozzarella, cheddar, parmesan, gorgonzola}},
	{"Mushroom Truffle", 12.49, []string{oliveOil, mozzarella, mixedMushrooms, truffleOil}},
}

var historicalOrders = generateHistoricalOrders()

func generateHistoricalOrders() []HistoricalOrder {
	const numOrders = 3000
	orders := make([]HistoricalOrder, numOrders)
	numPizzas := len(availablePizzas)
	pizzaNames := make([]string, numPizzas)
	for i, p := range availablePizzas {
		pizzaNames[i] = p.Name
	}
	for i := range numOrders {
		numItems := 2 + (i % 3)
		items := make([]string, numItems)
		for j := range numItems {
			idx := (i*5 + j*13) % numPizzas
			items[j] = pizzaNames[idx]
		}
		orders[i] = HistoricalOrder{Items: items}
	}
	return orders
}

func getPopularities() map[string]int {
	pop := make(map[string]int)
	for _, p := range availablePizzas {
		pop[p.Name] = 0
	}
	for _, order := range historicalOrders {
		for _, itemName := range order.Items {
			if _, exists := pop[itemName]; exists {
				pop[itemName]++
			}
		}
	}
	return pop
}

func displayMenu() {
	pop := getPopularities()

	type rankedPizza struct {
		Name  string
		Count int
	}
	ranks := make([]rankedPizza, 0, len(availablePizzas))
	for name, count := range pop {
		ranks = append(ranks, rankedPizza{Name: name, Count: count})
	}
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].Count > ranks[j].Count
	})

	fmt.Println("=== Popular Picks ===")
	for i := 0; i < 5 && i < len(ranks); i++ {
		fmt.Printf("- %s (ordered %d times)\n", ranks[i].Name, ranks[i].Count)
	}

	fmt.Println("\n=== Full Menu ===")
	for i, p := range availablePizzas {
		toppings := strings.Join(p.Toppings, ", ")
		fmt.Printf("%3d. %-25s $%6.2f  %s\n", i+1, p.Name, p.Price, toppings)
	}
}

func calculateTotal(cart Cart) float64 {
	total := 0.0
	for _, item := range cart.Items {
		total += item.Pizza.Price * float64(item.Quantity)
	}
	return total
}

func viewCart(cart Cart) {
	if len(cart.Items) == 0 {
		fmt.Println("Your cart is empty.")
		return
	}
	fmt.Println("=== Your Cart ===")
	for i, item := range cart.Items {
		subtotal := item.Pizza.Price * float64(item.Quantity)
		fmt.Printf("%d. %s x %d - $%.2f\n", i+1, item.Pizza.Name, item.Quantity, subtotal)
	}
	fmt.Printf("Total: $%.2f\n", calculateTotal(cart))
}

func checkout(cart *Cart) {
	if len(cart.Items) == 0 {
		fmt.Println("Cannot checkout empty cart.")
		return
	}
	total := calculateTotal(*cart)
	fmt.Println("Processing your order...")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Payment successful!")

	fmt.Println("\n=== Order Receipt ===")
	for _, item := range cart.Items {
		subtotal := item.Pizza.Price * float64(item.Quantity)
		fmt.Printf("%s x %d - $%.2f\n", item.Pizza.Name, item.Quantity, subtotal)
	}
	fmt.Printf("Grand Total: $%.2f\n", total)
	fmt.Println("Thank you for ordering with PizzaCLI!")

	cart.Items = cart.Items[:0]
}

func handleAddCommand(fields []string, cart *Cart) {
	if len(fields) != 3 {
		fmt.Println("Usage: add <number> <quantity>")
		return
	}

	idx, err1 := strconv.Atoi(fields[1])
	qty, err2 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil || idx < 1 || idx > len(availablePizzas) || qty <= 0 {
		fmt.Println("✗ Invalid pizza number or quantity")
		return
	}

	pizza := availablePizzas[idx-1]
	cart.Items = append(cart.Items, OrderItem{Pizza: pizza, Quantity: qty})
	fmt.Printf("✓ Added %d × %s\n", qty, pizza.Name)
}

func handleRemoveCommand(fields []string, cart *Cart) {
	if len(fields) != 2 {
		fmt.Println("Usage: remove <position>")
		return
	}

	pos, err := strconv.Atoi(fields[1])
	if err != nil || pos < 1 || pos > len(cart.Items) {
		fmt.Println("✗ Invalid position")
		return
	}

	removed := cart.Items[pos-1].Pizza.Name
	cart.Items = append(cart.Items[:pos-1], cart.Items[pos:]...)
	fmt.Printf("✓ Removed %s\n", removed)
}

func processCommand(fields []string, cart *Cart) bool {
	cmd := fields[0]

	switch cmd {
	case "exit", "quit":
		fmt.Println("Goodbye!")
		return true
	case "viewcart", "cart":
		viewCart(*cart)
	case "checkout":
		checkout(cart)
	case "add":
		handleAddCommand(fields, cart)
	case "remove", "removecart":
		handleRemoveCommand(fields, cart)
	default:
		fmt.Println("Unknown command. Try: add <num> <qty>, remove <pos>, viewcart, checkout, exit")
	}
	return false
}

func main() {
	// Start CPU profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile file:", err)
		return
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		fmt.Println("Error starting CPU profiling:", err)
		return
	}
	defer pprof.StopCPUProfile()

	var cart Cart
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to PizzaCLI!")

	for {
		displayMenu()

		fmt.Print("\nEnter command: ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(strings.ToLower(line))
		if processCommand(fields, &cart) {
			return
		}
	}
}
