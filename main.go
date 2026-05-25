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
	SignupDate string `json:"signupDate"`
}
type CustomerPatch struct {
	NewCustomerName string `json:"newCustomerName"`
	NewSubscriber bool `json:"newSubscriber"`
	NewCohort string `json:"newCohort"`
	NewSignupDate string `json:"newSignupDate"`
}
type Product struct {
		Id          int    `json:"id"`
		ProductName string `json:"productName"`
		Description string `json:"description"`
		StockAtWarehouse int `json:"stockAtWarehouse"`
		StockAtManufacturer int `json:"stockAtManufacturer"`
		Price float32 `json:"price"`
		ReorderPoint int `json:"reorderPoint"`
		LeadTimeDays int `json:"leadTimeDays"`
	}
	type ProductPatch struct {
		NewProductName string `json:"newProductName"`
		NewDescription string `json:"newDescription"`
		NewStockAtWarehouse int `json:"newStockAtWarehouse"`
		NewStockAtManufacturer int `json:"newStockAtManufacturer"`
		NewPrice float32 `json:"newPrice"`
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
}
type EmployeePatch struct {
	NewEmployeeName string `json:"newEmployeeName"`
	NewJobTitle     string `json:"newJobTitle"`
	NewDepartment   string `json:"newDepartment"`
	NewStartDate    string `json:"newStartDate"`
	NewPhone        string `json:"newPhone"`
	NewEmail        string `json:"newEmail"`
	NewSalary float32 `json:"newSalary"`
}
type Order struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customerId"`
	ProductId  int    `json:"productId"`
	Quantity   int    `json:"quantity"`
	Total      float32    `json:"total"`
	Progress   string `json:"progress"`
	OrderDate   string `json:"orderDate"`
}
type OrderPatch struct {
	NewCustomerId int    `json:"newCustomerId"`
	NewProductId  int    `json:"newProductId"`
	NewQuantity   int    `json:"newQuantity"`
	NewProgress   string `json:"newProgress"`
	NewOrderDate  string `json:"newOrderDate"`
}
type Review struct {
	OrderId      int    `json:"id"`
	Rating       int    `json:"rating"`
	Review       string `json:"review"`
}
type Transaction struct {
	Id int `json:"id"`
	Counterparty string `json:"counterparty"`
	TransactionDate string `json:"transactionDate"`
	Amount float32 `json:"amount"`
	AccountType string `json:"accountType"`
	Category string `json:"rating"`
	Notes string `json:"notes"`
}
type TransactionPatch struct {
	NewCounterparty string `json:"newCounterparty"`
	NewTransactionDadete string `json:"newTransactionDate"`
	NewAmount float32 `json:"newAmount"`
	NewAccountType string `json:"newAccountType"`
	NewCategory string `json:"newRating"`
	NewNotes string `json:"newNotes"`
}
type CashFlowStatements struct {
    TimePeriod string `json:"notes"`
    OperatingCashFlow float32 `json:"operatingCashFlow"`
    TotalSales float32 `json:"totalSales"`
    CashSpentOnAssets float32 `json:"cashSpentOnAssets"`
    OperatingExpenses float32 `json:"operatingExpenses"`
}
type IncomeStatements struct {
    TimePeriod string `json:"timePeriod"`
    TotalSales float32 `json:"totalSales"`
    CostOfGoodsSold float32 `json:"costOfGoodsSold"`
    Profit float32 `json:"profit"`
    PromotionExpenses float32 `json:"promotionExpenses"`
    SellingAndGeneralAdministrativeExpenses float32 `json:"sellingAndGeneralAdministrativeExpenses"`
    DepreciationAndAmoritization float32 `json:"depreciationAndAmoritization"`
}
type BalanceSheets struct {
    TimePeriod string `json:"timePeriod"`
    Cash float32 `json:"cash"`
    AccountsReceivable float32 `json:"accountsReceivable"`
    PrepaidExpenses float32 `json:"prepaidExpenses"`
    Inventory float32 `json:"inventory"`
    PropertyAndEquipment float32 `json:"propertyAndEquipment"`
    Goodwill float32 `json:"goodwill"`
    AccountsPayable float32 `json:"accountsPayable"`
    AccruedExpenses float32 `json:"accruedExpenses"`
    UnearnedRevenue float32 `json:"unearnedRevenue"`
    LongTermDebt float32 `json:"longTermDebt"`
}
type Pagination struct {
	ItemsPerPage int `json:"itemsPerPage"`
	PageNumber   int `json:"pageNumber"`
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
		pageSignals := &Pagination{}
		if err := datastar.ReadSignals(r, pageSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		customerSignals := &Customer{}
		if err := datastar.ReadSignals(r, customerSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM customers")
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t := template.Must(template.ParseFiles("templates/base.html", "templates/customers.html"))
		ipp := max(5, pageSignals.ItemsPerPage)
		pn := max(1, pageSignals.PageNumber)
		page := customerRecords[ipp*(min(pn-1, len(customerRecords)/ipp)) : min(pn*ipp, len(customerRecords))]
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
		customerPatchSignals := &CustomerPatch{}
		if err := datastar.ReadSignals(r, customerPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		if customerPatchSignals.NewCustomerName != "" {
			result, err := db.Exec("UPDATE customers SET customer_name=? subscriber=?, cohort=?, signupDate=? WHERE id=?", customerPatchSignals.NewCustomerName, customerPatchSignals.NewSubscriber, customerPatchSignals.NewCohort, customerPatchSignals.NewSignupDate, r.FormValue("id"))
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
	type CustomerRow struct {
		Customer Customer `json:"row18"`
	}
	signals := &CustomerRow{}
	if err := datastar.ReadSignals(r, signals); err != nil {
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
					<option value="0">engaged prospect</option>
					<option value="0">engaged prospect</option>
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
		`, r.FormValue("id"), signals.Customer.Id, signals.Customer.CustomerName, signals.Customer.Subscriber, signals.Customer.Cohort, signals.Customer.SignupDate, signals.Customer.Id))
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
		var customerFields Customer
		if err := rows.Scan(&customerFields.Id, &customerFields.CustomerName); err != nil {
			fmt.Errorf("search error: ", err.Error())
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

func products(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	pageSignals := &Pagination{}
	if err := datastar.ReadSignals(r, pageSignals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		productRecords = append(productRecords, productFields)
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
		ipp := max(5, pageSignals.ItemsPerPage)
		pn := max(1, pageSignals.PageNumber)
		page := productRecords[ipp*(min(pn-1, len(productRecords)/ipp)) : min(pn*ipp, len(productRecords))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO products (product_name, description, price, stock_at_warehouse, stock_at_manufacturer) VALUES (?,?,?,?,?)", r.FormValue("product-name"), r.FormValue("description"), r.FormValue("price"), r.FormValue("stock-at-warehouse"), r.FormValue("stock-at-manufacturer"))
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
		productPatchSignals := &ProductPatch{}
		if err := datastar.ReadSignals(r, productPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE products SET product_name=?, description=?, price=?, stock_at_warehouse=?, stock_at_manufacturer=? WHERE id=?", productPatchSignals.NewProductName, productPatchSignals.NewDescription, productPatchSignals.NewPrice, productPatchSignals.NewStockAtWarehouse, productPatchSignals.NewStockAtManufacturer, r.FormValue("id"))
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse.Redirect("/products")
	}
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	type ProductRow struct {
			Product Product `json:"row2"`
	}
	signals := &ProductRow{}
	if err := datastar.ReadSignals(r, signals); err != nil {
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
		`, r.FormValue("id"), signals.Product.Id, signals.Product.ProductName, signals.Product.Description, signals.Product.StockAtWarehouse, signals.Product.StockAtManufacturer, signals.Product.Price, signals.Product.ReorderPoint, signals.Product.LeadTimeDays, signals.Product.Id))
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
		var productFields Product
		if err := rows.Scan(&productFields.Id, &productFields.ProductName, &productFields.Description, &productFields.StockAtWarehouse, &productFields.StockAtManufacturer, &productFields.Price, &productFields.ReorderPoint, &productFields.LeadTimeDays, ); err != nil {
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

func employees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	pageSignals := &Pagination{}
	if err := datastar.ReadSignals(r, pageSignals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	employeeSignals := &Employee{}
	if err := datastar.ReadSignals(r, employeeSignals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t := template.Must(template.ParseFiles("templates/base.html", "templates/employees.html"))
		ipp := max(5, pageSignals.ItemsPerPage)
		pn := max(1, pageSignals.PageNumber)
		page := employeeRecords[ipp*(min(pn-1, len(employeeRecords)/ipp)) : min(pn*ipp, len(employeeRecords))]
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
			fmt.Errorf("error: %v, %v", id, err)
		}
		sse.Redirect("/employees")
	}
}

func editEmployee(w http.ResponseWriter, r *http.Request) {
	type EmployeeRow struct {
			Employee Employee `json:"row1"`
	}
	signals := &EmployeeRow{}
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
			<td><input name="new-start-date" type="date" value="%v" data-bind:new-start-date required></td>
			<td><input name="new-phone" type="phone" value="%v" data-bind:new-phone required></td>
			<td><input name="new-email" type="email" value="%v" data-bind:new-email required></td>
			<td><input name="salary" type="number" value="%v" data-bind:new-salary required></td>
			<td><button data-on:click="@get('/employees')">Cancel</button></td>
			<td><button data-on:click="@patch('/employees?id=%v')">Update</button></td>
		</tr>
		`,r.FormValue("id"), signals.Employee.Id, signals.Employee.EmployeeName, signals.Employee.JobTitle, signals.Employee.Department, signals.Employee.StartDate, signals.Employee.Phone, signals.Employee.Email, signals.Employee.Salary, r.FormValue("id")))
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

func orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		orderSignals := &Order{}
		if err := datastar.ReadSignals(r, orderSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pageSignals := &Pagination{}
		if err := datastar.ReadSignals(r, pageSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM orders")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rowErr := rows.Close()
		if rowErr != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t := template.Must(template.ParseFiles("templates/base.html", "templates/orders.html"))
		ipp := max(5, pageSignals.ItemsPerPage)
		pn := max(1, pageSignals.PageNumber)
		page := orderRecords[ipp*(min(pn-1, len(orderRecords)/ipp)) : min(pn*ipp, len(orderRecords))]
		t.Execute(w, map[string]interface{}{
			"data": page,
		})
	case "POST":
		result, err := db.Exec("INSERT INTO orders (customer_id, product_id, quantity, progress, order_date) VALUES (?,?,?,?, CURRENT_DATE())", r.FormValue("customer-id"), r.FormValue("product-id"), r.FormValue("quantity"), r.FormValue("progress"))
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
		result, err := db.Exec("DELETE FROM orders WHERE order_id=?", r.FormValue("id"))
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
		orderPatchSignals := &OrderPatch{}
		if err := datastar.ReadSignals(r, orderPatchSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sse := datastar.NewSSE(w, r)
		result, err := db.Exec("UPDATE orders SET customer_id=?, product_id=?, quantity=?, progress=? WHERE id=?", orderPatchSignals.NewCustomerId, orderPatchSignals.NewProductId, orderPatchSignals.NewQuantity, orderPatchSignals.NewProgress, r.FormValue("id"))
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

func editOrder(w http.ResponseWriter, r *http.Request) {
	type OrderRow struct {
			Order Order `json:"row1"`
	}
	signals := &OrderRow{}
	if err := datastar.ReadSignals(r, signals); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`
		<tr id="row-%v">
			<td>%v</td>
			<td><input name="new-customer-id" type="number" value="%v" data-bind:new-customer-id required>
			<td><input name="new-product-id" type="number" value="%v" data-bind:new-product-id required></td>
			<td><input name="new-quantity" type="text" value="%v" data-bind:new-quantity required></td>
			<td>%v</td>
			<td>
				<select name="new-progress" data-bind:new-progress>
					<option value="%v">keep original</option>
					<option value="received">received</option>
					<option value="in-progress">in progress</option>
					<option value="delivered">delivered</option>
				</select>
			</td>
			<td><input name="new-order-date" type="date" value="%v" data-bind:new-order-date required></td>
			<td><button data-on:click="@get('/orders')">Cancel</button></td>
			<td><button data-on:click="@patch('/orders?id=%v')">Update</button></td>
		</tr>
		`, r.FormValue("id"), signals.Order.Id, signals.Order.CustomerId, signals.Order.ProductId, signals.Order.Quantity, signals.Order.Total, signals.Order.Progress, signals.Order.OrderDate, r.FormValue("id")))
}

func reviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		reviewSignals := &Review{}
		if err := datastar.ReadSignals(r, reviewSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pageSignals := &Pagination{}
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t := template.Must(template.ParseFiles("templates/base.html", "templates/reviews.html"))
		ipp := max(5, pageSignals.ItemsPerPage)
		pn := max(1, pageSignals.PageNumber)
		page := reviewRecords[ipp*(min(pn-1, len(reviewRecords)/ipp)) : min(pn*ipp, len(reviewRecords))]
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

// func editReview(w http.ResponseWriter, r *http.Request) {
// 	signals := &Review{}
// 	if err := datastar.ReadSignals(r, signals); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	sse := datastar.NewSSE(w, r)
// 	sse.PatchElements(fmt.Sprintf(`
// 		<tr id="row-%v">
// 			<td>%v</td>
// 			<td><input name="new-order-id" type="number" value="%v" data-bind:order-id required>
// 			<td><input name="new-rating" type="number" value="%v" min="1" max="5" data-bind:rating required></td>
// 			<td><input name="new-review" type="text" value="%v" data-bind:review required></td>
// 			<td><button data-on:click="@get('/orders')">Cancel</button></td>
// 			<td><button data-on:click="@patch('/orders?id=%v')">Update</button></td>
// 		</tr>
// 		`, signals.NewOrderId, signals.NewRating, signals.NewReview, signals.OrderId))
// }

func finances(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		transactionSignals := &Transaction{}
		if err := datastar.ReadSignals(r, transactionSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pageSignals := &Pagination{}
		if err := datastar.ReadSignals(r, pageSignals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := db.Query("SELECT * FROM balance_sheets JOIN cash_flow_statements JOIN income_statements")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var balanceSheetRecords []BalanceSheets
		var cashFlowStatementRecords []CashFlowStatements
		var incomeStatementRecords []IncomeStatements
		for rows.Next() {
			var balanceSheetFields BalanceSheets
			var cashFlowStatementFields CashFlowStatements
			var incomeStatementFields IncomeStatements
			if err := rows.Scan(&balanceSheetFields.TimePeriod, &balanceSheetFields.Cash, &balanceSheetFields.AccountsReceivable, &balanceSheetFields.PrepaidExpenses, &balanceSheetFields.Inventory, &balanceSheetFields.PropertyAndEquipment, &balanceSheetFields.Goodwill, &balanceSheetFields.AccountsPayable, &balanceSheetFields.AccruedExpenses, &balanceSheetFields.UnearnedRevenue, &balanceSheetFields.LongTermDebt, &cashFlowStatementFields.TimePeriod, &cashFlowStatementFields.OperatingCashFlow, &cashFlowStatementFields.TotalSales, &cashFlowStatementFields.CashSpentOnAssets, &cashFlowStatementFields.OperatingExpenses,&incomeStatementFields.TimePeriod, &incomeStatementFields.TotalSales, &incomeStatementFields.CostOfGoodsSold, &incomeStatementFields.Profit, &incomeStatementFields.PromotionExpenses, &incomeStatementFields.SellingAndGeneralAdministrativeExpenses, &incomeStatementFields.DepreciationAndAmoritization,); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			balanceSheetRecords = append(balanceSheetRecords, balanceSheetFields)
			cashFlowStatementRecords = append(cashFlowStatementRecords, cashFlowStatementFields)
			incomeStatementRecords = append(incomeStatementRecords, incomeStatementFields)
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
		rows, err = db.Query("SELECT * FROM transactions")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		var transactionRecords []Transaction
		for rows.Next() {
			var transactionFields Transaction
			if err := rows.Scan(&transactionFields.Id, &transactionFields.Counterparty, &transactionFields.TransactionDate, &transactionFields.Amount, &transactionFields.AccountType, &transactionFields.Category, &transactionFields.Notes); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			transactionRecords = append(transactionRecords, transactionFields)
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
		t := template.Must(template.ParseFiles("templates/finances.html"))
		ipp := max(5, pageSignals.ItemsPerPage)
		pn := max(1, pageSignals.PageNumber)
		transactionPage := transactionRecords[ipp*(min(pn-1, len(transactionRecords)/ipp)) : min(pn*ipp, len(transactionRecords))]
		t.Execute(w, map[string]interface{}{
			"transactionPage": transactionPage,
			"incomeStatements": incomeStatementRecords,
			"balanceSheets": balanceSheetRecords,
			"cashFlowStatements": cashFlowStatementRecords,
		})
		case "POST":
		result, err := db.Exec("INSERT INTO transactions (counterparty, transaction_date, amount, account_type, category, notes) VALUES (?, ?, ?, ?, ?, ?)", r.FormValue("counterparty"), r.FormValue("transaction-date"), r.FormValue("amount"), r.FormValue("account-type"), r.FormValue("category"), r.FormValue("notes"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Errorf("error: %v, %v", id, err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
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
	// http.HandleFunc("/reviews/edit", editReview)

	http.HandleFunc("/finances", finances)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
