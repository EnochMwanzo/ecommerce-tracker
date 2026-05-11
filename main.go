package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/starfederation/datastar-go/datastar"
	"gopkg.in/ini.v1"
)

var db *sql.DB
var rows *sql.Rows
var err error

type Customer struct {
	Id           int    `json:"id"`
	CustomerName string `json:"customerName"`
	Subscriber bool `json:"subscriber"`
	Cohort string `json:"cohort"`
	SignupDate string `signupDate`

	NewCustomerName string `json:"newCustomerName"`
	NewSubscriber bool `json:"newSubscriber"`
	NewCohort string `json:"newCohort"`
	NewSignupDate string `newSignupDate`

	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber   int `json: "pageNumber"`
}

type Product struct {
	Id          int    `json:"id"`
	ProductName string `json:"productName"`
	Description string `json:"description"`
	StockInWareouse int `json:"stockInWareouse"`
	StockAtManufacturer int `json:"stockAtManufacturer"`
	StockAtRetailers int `json:"stockAtRetailers"`
	Price float `json:"price"`
	ReorderPoint bool `json:"reorderPoint"`

	NewProductName string `json:"newProductName"`
	NewDescription string `json:"newDescription"`

	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber   int `json: "pageNumber"`
}

type Order struct {
	Id         int    `json:"Id"`
	CustomerId int    `json:"customerId"`
	ProductId  int    `json:"productId"`
	Quantity   int    `json:"quantity"`
	Total      int    `json:"total"`
	Progress   string `json:"progress"`

	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber   int `json: "pageNumber"`
}

type OrderPatch struct {
	NewCustomerId int    `json:"newCustomerId"`
	NewProductId  int    `json:"newProductId"`
	NewQuantity   int    `json:"newQuantity"`
	NewProgress   string `json:"newProgress"`
}

type Employee struct {
	Id           int    `json:"id"`
	EmployeeName string `json:"employeeName"`
	JobTitle     string `json:"jobTitle"`
	Department   string `json:"department"`
	StartDate    string `json:"startDate"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Salary float32 `json:"salary"`

	NewEmployeeName string `json:"newEmployeeName"`
	NewJobTitle     string `json:"newJobTitle"`
	NewDepartment   string `json:"newDepartment"`
	NewStartDate    string `json:"newStartDate"`
	NewPhone        string `json:"newPhone"`
	NewEmail        string `json:"newEmail"`
	NewSalary float32 `json:"newSalary"`

	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber   int `json: "pageNumber"`
}

type Review struct {
	OrderId      int    `json:"id"`
	Rating       int    `json:"rating"`
	Review       string `json:"review"`
	ItemsPerPage int    `json: "itemsPerPage"`
	PageNumber   int    `json: "pageNumber"`

	NewOrderId int    `json:"newId"`
	NewRating  int    `json:"newRating"`
	NewReview  string `json:"newReview"`
}

type Transactions struct {
	Id int `json:"id"`
	Counterparty string `json:"counterparty"`
	TransactionDate string `json:"transactionDate"`
	Amount float `json:"amount"`
	AccountType string `json:"accountType"`
	Category string `json:"rating"`
	notes string `json:"notes"`
}

type CashFlowStatements struct {
    TimePeriod string `json:"notes"`
    OperatingCashFlow float `json:"operatingCashFlow"`
    TotalSales float `json:"totalSales"`
    CashSpentOnAssets float `json:"cashSpentOnAssets"`
    OperatingExpenses float `json:"operatingExpenses"`
};

type IncomeStatements struct {
    TimePeriod string `json:"timePeriod"`
    TotalSales float `json:"totalSales"`
    CostOfGoodsSold float `json:"costOfGoodsSold"`
    Profit float `json:"profit"`
    PromotionExpenses float `json:"promotionExpenses"`
    SellingGeneralAdministraticeExpenses float `json:"sellingGeneralAdministraticeExpenses"`
    DepreciationAndAmoritization float `json:"depreciationAndAmoritization"`
};

type BalanceSheets struct {
    TimePeriod string `json:"timePeriod"`
    Cash float `json:"cash"`
    AccountsReceivable float `json:"accountsReceivable"`
    PrepaidExpenses float `json:"prepaidExpenses"`
    Inventory float `json:"inventory"`
    PropertyAndEquipment float `json:"propertyAndEquipment"`
    Goodwill float `json:"goodwill"`
    AccountsPayable float `json:"accountsPayable"`
    AccruedExpenses float `json:"accruedExpenses"`
    UnearnedRevenue float `json:"unearnedRevenue"`
    LongTermDebt float `json:"longTermDebt"`
};


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
		signals := &Customer{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM customers")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var customer_data []Customer
		for rows.Next() {
			var customer_field Customer
			if err := rows.Scan(&customer_field.Id, &customer_field.CustomerName, &customer_field.Subscriber, &customer_field.Cohort, &customer_field.SignupDate); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			customer_data = append(customer_data, customer_field)
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/customers.html"))
		ipp := max(5, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := customer_data[ipp*(min(pn-1, len(customer_data)/ipp)) : min(pn*ipp, len(customer_data))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO customers (customer_name, subscriber, cohort, signup_date) VALUES (?, ?, ?, ?)", r.FormValue("customer-name"), r.FormValue("subscriber"), r.FormValue("cohort"), r.FormValue("signup-date"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
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
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/customers")
	case "PATCH":
		signals := &Customer{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		if signals.NewCustomerName != "" {
			result, err := db.Exec("UPDATE customers SET customer_name=? subscriber=?, cohort=?, signupDate=? WHERE id=?", signals.NewCustomerName, r.FormValue("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			id, err := result.LastInsertId()
			if err != nil {
				fmt.Errorf("error: %v, %v", id, err)

			}
			sse.Redirect("/customers")
		}
	}
}

func editCustomer(w http.ResponseWriter, r *http.Request) {
	signals := &Customer{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="row-%v">
			<td>%v</td>
			<td><input name="customer-name" type="text" value="%v" data-bind:new-customer-name required></td>
			<td>
				<select name="subscriber" type="text" data-bind:subsrciber>
					<option value="false">false</option>
					<option value="true">true</option>
				</select>
			</td>
		    <td>
				<select name="cohort" data-bind:cohort>
					<option value="received">received</option>
					<option value="in-progress">in progress</option>
					<option value="delivered">delivered</option>
				</select>
			</td>
		    <td><input name="signup-date" data-bind:signup-date>%v</td>
			<td><button data-on:click="@get('/customers')">Cancel</button></td>
			<td><button data-on:click="@patch('/customers?id=%v')">Update</button></td>
		</tr>
		`, r.FormValue("id"), r.FormValue("id"), r.FormValue("name"), r.FormValue("id")))
}

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
		var customer_field Customer
		if err := rows.Scan(&customer_field.Id, &customer_field.CustomerName); err != nil {
			fmt.Errorf("search error: ", err.Error())
			return
		}
		results = append(results, customer_field)
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

