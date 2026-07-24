package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/starfederation/datastar-go/datastar"
	"gopkg.in/ini.v1"
)

var db *sql.DB
var rows *sql.Rows
var err error

type CurrentPage struct {
	Page int `json:"page"`
}

type Customer struct {
	Id             int    `json:"id"`
	CustomerName   string `json:"customerName"`
	Subscriber     bool   `json:"subscriber"`
	Cohort         string `json:"cohort"`
	SignupDate     string `json:"signupDate"`
	LifetimeValue  int
	NumberOfOrders int
	NetProfit      int
}
type CustomerPatch struct {
	NewCustomerName string `json:"newCustomerName"`
	NewSubscriber   bool   `json:"newSubscriber"`
	NewCohort       string `json:"newCohort"`
	NewSignupDate   string `json:"newSignupDate"`
}
type Product struct {
	Id                   int    `json:"id"`
	ProductName          string `json:"productName"`
	Description          string `json:"description"`
	Stock                int    `json:"stock"`
	MinimumOrderQuantity int    `json:"minimumOrderQuanity"`
	RetailPrice          int    `json:"retailPrice"`
}
type ProductPatch struct {
	NewProductName         string `json:"newProductName"`
	NewDescription         string `json:"newDescription"`
	NewStockAtWarehouse    int    `json:"newStockAtWarehouse"`
	NewMinimumOrderQuanity int    `json:"newMinimumOrderQuanity"`
	NewPrice               int    `json:"newPrice"`
}
type InventoryCosting struct {
	TimePeriod    string `json:"timePeriod"`
	ProductId     int    `json:"productId"`
	InventoryType string `json:"inventoryType"`
	Quantity      int    `json:"quantity"`
	UnitCost      int    `json:"unitCost"`
	Total         int    `json:"total"`
	QuantityCount int    `json:"quantityCount"`
	TotalCount    int    `json:"totalCount"`
}
type Employee struct {
	Id           int    `json:"id"`
	EmployeeName string `json:"employeeName"`
	JobTitle     string `json:"jobTitle"`
	Department   string `json:"department"`
	StartDate    string `json:"startDate"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Salary       int    `json:"salary"`
}
type EmployeePatch struct {
	NewEmployeeName string `json:"newEmployeeName"`
	NewJobTitle     string `json:"newJobTitle"`
	NewDepartment   string `json:"newDepartment"`
	NewStartDate    string `json:"newStartDate"`
	NewPhone        string `json:"newPhone"`
	NewEmail        string `json:"newEmail"`
	NewSalary       int    `json:"newSalary"`
}
type Order struct {
	Id                int    `json:"id"`
	CustomerId        int    `json:"customerId"`
	ProductId         int    `json:"productId"`
	Quantity          int    `json:"quantity"`
	Total             int    `json:"total"`
	Progress          string `json:"progress"`
	OrderDate         string `json:"orderDate"`
	CostOfGoodsSold   int    `json:"costOfGoodsSold"`
	OrderProfitMargin int    `json:"orderProfitMargin"`
}
type OrderPatch struct {
	NewCustomerId int    `json:"newCustomerId"`
	NewProductId  int    `json:"newProductId"`
	NewQuantity   int    `json:"newQuantity"`
	NewProgress   string `json:"newProgress"`
	NewOrderDate  string `json:"newOrderDate"`
}
type Review struct {
	OrderId int    `json:"id"`
	Rating  int    `json:"rating"`
	Review  string `json:"review"`
}

//	type Transaction struct {
//		Id              int     `json:"id"`
//		Counterparty    string  `json:"counterparty"`
//		TransactionDate string  `json:"transactionDate"`
//		Amount          float32 `json:"amount"`
//		AccountType     string  `json:"accountType"`
//		Category        string  `json:"rating"`
//		Notes           string  `json:"notes"`
//	}
//
//	type TransactionPatch struct {
//		NewCounterparty      string  `json:"newCounterparty"`
//		NewTransactionDadete string  `json:"newTransactionDate"`
//		NewAmount            float32 `json:"newAmount"`
//		NewAccountType       string  `json:"newAccountType"`
//		NewCategory          string  `json:"newRating"`
//		NewNotes             string  `json:"newNotes"`
//	}
type CashFlowStatements struct {
	TimePeriod        string `json:"notes"`
	OperatingCashFlow int    `json:"operatingCashFlow"`
	TotalSales        int    `json:"totalSales"`
	CashSpentOnAssets int    `json:"cashSpentOnAssets"`
	OperatingExpenses int    `json:"operatingExpenses"`
}
type IncomeStatements struct {
	TimePeriod                              string `json:"timePeriod"`
	TotalSales                              int    `json:"totalSales"`
	CostOfGoodsSold                         int    `json:"costOfGoodsSold"`
	Profit                                  int    `json:"profit"`
	PromotionExpenses                       int    `json:"promotionExpenses"`
	SellingAndGeneralAdministrativeExpenses int    `json:"sellingAndGeneralAdministrativeExpenses"`
	DepreciationAndAmoritization            int    `json:"depreciationAndAmoritization"`
	NetIncome                               int
}
type BalanceSheets struct {
	TimePeriod           string `json:"timePeriod"`
	Cash                 int    `json:"cash"`
	AccountsReceivable   int    `json:"accountsReceivable"`
	PrepaidExpenses      int    `json:"prepaidExpenses"`
	Inventory            int    `json:"inventory"`
	PropertyAndEquipment int    `json:"propertyAndEquipment"`
	Goodwill             int    `json:"goodwill"`
	AccountsPayable      int    `json:"accountsPayable"`
	AccruedExpenses      int    `json:"accruedExpenses"`
	UnearnedRevenue      int    `json:"unearnedRevenue"`
	LongTermDebt         int    `json:"longTermDebt"`
	Equity               int    `json:"equity"`
}
type ChartOfAccounts struct {
	AccountCode        int
	AccountName        string
	AccountType        string
	FinancialStatement string
}

type GeneralLedger struct {
	AccountCode int
	AccountName string
	Ledger      struct {
		JournalEntryID [100]int
		CreditOrDebit  [100]string
		Balance        [100]int
	}
}
type FinancialRatios struct {
	TimePeriod                 string
	CurrentRatio               float32
	AcidTestRatio              float32
	InventoryTurnover          float32
	InventoryDays              float32
	AccountsReceivableTurnover float32
	CollectionPeriod           float32
	DebtRatio                  float32
	GrossProfitPercentage      float32
	ReturnOnSales              float32
	ReturnOnAssets             float32
	CashConversionCycle        float32
}
type JournalEntries struct {
	Id            int
	Date          string
	DebitAccount  int
	CreditAccount int
	Debit         int
	Credit        int
	Notes         string
}
type TrialBalance struct {
	AccountCode       int
	AccountName       string
	UnadjustedBalance int
	AdjustingEntry    int
	AdjustedBalance   int
}

// type Pagination struct {
// 	ItemsPerPage                  int `json:"itemsPerPage"`
// 	PageNumber                    int `json:"pageNumber"`
// 	ItemsFetchedFromDatabaseCount int
// }

// type Search struct {
// 	Query        string `json:"query"`
// 	SearchBy     string `json:"searchBy"`
// 	ItemsPerPage int    `json:"itemsPerPage"`
// 	PageNumber   int    `json:"pageNumber"`
// }

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		config, err := ini.Load("config.ini")
		if err != nil {
			log.Fatal("error: ", err)
		}
		http.Redirect(w, r, "/"+config.Section("frontend").Key("DefaultPage").String(), http.StatusFound)
	}
}

func cancel(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.Redirectf("/%v", r.FormValue("resource"))
}

func customers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT * FROM customers ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var customerRecords []Customer
		for rows.Next() {
			var customerFields Customer
			if err := rows.Scan(&customerFields.Id, &customerFields.CustomerName, &customerFields.Subscriber, &customerFields.Cohort, &customerFields.SignupDate); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			customerRecords = append(customerRecords, customerFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusBadRequest)
			return
		}
		rows, err = db.Query("SELECT customers.id, customer_name, COALESCE(ROUND(SUM(total),2),0) AS lifetime_value, (SELECT COUNT(*) FROM orders WHERE customer_id = customers.id) AS number_of_orders, (SELECT SUM(total) - SUM(cost_of_goods_sold) FROM orders GROUP BY customer_id) AS net_profit FROM customers JOIN orders on customer_id = customers.id GROUP BY customers.id ORDER BY customers.id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var rfm []Customer
		for rows.Next() {
			var rfmRecord Customer
			if err := rows.Scan(&rfmRecord.Id, &rfmRecord.CustomerName, &rfmRecord.LifetimeValue, &rfmRecord.NumberOfOrders, &rfmRecord.NetProfit); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rfm = append(rfm, rfmRecord)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr = rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusBadRequest)
			return
		}
		rows, err = db.Query("SELECT customer_id, SUM(total) - SUM(cost_of_goods_sold) AS net_profit FROM orders GROUP BY customer_id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()

		jinja, err := gonja.FromFile("templates/customers.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]any{
			"customerRecords": customerRecords,
			"rfm":             rfm,
		})
		jinja.Execute(w, data)
	case "POST":
		result, err := db.Exec("INSERT INTO customers (customer_name, subscriber, cohort, signup_date) VALUES (?, ?, ?, CURRENT_DATE())", r.FormValue("customer-name"), r.FormValue("subscriber"), r.FormValue("cohort"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/", http.StatusFound)
	case "DELETE":
		result, err := db.Exec("DELETE FROM customers WHERE id=?", r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/customers")
	case "PATCH":
		customerPatchSignals := &CustomerPatch{}
		if err := datastar.ReadSignals(r, customerPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE customers SET customer_name=? subscriber=?, cohort=?, signupDate=? WHERE id=?", customerPatchSignals.NewCustomerName, customerPatchSignals.NewSubscriber, customerPatchSignals.NewCohort, customerPatchSignals.NewSignupDate, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse.Redirect("/customers")
	}
}

func editCustomer(w http.ResponseWriter, r *http.Request) {
	var customerBeingEdited Customer
	id := r.FormValue("id")
	err := db.QueryRow("SELECT * FROM customers WHERE id = ?", id).Scan(&customerBeingEdited.Id, &customerBeingEdited.CustomerName, &customerBeingEdited.Subscriber, &customerBeingEdited.Cohort, &customerBeingEdited.SignupDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="%v">
			<td>%v</td>
			<td><input data-bind:new-customer-name name="new-customer-name" type="text" value="%v" required></td>
			<td>
				<select data-bind:new-subscriber name="new-subscriber" value="%v">
					<option value=0>false</option>
					<option value=1>true</option>
				</select>
			</td>
		    <td>
				<select data-bind:new-cohort name="new-cohort">
					<option value="%v">%v (original)</option>
					<option value="1">lapsed prospect</option>
					<option value="2">one purchase</option>
					<option value="3">two purchases</option>
					<option value="4">VIP</option>
					<option value="5">churned</option>
				</select>
			</td>
		    <td><input data-bind:new-signup-date name="new-signup-date" type="date" value="%v"></td>
			<td><button data-on:click="@get('/customers')">Cancel</button></td>
			<td><button data-on:click="@patch('/customers?id=%v')">Update</button></td>
		</tr>
		`, id, customerBeingEdited.Id, customerBeingEdited.CustomerName, customerBeingEdited.Subscriber, customerBeingEdited.Cohort, customerBeingEdited.Cohort, customerBeingEdited.SignupDate, id))
}

