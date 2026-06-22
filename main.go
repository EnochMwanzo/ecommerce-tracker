package main

import (
	"database/sql"
	"fmt"
	"html/template"
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
	Id            int    `json:"id"`
	CustomerName  string `json:"customerName"`
	Subscriber    bool   `json:"subscriber"`
	Cohort        string `json:"cohort"`
	SignupDate    string `json:"signupDate"`
	LifetimeValue float32
}
type CustomerPatch struct {
	NewCustomerName string `json:"newCustomerName"`
	NewSubscriber   bool   `json:"newSubscriber"`
	NewCohort       string `json:"newCohort"`
	NewSignupDate   string `json:"newSignupDate"`
}
type Product struct {
	Id                  int     `json:"id"`
	ProductName         string  `json:"productName"`
	Description         string  `json:"description"`
	StockAtWarehouse    int     `json:"stockAtWarehouse"`
	StockAtManufacturer int     `json:"stockAtManufacturer"`
	Price               float32 `json:"price"`
	ReorderPoint        int     `json:"reorderPoint"`
	LeadTimeDays        int     `json:"leadTimeDays"`
}
type ProductPatch struct {
	NewProductName         string  `json:"newProductName"`
	NewDescription         string  `json:"newDescription"`
	NewStockAtWarehouse    int     `json:"newStockAtWarehouse"`
	NewStockAtManufacturer int     `json:"newStockAtManufacturer"`
	NewPrice               float32 `json:"newPrice"`
}
type InventoryCosting struct {
	TimePeriod    string  `json:"timePeriod"`
	InventoryType string  `json:"inventoryType"`
	Quantity      int     `json:"quantity"`
	UnitCost      float32 `json:"unitCost"`
	Total         float32 `json:"total"`
	QuantityCount int     `json:"quantityCount"`
	TotalCount    float32 `json:"totalCount"`
}
type Employee struct {
	Id           int     `json:"id"`
	EmployeeName string  `json:"employeeName"`
	JobTitle     string  `json:"jobTitle"`
	Department   string  `json:"department"`
	StartDate    string  `json:"startDate"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Salary       float32 `json:"salary"`
}
type EmployeePatch struct {
	NewEmployeeName string  `json:"newEmployeeName"`
	NewJobTitle     string  `json:"newJobTitle"`
	NewDepartment   string  `json:"newDepartment"`
	NewStartDate    string  `json:"newStartDate"`
	NewPhone        string  `json:"newPhone"`
	NewEmail        string  `json:"newEmail"`
	NewSalary       float32 `json:"newSalary"`
}
type Order struct {
	Id         int     `json:"id"`
	CustomerId int     `json:"customerId"`
	ProductId  int     `json:"productId"`
	Quantity   int     `json:"quantity"`
	Total      float32 `json:"total"`
	Progress   string  `json:"progress"`
	OrderDate  string  `json:"orderDate"`
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
type Transaction struct {
	Id              int     `json:"id"`
	Counterparty    string  `json:"counterparty"`
	TransactionDate string  `json:"transactionDate"`
	Amount          float32 `json:"amount"`
	AccountType     string  `json:"accountType"`
	Category        string  `json:"rating"`
	Notes           string  `json:"notes"`
}
type TransactionPatch struct {
	NewCounterparty      string  `json:"newCounterparty"`
	NewTransactionDadete string  `json:"newTransactionDate"`
	NewAmount            float32 `json:"newAmount"`
	NewAccountType       string  `json:"newAccountType"`
	NewCategory          string  `json:"newRating"`
	NewNotes             string  `json:"newNotes"`
}
type CashFlowStatements struct {
	TimePeriod        string  `json:"notes"`
	OperatingCashFlow float32 `json:"operatingCashFlow"`
	TotalSales        float32 `json:"totalSales"`
	CashSpentOnAssets float32 `json:"cashSpentOnAssets"`
	OperatingExpenses float32 `json:"operatingExpenses"`
}
type IncomeStatements struct {
	TimePeriod                              string  `json:"timePeriod"`
	TotalSales                              float32 `json:"totalSales"`
	CostOfGoodsSold                         float32 `json:"costOfGoodsSold"`
	Profit                                  float32 `json:"profit"`
	PromotionExpenses                       float32 `json:"promotionExpenses"`
	SellingAndGeneralAdministrativeExpenses float32 `json:"sellingAndGeneralAdministrativeExpenses"`
	DepreciationAndAmoritization            float32 `json:"depreciationAndAmoritization"`
}
type BalanceSheets struct {
	TimePeriod           string  `json:"timePeriod"`
	Cash                 float32 `json:"cash"`
	AccountsReceivable   float32 `json:"accountsReceivable"`
	PrepaidExpenses      float32 `json:"prepaidExpenses"`
	Inventory            float32 `json:"inventory"`
	PropertyAndEquipment float32 `json:"propertyAndEquipment"`
	Goodwill             float32 `json:"goodwill"`
	AccountsPayable      float32 `json:"accountsPayable"`
	AccruedExpenses      float32 `json:"accruedExpenses"`
	UnearnedRevenue      float32 `json:"unearnedRevenue"`
	LongTermDebt         float32 `json:"longTermDebt"`
}
type ChartOfAccounts struct {
	AccountCode        int
	AccountName        string
	AccountType        string
	FinancialStatement string
}

type GeneralLedger struct {
	AccountCode    []int
	AccountName    []string
	JournalEntryID []int
	CreditOrDebit  []string
	Balance        []float32
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
	Debit         float32
	Credit        float32
	Notes         string
}
type TrialBalance struct {
	AccountCode       int
	AccountName       string
	UnadjustedBalance float32
	AdjustingEntry    float32
	AdjustedBalance   float32
}
type Pagination struct {
	ItemsPerPage                  int `json:"itemsPerPage"`
	PageNumber                    int `json:"pageNumber"`
	ItemsFetchedFromDatabaseCount int
}

type Search struct {
	Query        string `json:"query"`
	SearchBy     string `json:"searchBy"`
	ItemsPerPage int    `json:"itemsPerPage"`
	PageNumber   int    `json:"pageNumber"`
}

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/customers", http.StatusFound)
	}
}

func customers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT customers.id, customer_name, subscriber, cohort, signup_date, COALESCE(ROUND(SUM(total),2),0) AS lifetime_value FROM customers JOIN orders on customer_id = customers.id GROUP BY customers.id ORDER BY customers.id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var customerRecords []Customer
		for rows.Next() {
			var customerFields Customer
			if err := rows.Scan(&customerFields.Id, &customerFields.CustomerName, &customerFields.Subscriber, &customerFields.Cohort, &customerFields.SignupDate, &customerFields.LifetimeValue); err != nil {
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
		var templateVariables = map[string]any{
			"data": customerRecords,
		}
		t := template.Must(template.ParseFiles("templates/base.html", "templates/customers.html", "templates/styles.css"))
		t.Execute(w, templateVariables)
	case "POST":
		result, err := db.Exec("INSERT INTO customers (customer_name, subscriber, cohort, signup_date) VALUES (?, ?, ?, ?)", r.FormValue("customer-name"), r.FormValue("subscriber"), r.FormValue("cohort"), r.FormValue("signup-date"))
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
	err = db.QueryRow("SELECT * FROM customers WHERE id = ?", id).Scan(&customerBeingEdited.LifetimeValue)
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
			<td>%v</td>
			<td><button data-on:click="@get('/customers')">Cancel</button></td>
			<td><button data-on:click="@patch('/customers?id=%v')">Update</button></td>
		</tr>
		`, id, customerBeingEdited.Id, customerBeingEdited.CustomerName, customerBeingEdited.Subscriber, customerBeingEdited.Cohort, customerBeingEdited.Cohort, customerBeingEdited.SignupDate, customerBeingEdited.LifetimeValue, id))
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
		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var productRecords []Product
		for rows.Next() {
			var productFields Product
			if err := rows.Scan(&productFields.Id, &productFields.ProductName, &productFields.Description, &productFields.Price, &productFields.StockAtWarehouse, &productFields.StockAtManufacturer, &productFields.ReorderPoint, &productFields.LeadTimeDays); err != nil {
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
			if err := rows.Scan(&inventoryCostingRecord.TimePeriod, &inventoryCostingRecord.InventoryType, &inventoryCostingRecord.Quantity, &inventoryCostingRecord.UnitCost, &inventoryCostingRecord.Total); err != nil {
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
			q := float32(count)
			inventoryCosting[i].QuantityCount = count
			inventoryCosting[i].TotalCount = q * inventoryCosting[i].UnitCost
		}
		jinja, err := gonja.FromFile("templates/products.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]any{
			"productRecords":   productRecords,
			"inventoryCosting": inventoryCosting,
		})
		jinja.Execute(w, data)
	case "POST":
		if r.FormValue("form") == "product" {
			result, err := db.Exec("INSERT INTO products (product_name, description, price, stock_at_warehouse, stock_at_manufacturer) VALUES (?,?,?,?,?)", r.FormValue("product-name"), r.FormValue("description"), r.FormValue("price"), r.FormValue("stock-at-warehouse"), r.FormValue("stock-at-manufacturer"))
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
			result, err := db.Exec("INSERT INTO inventory_costing (time_period, inventory_type, quantity, unit_cost, total) VALUES (?,?,?,?,?*?)",
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
		result, err := db.Exec("UPDATE products SET product_name=?, description=?, price=?, stock_at_warehouse=?, stock_at_manufacturer=? WHERE id=?", productPatchSignals.NewProductName, productPatchSignals.NewDescription, productPatchSignals.NewPrice, productPatchSignals.NewStockAtWarehouse, productPatchSignals.NewStockAtManufacturer, r.FormValue("id"))
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
	err := db.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan(&productBeingEdited.Id, &productBeingEdited.ProductName, &productBeingEdited.Description, &productBeingEdited.Price, &productBeingEdited.StockAtWarehouse, &productBeingEdited.StockAtManufacturer, &productBeingEdited.ReorderPoint, &productBeingEdited.LeadTimeDays)
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
	        <td><input type="number" data-bind:new-stock-at-warehouse name="new-stock-at-warehouse" value="%v" required></td>
	        <td><input type="number" data-bind:new-stock-at-manufacturer name="new-stock-at-manufacturer" value="%v" required></td>
	        <td><input type="number" data-bind:new-price name="new-price" value="%v" required></td>
	        <td>%v</td>
	        <td>%v</td>
         	<td><button data-on:click="@get('/products')">Cancel</button></td>
			<td><button data-on:click="@patch('/products?id=%v')">Update</button></td>
        </tr>
		`, id, productBeingEdited.Id, productBeingEdited.ProductName, productBeingEdited.Description, productBeingEdited.StockAtWarehouse, productBeingEdited.StockAtManufacturer, productBeingEdited.Price, productBeingEdited.ReorderPoint, productBeingEdited.LeadTimeDays, id))
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/employees.html", "templates/styles.css"))
		t.Execute(w, map[string]any{
			"data": employeeRecords,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO employees (employee_name, job_title, department, phone, company_email, salary) VALUES (?,?,?,?,?,?)", r.FormValue("employee-name"), r.FormValue("job-title"), r.FormValue("department"), r.FormValue("phone"), r.FormValue("company-email"), r.FormValue("salary"))
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
		result, err := db.Exec("UPDATE employees SET employee_name=?, job_title=?, department=?, start_date=?, phone=?, company_email=?, salary=? WHERE id=?", employeeSignalsPatch.NewEmployeeName, employeeSignalsPatch.NewJobTitle, employeeSignalsPatch.NewDepartment, employeeSignalsPatch.NewStartDate, employeeSignalsPatch.NewPhone, employeeSignalsPatch.NewEmail, employeeSignalsPatch.NewSalary, r.FormValue("id"))
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
	var employeeBeingEdited Employee
	id := r.FormValue("id")
	err := db.QueryRow("SELECT * FROM employees WHERE id = ?", id).Scan(&employeeBeingEdited.Id, &employeeBeingEdited.EmployeeName, &employeeBeingEdited.JobTitle, &employeeBeingEdited.Department, &employeeBeingEdited.StartDate, &employeeBeingEdited.Phone, &employeeBeingEdited.Email, &employeeBeingEdited.Salary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="row-%v">
			<td>%v</td>
			<td><input name="new-employee-name" type="text" value="%v" data-bind:new-employee-name required>
			<td><input name="new-job-title" type="text" value="%v" data-bind:new-job-title required></td>
			<td><input name="new-department" type="text" value="%v" data-bind:new-department required></td>
			<td><input name="new-start-date" type="date" value="%v" data-bind:new-start-date required></td>
			<td><input name="new-phone" type="phone" value="%v" data-bind:new-phone required></td>
			<td><input name="new-email" type="email" value="%v" data-bind:new-email required></td>
			<td><input name="salary" type="number" value="%v" data-bind:new-salary required></td>
			<td><button data-on:click="@get('/employees')">Cancel</button></td>
			<td><button data-on:click="@patch('/employees?id=%v')">Update</button></td>
		</tr>
		`, id, employeeBeingEdited.Id, employeeBeingEdited.EmployeeName, employeeBeingEdited.JobTitle, employeeBeingEdited.Department, employeeBeingEdited.StartDate, employeeBeingEdited.Phone, employeeBeingEdited.Email, employeeBeingEdited.Salary, id))
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
		orderSignals := &Order{}
		if err := datastar.ReadSignals(r, orderSignals); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, err := db.Query("SELECT * FROM orders")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var orderRecords []Order
		for rows.Next() {
			var orderFields Order
			if err := rows.Scan(&orderFields.Id, &orderFields.CustomerId, &orderFields.ProductId, &orderFields.Quantity, &orderFields.Total, &orderFields.Progress, &orderFields.OrderDate); err != nil {
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/orders.html", "templates/styles.css"))
		t.Execute(w, map[string]interface{}{
			"data": orderRecords,
		})
	case "POST":
		var productPrice float32
		pid := r.FormValue("product-id")
		db.QueryRow("SELECT price FROM products WHERE id = ?", pid).Scan(&productPrice)
		//, quantity FROM products JOIN orders ON products.id = orders.product_id
		result, err := db.Exec("INSERT INTO orders (customer_id, product_id, quantity, total, progress, order_date) VALUES (?,?,?,?*?, ?, CURRENT_DATE())",
			r.FormValue("customer-id"),
			r.FormValue("product-id"),
			r.FormValue("quantity"),
			r.FormValue("quantity"), productPrice,
			r.FormValue("progress"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/orders", http.StatusFound)
	case "DELETE":
		result, err := db.Exec("DELETE FROM orders WHERE id=?", r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/orders")
	case "PATCH":
		orderPatchSignals := &OrderPatch{}
		if err := datastar.ReadSignals(r, orderPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var productPrice float32
		db.QueryRow("SELECT price FROM products WHERE id = ?", orderPatchSignals.NewProductId).Scan(&productPrice)
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE orders SET customer_id=?, product_id=?, quantity=?, progress=?, total=?*? WHERE id=?", orderPatchSignals.NewCustomerId, orderPatchSignals.NewProductId, orderPatchSignals.NewQuantity, orderPatchSignals.NewProgress, orderPatchSignals.NewQuantity, productPrice, r.FormValue("id"))
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/reviews.html", "templates/styles.css"))
		t.Execute(w, map[string]any{
			"data": reviewRecords,
		})
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
	switch r.Method {
	case "GET":
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

		var accountBalance []TrialBalance
		for _, v := range chartOfAccountsRecords {
			var accountBalanceRecord TrialBalance
			if err = db.QueryRow("SELECT account_code, (SELECT COALESCE(SUM(debit),0) FROM journal_entries WHERE debit_account = ? ORDER BY debit_account) - (SELECT COALESCE(SUM(credit),0) FROM journal_entries WHERE credit_account = ? ORDER BY credit_account) AS balance FROM chart_of_accounts WHERE account_code = ? ORDER BY account_code", v.AccountCode, v.AccountCode, v.AccountCode).Scan(&accountBalanceRecord.AccountCode, &accountBalanceRecord.UnadjustedBalance); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			accountBalance = append(accountBalance, accountBalanceRecord)
		}

		generalLedger := &GeneralLedger{}
		ledgerAccountCode := r.FormValue("ledgerAccountCode")
		for journal_index, _ := range journalEntryRecords {
			// there is a credit and a debit that need to go to their respective accounts on the ledger
			if accountMap[ledgerAccountCode] == journalEntryRecords[journal_index].DebitAccount {
				generalLedger.JournalEntryID = append(generalLedger.JournalEntryID, journalEntryRecords[journal_index].Id)
				generalLedger.AccountCode = append(generalLedger.AccountCode, journalEntryRecords[journal_index].DebitAccount)
				generalLedger.AccountName = append(generalLedger.AccountName, accountCodeMap[journalEntryRecords[journal_index].DebitAccount])
				if journal_index == 0 {
					generalLedger.Balance = append(generalLedger.Balance, journalEntryRecords[journal_index].Debit)
				} else {
					generalLedger.Balance[journal_index] = generalLedger.Balance[journal_index-1] + journalEntryRecords[journal_index].Debit
				}
			} else if accountMap[ledgerAccountCode] == journalEntryRecords[journal_index].CreditAccount {
				generalLedger.JournalEntryID = append(generalLedger.JournalEntryID, journalEntryRecords[journal_index].Id)
				generalLedger.AccountCode = append(generalLedger.AccountCode, journalEntryRecords[journal_index].CreditAccount)
				generalLedger.AccountName = append(generalLedger.AccountName, accountCodeMap[journalEntryRecords[journal_index].DebitAccount])
				if journal_index == 0 {
					generalLedger.Balance = append(generalLedger.Balance, journalEntryRecords[journal_index].Credit)
				} else {
					generalLedger.Balance[journal_index] = generalLedger.Balance[journal_index-1] + journalEntryRecords[journal_index].Credit
				}
			}
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
		balanceSheet.TimePeriod = "Q2-2026"
		var incomeStatement IncomeStatements
		incomeStatement.TimePeriod = "Q2-2026"
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
			case accountMap["Total Sales"]:
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

		cashFlowStatement := &CashFlowStatements{}
		cashFlowStatement.TimePeriod = "Q2-2026"

		// if the account is an asset, add it. query or for loop
		cashFlowStatement.CashSpentOnAssets = 0

		cashFlowStatement.OperatingCashFlow = incomeStatement.TotalSales - cashFlowStatement.CashSpentOnAssets

		// sum of all cash expenses
		cashFlowStatement.OperatingExpenses = 0

		jinja, err := gonja.FromFile("templates/finances.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := exec.NewContext(map[string]interface{}{
			"journalEntries":    journalEntryRecords,
			"trialBalance":      trialBalance,
			"balanceSheet":      balanceSheet,
			"chartOfAccounts":   chartOfAccountsRecords,
			"incomeStatement":   incomeStatement,
			"cashFlowStatement": cashFlowStatement,
			"generalLedger":     generalLedger,
		})
		jinja.Execute(w, data)
	case "POST":
		result, err := db.Exec("INSERT INTO journal_entries (entry_date, debit_account, credit_account, debit, credit, notes) VALUES (?, ?, ?, ?, ?, ?)", r.FormValue("date"), r.FormValue("debit-account"), r.FormValue("cerdit-account"), r.FormValue("debit"), r.FormValue("credit"), r.FormValue("notes"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("error: %v, %v", id, err.Error())
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func main() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Println("error: ", err)
	}
	cfg := mysql.NewConfig()
	cfg.User = config.Section("database").Key("User").String()
	cfg.Passwd = config.Section("database").Key("Passwd").String()
	cfg.Net = config.Section("database").Key("Net").String()
	cfg.Addr = config.Section("database").Key("Addr").String()
	cfg.DBName = config.Section("database").Key("DBName").String()

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

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