func products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		signals := &Product{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var product_data []Product
		for rows.Next() {
			var product_field Product
			if err := rows.Scan(&product_field.Id, &product_field.ProductName, &product_field.Description); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			product_data = append(product_data, product_field)
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/products.html"))
		ipp := max(5, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := product_data[ipp*(min(pn-1, len(product_data)/ipp)) : min(pn*ipp, len(product_data))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO products (product_name, product_description) VALUES (?,?)", r.FormValue("product-name"), r.FormValue("product-description"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		http.Redirect(w, r, "/products", http.StatusFound)
	case "DELETE":
		result, err := db.Exec("DELETE FROM products WHERE id=?", r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/products")
	case "PATCH":
		signals := &Product{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE products SET product_name=?, product_description=? WHERE id=?", signals.NewProductName, signals.NewDescription, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse.Redirect("/products")
	}
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	signals := &Product{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="row-%v">
			<td>%v</td>
			<td><input name="new-product-name" type="text" value="%v" data-bind:new-product-name required></td>
			<td><input name="product-description" type="text" value="%v" data-bind:new-description required></td>
			<td><button data-on:click="@get('/products')">Cancel</button></td>
			<td><button data-on:click="@patch('/products?id=%v')">Update</button></td>
		</tr>
		`, signals.Id, signals.Id, signals.ProductName, signals.Id))
}

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
		var product_field Product
		if err := rows.Scan(&product_field.Id, &product_field.ProductName, product_field.Description); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		results = append(results, product_field)
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

func employees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		signals := &Employee{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM employees")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var employee_data []Employee
		for rows.Next() {
			var employee_field Employee
			if err := rows.Scan(&employee_field.Id, &employee_field.EmployeeName, &employee_field.JobTitle, &employee_field.Department, &employee_field.StartDate, &employee_field.Phone, &employee_field.Email, &employee_field.Salary); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			employee_data = append(employee_data, employee_field)
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/employees.html"))
		ipp := max(5, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := employee_data[ipp*(min(pn-1, len(employee_data)/ipp)) : min(pn*ipp, len(employee_data))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO employees (employee_name, job_title, department, phone, company_email, salary) VALUES (?,?,?,?,?,?)", r.FormValue("employee-name"), r.FormValue("job-title"), r.FormValue("department"), r.FormValue("phone"), r.FormValue("company-email"), r.FormValue("salary"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
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
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/employees")
	case "PATCH":
		signals := &Employee{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE employees SET employee_name=?, job_title=?, department=?, phone=?, company_email=?, salary=? WHERE id=?", signals.NewEmployeeName, signals.NewJobTitle, signals.NewDepartment, signals.NewPhone, signals.NewEmail, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse.Redirect("/employees")
	}
}

func editEmployee(w http.ResponseWriter, r *http.Request) {
	signals := &Employee{}
	if err := datastar.ReadSignals(r, signals); err != nil {
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
			<td><input name="new-phone" type="phone" value="%v" data-bind:new-phone required></td>
			<td><input name="new-email" type="email" value="%v" data-bind:new-email required></td>
			<td><input name="salary" type="number" value="%v" data-bind:new-salary required></td>
			<td><button data-on:click="@get('/employees')">Cancel</button></td>
			<td><button data-on:click="@patch('/employees?id=%v')">Update</button></td>
		</tr>
		`, signals.EmployeeName, signals.JobTitle, signals.Department, signals.StartDate, signals.Phone, signals.Email))
}

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
			var employee_field Employee
			if err := rows.Scan(&employee_field.Id, &employee_field.EmployeeName, &employee_field.JobTitle, &employee_field.Department, &employee_field.StartDate, &employee_field.Phone, &employee_field.Email, &employee_field.Salary); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			results = append(results, employee_field)
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

func orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		signals := &Order{}
		signalsToBePatched := &OrderPatch{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := datastar.ReadSignals(r, signalsToBePatched); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM orders")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var order_data []Order
		for rows.Next() {
			var order_field Order
			if err := rows.Scan(&order_field.Id, &order_field.CustomerId, &order_field.ProductId, &order_field.Quantity, &order_field.Total, &order_field.Progress); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			order_data = append(order_data, order_field)
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/orders.html"))
		ipp := max(5, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := order_data[ipp*(min(pn-1, len(order_data)/ipp)) : min(pn*ipp, len(order_data))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO orders (customer_id, product_id, quantity, progress) VALUES (?,?,?,?)", r.FormValue("customer-id"), r.FormValue("product-id"), r.FormValue("quantity"), r.FormValue("progress"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		http.Redirect(w, r, "/orders", http.StatusFound)
	case "DELETE":
		result, err := db.Exec("DELETE FROM orders WHERE order_id=?", r.FormValue("order-id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/orders")
	case "PATCH":
		signalsToBePatched := &OrderPatch{}
		if err := datastar.ReadSignals(r, signalsToBePatched); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE orders SET customer_id=?, product_id=?, quantity=?, progress=? WHERE id=?", signalsToBePatched.NewCustomerId, signalsToBePatched.NewProductId, signalsToBePatched.NewQuantity, signalsToBePatched.NewProgress, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse.Redirect("/orders")
	}
}

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
			var order_field Order
			if err := rows.Scan(&order_field.Id, &order_field.CustomerId, &order_field.ProductId, &order_field.Quantity, &order_field.Total, &order_field.Progress); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			results = append(results, order_field)
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

func editOrder(w http.ResponseWriter, r *http.Request) {
	signals := &Order{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	signalsToBePatched := &OrderPatch{}
	if err := datastar.ReadSignals(r, signalsToBePatched); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="row-%v">
			<td>%v</td>
			<td><input name="new-customer-id" type="number" value="%v" data-bind:new-customer-id required>
			<td><input name="new-product-id" type="number" value="%v" min="1" max="5" data-bind:new-product-id required></td>
			<td><input name="new-quantity" type="text" value="%v" data-bind:new-quantity required></td>
			<td>%v</td>
			<td>
				<select data-bind:new-progress>
					<option value="received">received</option>
					<option value="in progress">in progress</option>
					<option value="delivered">delivered</option>
				</select>
			</td>
			<td><button data-on:click="@get('/orders')">Cancel</button></td>
			<td><button data-on:click="@patch('/orders?id=%v')">Update</button></td>
		</tr>
		`, r.FormValue("id"), signals.CustomerId, signals.ProductId, signals.Id))
}

func reviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		signals := &Review{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM reviews")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var review_data []Review
		for rows.Next() {
			var review_field Review
			if err := rows.Scan(&review_field.OrderId, &review_field.Rating, &review_field.Review); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			review_data = append(review_data, review_field)
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
		t := template.Must(template.ParseFiles("templates/base.html", "templates/reviews.html"))
		ipp := max(5, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := review_data[ipp*(min(pn-1, len(review_data)/ipp)) : min(pn*ipp, len(review_data))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO reviews (order_id, rating, review) VALUES (?,?,?)", r.FormValue("order-id"), r.FormValue("rating"), r.FormValue("review"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
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
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse := datastar.NewSSE(w, r)
		sse.Redirect("/reviews")
	case "PATCH":
		signals := &Review{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE reviews SET order_id=?, rating=?, review=? WHERE order_id=?", signals.NewOrderId, signals.NewRating, signals.NewReview, r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse.Redirect("/reviews")
	}
}

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
			var review_field Review
			if err := rows.Scan(&review_field.OrderId, &review_field.Rating, &review_field.Review); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			results = append(results, review_field)
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

func finances(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		transactionSignals := &Transactions{}
		if err := datastar.ReadSignals(r, signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM balance_sheets JOIN cash_flow_statements JOIN income_statements")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var balance_sheet_data []BalanceSheets
		var cash_flow_statement_data []CashFlowStatements
		var income_statement_data []IncomeStatements
		for rows.Next() {
			var balance_sheet_field []BalanceSheets
			var cash_flow_statement_field []CashFlowStatements
			var income_statement_field []IncomeStatements
			if err := rows.Scan(
				&balance_sheet_field.TimePeriod, &balance_sheet_field.Cash, &balance_sheet_field.AccountsReceivable, &balance_sheet_field.PrepaidExpenses, &balance_sheet_field.Inventory, &balance_sheet_field.PropertyAndEquipment, &balance_sheet_field.Goodwill, &balance_sheet_field.AccountsPayable, &balance_sheet_field.AccruedExpenses, &balance_sheet_field.UnearnedRevenue, &balance_sheet_field.LongTermDebt, &cash_flow_statement_field.TimePeriod, &cash_flow_statement_field.OperatingCashFlow, &cash_flow_statement_field.TotalSales, &cash_flow_statement_field.CashSpentOnAssets, &cash_flow_statement_field.OperatingExpenses,
				&income_statement_field.TimePeriod, &income_statement_field.TotalSales, &income_statement_field.CostOfGoodsSold, &income_statement_field.Profit, &income_statement_field.PromotionExpenses, &income_statement_field.SellingGeneralAdministraticeExpenses, &income_statement_field.DepreciationAndAmoritization,
			); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			balance_sheet_data = append(balance_sheet_data, balance_sheet_field)
			cash_flow_statement_data = append(cash_flow_statement_data, cash_flow_statement_field)
			income_statement_data = append(income_statement_data, income_statement_field)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = rows.Close()
		if rowErr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM transactions")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var transaction_data []Transaction
		for rows.Next(){
			var transaction_field Transaction
			if err := rows.Scan(&transaction_field.Id, &transaction_field.Counterparty, &transaction_field.TransactionDate, &transaction_field.Amount, &transaction_field.AccountType, &transaction_field.Category, )
		}
	}
}

func main() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("error: ", err)
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
		log.Fatalf("error: ", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf("error: ", pingErr)
	}
	fmt.Println("Connected to MySQL")

	http.HandleFunc("/", index)

	http.HandleFunc("/customers", customers)
	http.HandleFunc("/customers/search", searchCustomers)
	http.HandleFunc("/customers/edit", editCustomer)

	http.HandleFunc("/products", products)
	http.HandleFunc("/products/search", searchProducts)
	http.HandleFunc("/products/edit", editProduct)

	http.HandleFunc("/orders", orders)
	http.HandleFunc("/orders/search", searchOrders)
	http.HandleFunc("/orders/edit", editOrder)

	http.HandleFunc("/employees", employees)
	http.HandleFunc("/employees/search", searchEmployees)
	http.HandleFunc("/employees/edit", editEmployee)

	http.HandleFunc("/reviews", reviews)
	http.HandleFunc("/reviews/search", searchReviews)
	http.HandleFunc("/reviews/edit", editReview)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