/*
	func searchCustomers(w http.ResponseWriter, r *http.Request) {
		signals := &Search{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pattern := "%" + signals.Query + "%"
		switch signals.SearchBy {
		case "customer-name":
			rows, err = db.Query("SELECT * FROM customers WHERE customer_name LIKE ?", pattern)
		case "cohort":
			rows, err = db.Query("SELECT * FROM customers WHERE cohort LIKE ?", pattern)
		case "subscriber":
			rows, err = db.Query("SELECT * FROM customers WHERE subscriber LIKE 1", pattern)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var results []Customer
		for rows.Next() {
			var customerFields Customer
			if err := rows.Scan(&customerFields.Id, &customerFields.CustomerName); err != nil {
				msg := fmt.Errorf("search error: %v", err)
				fmt.Println(msg)
				return
			}
			results = append(results, customerFields)
		}
		ipp := max(1, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := results[ipp*(min(pn-1, len(results)/ipp)) : min(pn*ipp, len(results))]
		t, err := template.New("results").Parse(`
			<tbody id="current-table">
				{{range .pages}}
				<tr id="row-{{.ID}}">
				<td data-bind:id>{{.Id}}</td>
			    <td data-bind:customer-name>{{.CustomerName}}</td>
			    <td data-bind:subsrciber>{{.Subsrciber}}</td>
			    <td data-bind:cohort>{{.Cohort}}</td>
			    <td data-bind:signup-date>{{.SignupDate}}</td>
				<td><button data-on:click="confirm('Are you sure?') && @delete('/customers?id={{.ID}}')">Delete</button></td>
				<td><button data-on:click="@get('/customers/edit?id={{.ID}}&name={{.Name}}')">Edit</button></td>
				</tr>
				{{end}}
			</tbody>
			`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var builder strings.Builder
		resultsPage := map[string][]Customer{"pages": page}
		t.Execute(&builder, resultsPage)
		searchResult := builder.String()
		sse := datastar.NewSSE(w, r)
		sse.PatchElements(searchResult)
	}
*/
func products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT * FROM products ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var productRecords []Product
		for rows.Next() {
			var productFields Product
			if err := rows.Scan(&productFields.Id, &productFields.ProductName, &productFields.Description, &productFields.RetailPrice, &productFields.MinimumOrderQuantity); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			productRecords = append(productRecords, productFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, err = db.Query("SELECT * FROM inventory_costing")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var inventoryCosting []InventoryCosting
		for rows.Next() {
			var inventoryCostingRecord InventoryCosting
			if err := rows.Scan(&inventoryCostingRecord.ProductId, &inventoryCostingRecord.TimePeriod, &inventoryCostingRecord.InventoryType, &inventoryCostingRecord.Quantity, &inventoryCostingRecord.UnitCost, &inventoryCostingRecord.Total); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			inventoryCosting = append(inventoryCosting, inventoryCostingRecord)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var count int
		for i, _ := range inventoryCosting {
			if inventoryCosting[i].InventoryType == "Purchase" {
				count += inventoryCosting[i].Quantity
			} else if inventoryCosting[i].InventoryType == "Cost of Goods Sold" {
				count -= inventoryCosting[i].Quantity
			}
			q := count
			inventoryCosting[i].QuantityCount = count
			inventoryCosting[i].TotalCount = q * inventoryCosting[i].UnitCost
		}
		// the product stock is equal to the last quantity in IC
		for productRecordIndex, _ := range productRecords {
			for i := len(inventoryCosting) - 1; i >= 0; i-- {
				if inventoryCosting[i].ProductId == productRecords[productRecordIndex].Id {
					productRecords[productRecordIndex].Stock = inventoryCosting[i].Quantity
				}
			}
		}
		template, err := gonja.FromFile("templates/products.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]any{
			"productRecords":   productRecords,
			"inventoryCosting": inventoryCosting,
		})
		template.Execute(w, data)
	case "POST":
		if r.FormValue("form") == "product" {
			result, err := db.Exec("INSERT INTO products (product_name, description, price, minimum_order_quantity) VALUES (?,?,?*100,?)", r.FormValue("product-name"), r.FormValue("description"), r.FormValue("price"), r.FormValue("minimum-order-quantity"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := result.LastInsertId()
			if err != nil {
				log.Printf("error: %v, %v", id, err.Error())
			}
			http.Redirect(w, r, "/products", http.StatusFound)
		} else {
			result, err := db.Exec("INSERT INTO inventory_costing (product_id, time_period, inventory_type, quantity, unit_cost, total) VALUES (?,?,?,?,?*100,?*?*100)",
				r.FormValue("product"),
				r.FormValue("time-period"),
				r.FormValue("inventory-type"),
				r.FormValue("quantity"),
				r.FormValue("unit-cost"),
				r.FormValue("quantity"), r.FormValue("unit-cost"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := result.LastInsertId()
			if err != nil {
				log.Printf("error: %v, %v", id, err.Error())
			}
			// add journal entry with inventory account
			result, err = db.Exec("INSERT INTO journal_entries (entry_date, debit_account, credit_account, debit, credit, notes) VALUES (CURRENT_DATE(),?,?,?*?*100,?*?*100,?)",
				1104,
				2401,
				r.FormValue("quantity"), r.FormValue("unit-cost"),
				r.FormValue("quantity"), r.FormValue("unit-cost"),
				"increased inventory")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err = result.LastInsertId()
			if err != nil {
				log.Printf("error: %v, %v", id, err.Error())
			}
			http.Redirect(w, r, "/products", http.StatusFound)
		}
	case "DELETE":
		result, err := db.Exec("DELETE FROM products WHERE id=?", r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/products")
	case "PATCH":
		productPatchSignals := &ProductPatch{}
		if err := datastar.ReadSignals(r, productPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE products SET product_name=?, description=?, price=?*100, stock=?, minumum_order_quantity=? WHERE id=?", productPatchSignals.NewProductName, productPatchSignals.NewDescription, productPatchSignals.NewPrice, productPatchSignals.NewStockAtWarehouse, productPatchSignals.NewMinimumOrderQuanity, r.FormValue("id"))
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse.Redirect("/products")
	}
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	var productBeingEdited Product
	id := r.FormValue("id")
	err := db.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan(&productBeingEdited.Id, &productBeingEdited.ProductName, &productBeingEdited.Description, &productBeingEdited.RetailPrice, &productBeingEdited.Stock, &productBeingEdited.MinimumOrderQuantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="%v">
  			<td>%v</td>
	        <td><input data-bind:new-product-name name="new-product-name" value="%v" required></td>
			<td><input data-bind:new-description name="new-description" value="%v" required></td>
	        <td><input type="number" data-bind:stock name="new-stock" value="%v" required></td>
	        <td><input type="number" data-bind:new-minimum-order-quantity name="new-minimum-order-quantity" value="%v" required></td>
	        <td><input type="number" data-bind:new-price name="new-price" value="%v" required></td>
         	<td><button data-on:click="@get('/cancel?resource=products')">Cancel</button></td>
			<td><button data-on:click="@patch('/products?id=%v')">Update</button></td>
        </tr>
		`, id, productBeingEdited.Id, productBeingEdited.ProductName, productBeingEdited.Description, productBeingEdited.Stock, productBeingEdited.MinimumOrderQuantity, productBeingEdited.RetailPrice, id))
}

/*
	func searchProducts(w http.ResponseWriter, r *http.Request) {
		signals := &Search{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pattern := "%" + signals.Query + "%"
		switch signals.SearchBy {
		case "product-name":
			rows, err = db.Query("SELECT * FROM products WHERE product_name LIKE ?", pattern)
		case "description":
			rows, err = db.Query("SELECT * FROM products WHERE description LIKE ?", pattern)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var results []Product
		for rows.Next() {
			var productFields Product
			if err := rows.Scan(&productFields.Id, &productFields.ProductName, &productFields.Description, &productFields.StockAtWarehouse, &productFields.StockAtManufacturer, &productFields.Price, &productFields.ReorderPoint, &productFields.LeadTimeDays); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			results = append(results, productFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("rows: ", results)
		ipp := max(1, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := results[ipp*(min(pn-1, len(results)/ipp)) : min(pn*ipp, len(results))]
		t, err := template.New("results").Parse(`
			<tbody id="current-table">
				{{range .pages}}
				<tr id="row-{{.ID}}">
					<td data-signals="{ID: {{.ID}} }">{{.ID}}</td>
					<td data-bind:product-name>{{.Name}}</td>
					<td><button data-on:click="confirm('Are you sure?') && @delete('/customers?id={{.ID}}')">Delete</button></td>
					<td><button data-on:click="@get('/customers/edit?id={{.ID}}&name={{.Name}}')">Edit</button></td>
				</tr>
				{{end}}
			</tbody>`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var builder strings.Builder
		resultsPage := map[string][]Product{"pages": page}
		t.Execute(&builder, resultsPage)
		searchResult := builder.String()
		sse := datastar.NewSSE(w, r)
		sse.PatchElements(searchResult)
	}
*/
func employees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT * FROM employees")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var employeeRecords []Employee
		for rows.Next() {
			var employeeFields Employee
			if err := rows.Scan(&employeeFields.Id, &employeeFields.EmployeeName, &employeeFields.JobTitle, &employeeFields.Department, &employeeFields.StartDate, &employeeFields.Phone, &employeeFields.Email, &employeeFields.Salary); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			employeeRecords = append(employeeRecords, employeeFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusBadRequest)
			return
		}
		template, err := gonja.FromFile("templates/employees.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]any{
			"employeeRecords": employeeRecords,
		})
		template.Execute(w, data)
	case "POST":
		result, err := db.Exec("INSERT INTO employees (employee_name, job_title, department, start_data, phone, company_email, salary) VALUES (?,?,?,?,?,?*100)", r.FormValue("employee-name"), r.FormValue("job-title"), r.FormValue("department"), r.FormValue("start-date"), r.FormValue("phone"), r.FormValue("company-email"), r.FormValue("salary"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/employees", http.StatusFound)
	case "DELETE":
		result, err := db.Exec("DELETE FROM employees WHERE id=?", r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/employees")
	case "PATCH":
		employeeSignalsPatch := &EmployeePatch{}
		if err := datastar.ReadSignals(r, employeeSignalsPatch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE employees SET employee_name=?, job_title=?, department=?, start_date=?, phone=?, company_email=?, salary=?*100 WHERE id=?", employeeSignalsPatch.NewEmployeeName, employeeSignalsPatch.NewJobTitle, employeeSignalsPatch.NewDepartment, employeeSignalsPatch.NewStartDate, employeeSignalsPatch.NewPhone, employeeSignalsPatch.NewEmail, employeeSignalsPatch.NewSalary, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse.Redirect("/employees")
	}
}

func editEmployee(w http.ResponseWriter, r *http.Request) {
	// TODO: when table is replaced with form, currency should be formatted
	var employeeBeingEdited Employee
	id := r.FormValue("id")
	err := db.QueryRow("SELECT * FROM employees WHERE id = ?", id).Scan(&employeeBeingEdited.Id, &employeeBeingEdited.EmployeeName, &employeeBeingEdited.JobTitle, &employeeBeingEdited.Department, &employeeBeingEdited.StartDate, &employeeBeingEdited.Phone, &employeeBeingEdited.Email, &employeeBeingEdited.Salary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="%v">
			<td>%v</td>
			<td><input name="new-employee-name" type="text" value="%v" data-bind:new-employee-name required>
			<td><input name="new-job-title" type="text" value="%v" data-bind:new-job-title required></td>
			<td><input name="new-department" type="text" value="%v" data-bind:new-department required></td>
			<td><input name="new-start-date" type="date" value="%v" data-bind:new-start-date required></td>
			<td><input name="new-phone" type="phone" value="%v" data-bind:new-phone required></td>
			<td><input name="new-email" type="email" value="%v" data-bind:new-email required></td>
			<td class="currency"><input name="salary" type="number" value="%v" data-bind:new-salary required></td>
			<td><button data-on:click="@get('/employees')">Cancel</button></td>
			<td><button data-on:click="@patch('/employees?id=%v')">Update</button></td>
		</tr>
		`, id, employeeBeingEdited.Id, employeeBeingEdited.EmployeeName, employeeBeingEdited.JobTitle, employeeBeingEdited.Department, employeeBeingEdited.StartDate, employeeBeingEdited.Phone, employeeBeingEdited.Email, employeeBeingEdited.Salary/100, id))
}

/*
	func searchEmployees(w http.ResponseWriter, r *http.Request) {
		signals := &Search{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pattern := "%" + signals.Query + "%"
		switch signals.SearchBy {
		case "employee-name":
			rows, err = db.Query("SELECT * FROM employees WHERE employee_name LIKE ?", pattern)
		case "job-title":
			rows, err = db.Query("SELECT * FROM employees WHERE job_title LIKE ?", pattern)
		case "department":
			rows, err = db.Query("SELECT * FROM employees WHERE department LIKE ?", pattern)
		case "start-date":
			rows, err = db.Query("SELECT * FROM employees WHERE start_date LIKE ?", pattern)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer rows.Close()
			var results []Employee
			for rows.Next() {
				var employeeFields Employee
				if err := rows.Scan(&employeeFields.Id, &employeeFields.EmployeeName, &employeeFields.JobTitle, &employeeFields.Department, &employeeFields.StartDate, &employeeFields.Phone, &employeeFields.Email, &employeeFields.Salary); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				results = append(results, employeeFields)
			}
			if err = rows.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = rows.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ipp := max(1, signals.ItemsPerPage)
			pn := max(1, signals.PageNumber)
			page := results[ipp*(min(pn-1, len(results)/ipp)) : min(pn*ipp, len(results))]
			t, err := template.New("results").Parse(`
				<tbody id="current-table">
					{{range .pages}}
					<tr id="row-{{.Id}}">
						<td data-bind:id>{{.Id}}</td>
						<td data-bind:employee-name>{{.EmployeeName}}</td>
						<td data-bind:job-title>{{.JobTitle}}</td>
						<td data-bind:department>{{.Department}}</td>
						<td data-bind:phone>{{.Phone}}</td>
	    				<td data-bind:company-email>{{.CompanyEmail}}</td>
						<td><button data-on:click="confirm('Are you sure?') && @delete('/customers?id={{.Id}}')">Delete</button></td>
						<td><button data-on:click="@get('/customers/edit?id={{.Id}}&name={{.Name}}')">Edit</button></td>
					</tr>
					{{end}}
				</tbody>
				`)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var builder strings.Builder
			resultsPage := map[string][]Employee{"pages": page}
			t.Execute(&builder, resultsPage)
			searchResult := builder.String()
			sse := datastar.NewSSE(w, r)
			sse.PatchElements(searchResult)
		}
	}
*/
func orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT * FROM orders")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var orderRecords []Order
		for rows.Next() {
			var orderFields Order
			if err := rows.Scan(&orderFields.Id, &orderFields.CustomerId, &orderFields.ProductId, &orderFields.Quantity, &orderFields.Total, &orderFields.Progress, &orderFields.OrderDate, &orderFields.CostOfGoodsSold); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			orderRecords = append(orderRecords, orderFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusInternalServerError)
			return
		}
		// calculate profit margin
		for i, v := range orderRecords {
			orderRecords[i].OrderProfitMargin = v.Total - v.CostOfGoodsSold
		}
		// product name and id
		rows, err = db.Query("SELECT id, product_name FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var productRecords []Product
		for rows.Next() {
			var productFields Product
			if err := rows.Scan(&productFields.Id, &productFields.ProductName); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			productRecords = append(productRecords, productFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowErr = rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusInternalServerError)
			return
		}
		// customer name and id
		rows, err = db.Query("SELECT id, customer_name FROM customers")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var customerRecords []Customer
		for rows.Next() {
			var customerFields Customer
			if err := rows.Scan(&customerFields.Id, &customerFields.CustomerName); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			customerRecords = append(customerRecords, customerFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowErr = rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusInternalServerError)
			return
		}
		template, err := gonja.FromFile("templates/orders.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]any{
			"orderRecords":    orderRecords,
			"productRecords":  productRecords,
			"customerRecords": customerRecords,
		})
		template.Execute(w, data)
	case "POST":
		var productPrice int
		pid := r.FormValue("product-id")
		db.QueryRow("SELECT price FROM products WHERE id = ?", pid).Scan(&productPrice)
		// calculate weighted average cost of goods sold
		rows, err = db.Query("SELECT unit_cost, quantity, SUM(quantity) AS weights FROM inventory_costing WHERE product_id=? AND inventory_type=1 GROUP BY unit_cost, quantity", pid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		type Weight struct {
			UnitCost             int
			Quantity             int
			SumQuantity          int
			CostTimesQuantity    int
			SumCostTimesQuantity int
		}
		var weights []Weight
		for rows.Next() {
			var weight Weight
			if err := rows.Scan(&weight.UnitCost, &weight.Quantity, &weight.SumQuantity); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			weights = append(weights, weight)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusInternalServerError)
			return
		}
		for i, v := range weights {
			weights[i].CostTimesQuantity = v.UnitCost * int((float32(v.Quantity) / float32(v.SumQuantity)))
		}
		var sumCostTimesQuantity int
		for _, v := range weights {
			sumCostTimesQuantity += v.CostTimesQuantity
		}
		var weightedAverage int
		weightedAverage = int(float32(sumCostTimesQuantity) / float32(len(weights)))
		// insert order
		result, err := db.Exec("INSERT INTO orders (customer_id, product_id, quantity, total, progress, order_date, cost_of_goods_sold) VALUES (?,?,?,?*?, ?, CURRENT_DATE(), ?*?)",
			r.FormValue("customer-id"),
			r.FormValue("product-id"),
			r.FormValue("quantity"),
			r.FormValue("quantity"), productPrice,
			r.FormValue("progress"),
			r.FormValue("quantity"), weightedAverage,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		// record orders in journal entry; credit inventory/debit COGS, and then credit sales and debit cash/AR
		result, err = db.Exec("INSERT INTO journal_entries (entry_date, debit_account, credit_account, debit, credit, notes) VALUES (CURRENT_DATE(),?,?,?*?,?*?,?),(CURRENT_DATE(),?,?,?*?,?*?,?)",
			1101,
			2301,
			r.FormValue("quantity"), productPrice,
			r.FormValue("quantity"), productPrice,
			fmt.Sprintf("order from customer %v for product %v", r.FormValue("customer-id"), r.FormValue("product-id")),
			// inventory & cogs
			1104,
			2401,
			r.FormValue("quantity"), weightedAverage,
			r.FormValue("quantity"), weightedAverage,
			"decrease inventory from sale",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err = result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		// update inventory costing
		result, err = db.Exec("INSERT INTO inventory_costing (product_id, time_period, inventory_type, quantity, unit_cost, total) VALUES (?, CURRENT_DATE(),?,?,?, ?*?)",
			r.FormValue("product-id"),
			2,
			r.FormValue("quantity"),
			weightedAverage,
			r.FormValue("quantity"), weightedAverage,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err = result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/orders", http.StatusFound)
	case "DELETE":
		var productPrice int
		deletedOrder := Order{}
		if err := datastar.ReadSignals(r, deletedOrder); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pid := deletedOrder.ProductId
		cid := deletedOrder.CustomerId
		db.QueryRow("SELECT price FROM products WHERE id = ?", pid).Scan(&productPrice)
		result, err := db.Exec("DELETE FROM orders WHERE id=?", r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		// debit sales and credit cash
		// then debit inventory and credit cogs
		result, err = db.Exec("INSERT INTO journal_entries (entry_date, debit_account, credit_account, debit, credit, notes) VALUES (CURRENT_DATE(),?,?,?*?,?*?,?), (CURRENT_DATE(),?,?,?*?,?*?,?)",
			2301,
			1101,
			r.FormValue("quantity"), productPrice,
			r.FormValue("quantity"), productPrice,
			fmt.Sprintf("deleted order from customer %v for product %v", cid, pid))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err = result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		// add the inventory back. the unit cost is whatever the unit cost was when the order was made
		result, err = db.Exec("INSERT INTO inventory_costing (time_period, inventory_type, quantity, unit_cost, total) VALUES (CURRENT_DATE(),?,?,(SELECT AVG(unit_cost * quantity/SUM(quantity)) AS weighted_avg_cogs FROM inventory_costing), ?*(SELECT AVG(unit_cost * quantity/SUM(quantity)) AS weighted_avg_cogs FROM inventory_costing))",
			1,
			r.FormValue("quantity"),
			r.FormValue("quantity"),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err = result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/orders")
	case "PATCH":
		// todo: might reove the ability to edit orders since that would affect many other parts of the program
		orderPatchSignals := &OrderPatch{}
		if err := datastar.ReadSignals(r, orderPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var productPrice int
		db.QueryRow("SELECT price FROM products WHERE id = ?", orderPatchSignals.NewProductId).Scan(&productPrice)
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE orders SET customer_id=?, product_id=?, quantity=?, progress=?, total=?*?*100 WHERE id=?", orderPatchSignals.NewCustomerId, orderPatchSignals.NewProductId, orderPatchSignals.NewQuantity, orderPatchSignals.NewProgress, orderPatchSignals.NewQuantity, productPrice, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse.Redirect("/orders")
	}
}

/*
func searchOrders(w http.ResponseWriter, r *http.Request) {
	signals := &Search{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pattern := "%" + signals.Query + "%"
	switch signals.SearchBy {
	case "customer-id":
		rows, err = db.Query("SELECT * FROM orders WHERE customer_id LIKE ?", pattern)
	case "order-id":
		rows, err = db.Query("SELECT * FROM orders WHERE order_id LIKE ?", pattern)
	case "product-id":
		rows, err = db.Query("SELECT * FROM orders WHERE product_id LIKE ?", pattern)
	case "progress":
		rows, err = db.Query("SELECT * FROM orders WHERE progress LIKE ?", pattern)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	var results []Order
	for rows.Next() {
		var orderFields Order
		if err := rows.Scan(&orderFields.Id, &orderFields.CustomerId, &orderFields.ProductId, &orderFields.Quantity, &orderFields.Total, &orderFields.Progress); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		results = append(results, orderFields)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rowErr := rows.Close()
	if rowErr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ipp := max(1, signals.ItemsPerPage)
	pn := max(1, signals.PageNumber)
	page := results[ipp*(min(pn-1, len(results)/ipp)) : min(pn*ipp, len(results))]
	t, err := template.New("results").Parse(`
			<tbody id="current-table">
				{{range .pages}}
				<tr id="row-{{.Id}}">
					<td data-bind:order-id>{{.OrderId}}</td>
					<td data-bind:customer-id>{{.CustomerId}}</td>
					<td data-bind:product-id>{{.ProductId}}</td>
					<td data-bind:quantity>{{.Quantity}}</td>
					<td data-bind:total>{{.Total}}</td>
					<td data-bind:progress>{{.Progress}}</td>
					<td><button data-on:click="confirm('Are you sure?') && @delete('/orders?id={{.Id}}')">Delete</button></td>
					<td><button data-on:click="@get('/orders/edit?id={{.Id}}')">Edit</button></td>
				</tr>
				{{end}}
			</tbody>
			`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var builder strings.Builder
	resultsPage := map[string][]Order{"pages": page}
	t.Execute(&builder, resultsPage)
	searchResult := builder.String()
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(searchResult)
}
*/

func editOrder(w http.ResponseWriter, r *http.Request) {
	var orderBeingEdited Order
	id := r.FormValue("id")
	err := db.QueryRow("SELECT * FROM orders WHERE id = ?", id).Scan(&orderBeingEdited.Id, &orderBeingEdited.CustomerId, &orderBeingEdited.ProductId, &orderBeingEdited.Quantity, &orderBeingEdited.Total, &orderBeingEdited.Progress, &orderBeingEdited.OrderDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="row-%v">
			<td>%v</td>
			<td><input name="new-customer-id" type="number" value="%v" data-bind:new-customer-id required>
			<td><input name="new-product-id" type="number" value="%v" data-bind:new-product-id required></td>
			<td><input name="new-quantity" type="number" value="%v" data-bind:new-quantity required></td>
			<td>%v</td>
			<td>
				<select name="new-progress" data-bind:new-progress>
					<option value="%v">keep original</option>
					<option value="1">received</option>
					<option value="2">in progress</option>
					<option value="3">delivered</option>
				</select>
			</td>
			<td><input name="new-order-date" type="date" value="%v" data-bind:new-order-date required></td>
			<td><button data-on:click="@get('/orders')">Cancel</button></td>
			<td><button data-on:click="@patch('/orders?id=%v')">Update</button></td>
		</tr>
		`, id, orderBeingEdited.Id, orderBeingEdited.CustomerId, orderBeingEdited.ProductId, orderBeingEdited.Quantity, orderBeingEdited.Total, orderBeingEdited.Progress, orderBeingEdited.OrderDate, id))
}

func reviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT * FROM reviews")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var reviewRecords []Review
		for rows.Next() {
			var reviewFields Review
			if err := rows.Scan(&reviewFields.OrderId, &reviewFields.Rating, &reviewFields.Review); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reviewRecords = append(reviewRecords, reviewFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusBadRequest)
			return
		}
		rows, err = db.Query("SELECT orders.id, customers.customer_name FROM orders JOIN customers ON orders.customer_id = customers.id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		type OrderInformation struct {
			OrderId      int
			CustomerName string
		}
		var orderInformationRecords []OrderInformation
		for rows.Next() {
			var orderInformationFields OrderInformation
			if err := rows.Scan(&orderInformationFields.OrderId, &orderInformationFields.CustomerName); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			orderInformationRecords = append(orderInformationRecords, orderInformationFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr = rows.Close()
		if rowErr != nil {
			http.Error(w, rowErr.Error(), http.StatusBadRequest)
			return
		}
		template, err := gonja.FromFile("templates/reviews.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]any{
			"reviewRecords":           reviewRecords,
			"orderInformationRecords": orderInformationRecords,
		})
		template.Execute(w, data)
	case "POST":
		result, err := db.Exec("INSERT INTO reviews (order_id, rating, review) VALUES (?,?,?)", r.FormValue("order-id"), r.FormValue("rating"), r.FormValue("review"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/reviews", http.StatusFound)
	case "DELETE":
		result, err := db.Exec("DELETE FROM reviews WHERE order_id=?", r.FormValue("order-id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/reviews")
	}
}

/*
	func searchReviews(w http.ResponseWriter, r *http.Request) {
		signals := &Search{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pattern := "%" + signals.Query + "%"
		switch signals.SearchBy {
		case "customer-id":
			rows, err = db.Query("SELECT * FROM reviews WHERE order_id LIKE ?", pattern)
		case "rating":
			rows, err = db.Query("SELECT * FROM reviews WHERE order_id LIKE ?", pattern)
		case "review":
			rows, err = db.Query("SELECT * FROM reviews WHERE order_id LIKE ?", pattern)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var results []Review
		for rows.Next() {
			var reviewFields Review
			if err := rows.Scan(&reviewFields.OrderId, &reviewFields.Rating, &reviewFields.Review); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			results = append(results, reviewFields)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ipp := max(1, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := results[ipp*(min(pn-1, len(results)/ipp)) : min(pn*ipp, len(results))]
		t, err := template.New("results").Parse(`
				<tbody id="current-table">
					{{range .pages}}
					<tr id="row-{{.Id}}">
					    <td data-bind:order-id>{{.OrderId}}</td>
					    <td data-bind:rating>{{.Rating}}</td>
					    <td data-bind:review>{{.Review}}</td>
						<td><button data-on:click="confirm('Are you sure?') && @delete('/reviews?id={{.Id}}')">Delete</button></td>
						<td><button data-on:click="@get('/reviews/edit?id={{.Id}}&name={{.Name}}')">Edit</button></td>
					</tr>
					{{end}}
				</tbody>
				`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var builder strings.Builder
		resultsPage := map[string][]Review{"pages": page}
		t.Execute(&builder, resultsPage)
		searchResult := builder.String()
		sse := datastar.NewSSE(w, r)
		sse.PatchElements(searchResult)
	}

	func editReview(w http.ResponseWriter, r *http.Request) {
		signals := &Review{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		sse.PatchElements(fmt.Sprintf(`
			<tr id="row-%v">
				<td>%v</td>
				<td><input name="new-order-id" type="number" value="%v" data-bind:order-id required>
				<td><input name="new-rating" type="number" value="%v" min="1" max="5" data-bind:rating required></td>
				<td><input name="new-review" type="text" value="%v" data-bind:review required></td>
				<td><button data-on:click="@get('/orders')">Cancel</button></td>
				<td><button data-on:click="@patch('/orders?id=%v')">Update</button></td>
			</tr>
			`, signals.NewOrderId, signals.NewRating, signals.NewReview, signals.OrderId))
	}
*/
func finances(w http.ResponseWriter, r *http.Request) {
	var chartOfAccountsRecords []ChartOfAccounts
	accountMap := make(map[string]int)
	accountCodeMap := make(map[int]string)
	rows, err = db.Query("SELECT * FROM chart_of_accounts ORDER BY account_code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var chartOfAccounts ChartOfAccounts
		if err = rows.Scan(&chartOfAccounts.AccountCode, &chartOfAccounts.AccountName, &chartOfAccounts.AccountType, &chartOfAccounts.FinancialStatement); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chartOfAccountsRecords = append(chartOfAccountsRecords, chartOfAccounts)
		accountMap[chartOfAccounts.AccountName] = chartOfAccounts.AccountCode
		accountCodeMap[chartOfAccounts.AccountCode] = chartOfAccounts.AccountName
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowErr := rows.Close()
	if rowErr != nil {
		http.Error(w, rowErr.Error(), http.StatusInternalServerError)
		return
	}

	var journalEntryRecords []JournalEntries
	rows, err = db.Query("SELECT * FROM journal_entries ORDER BY entry_date")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var journalEntry JournalEntries
		if err = rows.Scan(&journalEntry.Id, &journalEntry.Date, &journalEntry.DebitAccount, &journalEntry.CreditAccount, &journalEntry.Debit, &journalEntry.Credit, &journalEntry.Notes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		journalEntryRecords = append(journalEntryRecords, journalEntry)
	}
	err = rows.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowErr = rows.Close()
	if rowErr != nil {
		http.Error(w, rowErr.Error(), http.StatusBadRequest)
		return
	}
	var timePeriod string
	if len(journalEntryRecords) > 0 {
		timePeriod = fmt.Sprintf("%v to %v", journalEntryRecords[0].Date, journalEntryRecords[len(journalEntryRecords)-1].Date)
	}
	// there is a credit and a debit that need to go to their respective accounts on the ledger
	var generalLedger [20]GeneralLedger
	count := 0
	for _, v := range chartOfAccountsRecords {
		generalLedger[count].AccountName = v.AccountName
		generalLedger[count].AccountCode = v.AccountCode
		for i, _ := range journalEntryRecords {
			if generalLedger[count].AccountCode == journalEntryRecords[i].DebitAccount {
				generalLedger[count].Ledger.JournalEntryID[i] = journalEntryRecords[i].Id
				generalLedger[count].Ledger.CreditOrDebit[i] = "Debit"
				generalLedger[count].Ledger.Balance[i] = journalEntryRecords[i].Debit
			} else if generalLedger[count].AccountCode == journalEntryRecords[i].CreditAccount {
				generalLedger[count].Ledger.JournalEntryID[i] = journalEntryRecords[i].Id
				generalLedger[count].Ledger.CreditOrDebit[i] = "Credit"
				generalLedger[count].Ledger.Balance[i] = journalEntryRecords[i].Credit
			}
		}
		count += 1
	}
	ledgerIndex := r.FormValue("ledger-index")
	if ledgerIndex == "" {
		ledgerIndex = "Cash"
	}
	var displayedLedger GeneralLedger
	for _, v := range generalLedger {
		if v.AccountName == ledgerIndex {
			displayedLedger = v
		}
	}

	switch r.Method {
	case "GET":
		var accountBalance []TrialBalance
		for _, v := range chartOfAccountsRecords {
			var accountBalanceRecord TrialBalance
			if err = db.QueryRow("SELECT account_code, (SELECT COALESCE(SUM(debit),0) FROM journal_entries WHERE debit_account = ? ORDER BY debit_account) - (SELECT COALESCE(SUM(credit),0) FROM journal_entries WHERE credit_account = ? ORDER BY credit_account) AS balance FROM chart_of_accounts WHERE account_code = ? ORDER BY account_code", v.AccountCode, v.AccountCode, v.AccountCode).Scan(&accountBalanceRecord.AccountCode, &accountBalanceRecord.UnadjustedBalance); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			accountBalance = append(accountBalance, accountBalanceRecord)
		}

		var trialBalance []TrialBalance
		for i, v := range chartOfAccountsRecords {
			var singleTrialBalance TrialBalance
			singleTrialBalance.AccountCode = v.AccountCode
			singleTrialBalance.AccountName = v.AccountName
			singleTrialBalance.UnadjustedBalance = accountBalance[i].UnadjustedBalance
			trialBalance = append(trialBalance, singleTrialBalance)
		}

		var balanceSheet BalanceSheets
		balanceSheet.TimePeriod = timePeriod
		var incomeStatement IncomeStatements
		incomeStatement.TimePeriod = timePeriod
		for _, v := range accountBalance {
			switch v.AccountCode {
			case accountMap["Cash"]:
				balanceSheet.Cash = v.UnadjustedBalance
			case accountMap["Accounts Receivable"]:
				balanceSheet.AccountsReceivable = v.UnadjustedBalance
			case accountMap["Prepaid Expenses"]:
				balanceSheet.PrepaidExpenses = v.UnadjustedBalance
			case accountMap["Inventory"]:
				balanceSheet.Inventory = v.UnadjustedBalance
			case accountMap["Property and Equipment"]:
				balanceSheet.PropertyAndEquipment = v.UnadjustedBalance
			case accountMap["Goodwill"]:
				balanceSheet.Goodwill = v.UnadjustedBalance
			case accountMap["Accounts Payable"]:
				balanceSheet.AccountsPayable = v.UnadjustedBalance
			case accountMap["Accured Expenses"]:
				balanceSheet.AccruedExpenses = v.UnadjustedBalance
			case accountMap["Unearned Revenue"]:
				balanceSheet.UnearnedRevenue = v.UnadjustedBalance
			case accountMap["Long Term Debt"]:
				balanceSheet.LongTermDebt = v.UnadjustedBalance
			case accountMap["Sales"]:
				incomeStatement.TotalSales = v.UnadjustedBalance
			case accountMap["Cost of Goods Sold"]:
				incomeStatement.CostOfGoodsSold = v.UnadjustedBalance
			case accountMap["Profit"]:
				incomeStatement.Profit = v.UnadjustedBalance
			case accountMap["Promotion Expenses"]:
				incomeStatement.PromotionExpenses = v.UnadjustedBalance
			case accountMap["Selling & General Administrative Expenses"]:
				incomeStatement.SellingAndGeneralAdministrativeExpenses = v.UnadjustedBalance
			case accountMap["Depreciation & Amoritization"]:
				incomeStatement.DepreciationAndAmoritization = v.UnadjustedBalance
			}
		}
		incomeStatement.NetIncome = incomeStatement.TotalSales - (incomeStatement.CostOfGoodsSold + incomeStatement.PromotionExpenses + incomeStatement.SellingAndGeneralAdministrativeExpenses + incomeStatement.DepreciationAndAmoritization)

		//operating cash flow
		cashFlowStatement := &CashFlowStatements{}
		cashFlowStatement.TimePeriod = timePeriod

		// if the account is an asset, add it. query or for loop
		if err = db.QueryRow("SELECT COALESCE(SUM(debit), 0) FROM journal_entries  WHERE debit_account IN (SELECT account_code FROM chart_of_accounts WHERE account_type = 1)").Scan(&cashFlowStatement.CashSpentOnAssets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cashFlowStatement.OperatingCashFlow = (incomeStatement.NetIncome + incomeStatement.DepreciationAndAmoritization) - cashFlowStatement.CashSpentOnAssets

		// sum of all cash expenses
		cashFlowStatement.OperatingExpenses = incomeStatement.CostOfGoodsSold + incomeStatement.SellingAndGeneralAdministrativeExpenses + incomeStatement.PromotionExpenses - incomeStatement.DepreciationAndAmoritization

		financialRatios := &FinancialRatios{}
		financialRatios.TimePeriod = timePeriod
		// sum of balance sheet accounts where acc type == assets over sum of balance sheet accounts where acc type == liability
		financialRatios.CurrentRatio = float32((balanceSheet.Cash + balanceSheet.Inventory + balanceSheet.AccountsReceivable + balanceSheet.PrepaidExpenses) / max(1, (balanceSheet.AccountsPayable+balanceSheet.AccruedExpenses+balanceSheet.UnearnedRevenue+balanceSheet.LongTermDebt)))
		// cash and a/r over sum bal sheet where acc type == liability
		financialRatios.AcidTestRatio = float32((balanceSheet.Cash + balanceSheet.AccountsReceivable) / max(1, (balanceSheet.AccountsPayable+balanceSheet.AccruedExpenses+balanceSheet.UnearnedRevenue+balanceSheet.LongTermDebt)))

		financialRatios.GrossProfitPercentage = float32((incomeStatement.TotalSales-cashFlowStatement.OperatingExpenses)/max(1, incomeStatement.TotalSales)) * 100

		data := exec.NewContext(map[string]interface{}{
			"journalEntries":    journalEntryRecords,
			"trialBalance":      trialBalance,
			"balanceSheet":      balanceSheet,
			"chartOfAccounts":   chartOfAccountsRecords,
			"incomeStatement":   incomeStatement,
			"cashFlowStatement": cashFlowStatement,
			"generalLedger":     generalLedger,
			"financialRatios":   financialRatios,
			"accountMap":        accountMap,
			"displayedLedger":   displayedLedger,
		})
		switch r.FormValue("statement") {
		case "income":
			jinja, err := gonja.FromFile("templates/income-statement-template.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jinja.Execute(w, data)
		case "balance":
			jinja, err := gonja.FromFile("templates/balance-sheet-template.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jinja.Execute(w, data)
		case "cash":
			jinja, err := gonja.FromFile("templates/cash-flow-statement-template.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jinja.Execute(w, data)
		case "ledger":
			jinja, err := gonja.FromFile("templates/general-ledger-template.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jinja.Execute(w, data)
		default:
			jinja, err := gonja.FromFile("templates/finances.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
			err = jinja.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println(err.Error())
				return
			}
		}
	case "POST":
		result, err := db.Exec("INSERT INTO journal_entries (entry_date, debit_account, credit_account, debit, credit, notes) VALUES (?, ?, ?, ?*100, ?*100, ?)", r.FormValue("date"), r.FormValue("debit-account"), r.FormValue("credit-account"), r.FormValue("debit"), r.FormValue("credit"), r.FormValue("notes"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/finances", http.StatusFound)
	}
}

func changeLedger(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.Redirectf("/finances?display-ledger=%v", r.FormValue("ledger-index"))
}

func main() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("error: ", err)
	}
	cfg := mysql.NewConfig()
	cfg.User = config.Section("database").Key("User").String()
	cfg.Passwd = config.Section("database").Key("Passwd").String()
	cfg.Net = config.Section("database").Key("Net").String()
	cfg.Addr = config.Section("database").Key("Addr").String()
	cfg.DBName = config.Section("database").Key("DBName").String()

	port := "localhost:" + config.Section("frontend").Key("Port").String()
	// Get a database handle.
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println("error: ", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("error: ", pingErr)
	}
	log.Println("Connected to MySQL")

	http.HandleFunc("/", index)
	http.HandleFunc("/cancel", cancel)

	http.HandleFunc("/customers", customers)
	//	http.HandleFunc("/customers/search", searchCustomers)
	http.HandleFunc("/customers/edit", editCustomer)

	http.HandleFunc("/products", products)
	//	http.HandleFunc("/products/search", searchProducts)
	http.HandleFunc("/products/edit", editProduct)

	http.HandleFunc("/orders", orders)
	//	http.HandleFunc("/orders/search", searchOrders)
	http.HandleFunc("/orders/edit", editOrder)

	http.HandleFunc("/employees", employees)
	//	http.HandleFunc("/employees/search", searchEmployees)
	http.HandleFunc("/employees/edit", editEmployee)

	http.HandleFunc("/reviews", reviews)
	//	http.HandleFunc("/reviews/search", searchReviews)
	//  http.HandleFunc("/reviews/edit", editReview)

	http.HandleFunc("/finances", finances)
	http.HandleFunc("/finances/change-ledger", changeLedger)

	log.Printf("Visit http://%v for the program", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
