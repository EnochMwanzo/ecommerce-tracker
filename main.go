package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/starfederation/datastar-go/datastar"
	"gopkg.in/ini.v1"
)

var db *sql.DB

type Customer struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	NewName string `json:"newName"`
	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber int `json: "pageNumber"`
}

type Product struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	Description string `json:"Description"`
	NewName string `json:"newName"`
	NewDescription string `json:"newDescription"`
	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber int `json: "PageNumber"`
}

type Order struct {
	ID int `json:"id"`
	CustomerID int `json:"CustomerId"`
	ProductID int `json:"ProductId"`
	Quantity int `json:"quantity"`
	Total int `json:"total"`
	Progress string `json:"progress"`
	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber int `json: "PageNumber"`
}

type employees struct {
	ID int `json:"id"`
	EmployeeName string `json:"employeeName"`
	JobTitle string `json:"jobTitle"`
	Department string `json:"department"`
	DaysSinceStarting int `json:"daysSinceStarting"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber int `json: "PageNumber"`
}

type Search struct {
	Query string `json:"query"`
	SearchBy string `json:"searchBy"`
	ItemsPerPage int `json: "itemsPerPage"`
	PageNumber int `json: "PageNumber"`
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
				fmt.Errorf("error: ", err)
			}
			defer rows.Close()
			var customer_data []Customer
			for rows.Next() {
		      	var customer_field Customer
		       	if err := rows.Scan(&customer_field.ID, &customer_field.Name); err != nil {
		        	fmt.Errorf("error: ", err)
		        }
		        customer_data = append(customer_data, customer_field)
		   	}
		   	if err := rows.Err(); err != nil {
		        fmt.Errorf("error: ", err)
		    }
			rerr := rows.Close()
				if rerr != nil {
					fmt.Errorf("error: ", err)
				}
			t := template.Must(template.ParseFiles("templates/customers.html"))
			ipp := max(5, signals.ItemsPerPage)
			pn := max(1, signals.PageNumber)
			page := customer_data[ipp*(min(pn - 1, len(customer_data)/ipp)):min(pn*ipp, len(customer_data))]
			t.Execute(w, map[string]interface{}{
				"data": page,
			})
		case "POST":
			result, err := db.Exec("INSERT INTO customers (customer_name) VALUES (?)", r.FormValue("customer_name"))
		    if err != nil {
		        fmt.Errorf("error: ", err)
		    }
		    id, err := result.LastInsertId()
		    if err != nil {
		        fmt.Errorf("error: %v, %v", id, err)
		    }
			http.Redirect(w, r, "/", http.StatusFound)
		case "DELETE":
			result, err := db.Exec("DELETE FROM customers WHERE id=?", r.FormValue("id"))
			    if err != nil {
			        fmt.Errorf("error: ", err)
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
			if signals.NewName != "" {
				result, err := db.Exec("UPDATE customers SET customer_name=? WHERE id=?", signals.NewName, r.FormValue("id"))
				    if err != nil {
				        fmt.Errorf("error: ", err)
				    }
				id, err := result.LastInsertId()
				    if err != nil {
				        fmt.Errorf("error: %v, %v", id, err)

				    }
				sse.Redirect("/customers")
			}
	}
}

func editCustomers(w http.ResponseWriter, r *http.Request) {
	signals := &Customer{}
	if err := datastar.ReadSignals(r, signals); err != nil {
	    http.Error(w, err.Error(), http.StatusBadRequest)
	    return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`<tr id="row-%v"><td>%v</td><td><input name="customer-name" type="text" value="%v" data-bind:new-name required></td><td><button data-on:click="@get('/customers')">Cancel</button></td><td><button data-on:click="@patch('/customers?id=%v')">Update</button></td></form></tr>`, r.FormValue("id"), r.FormValue("id"), r.FormValue("name"), r.FormValue("id")))
}

func searchCustomers(w http.ResponseWriter, r *http.Request){
	signals := &Search{}
	if err := datastar.ReadSignals(r, signals); err != nil {
	    http.Error(w, err.Error(), http.StatusBadRequest)
	    return
	}
	switch signals.SearchBy {
		case "Name":
		pattern := "%" + signals.Query + "%"
		rows, err := db.Query("SELECT * FROM customers WHERE customer_name LIKE ?", pattern)
	    if err != nil {
	        fmt.Errorf("error: ", err)
	    }
		defer rows.Close()
		var results []Customer
		for rows.Next() {
	      	var customer_field Customer
	       	if err := rows.Scan(&customer_field.ID, &customer_field.Name); err != nil {
	        	fmt.Errorf("search error: ", err.Error())
	         	return
	        }
	        results = append(results, customer_field)
		}
		ipp := max(1, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := results[ipp*(min(pn - 1, len(results)/ipp)):min(pn*ipp, len(results))]
		t, err := template.New("results").Parse(`<tbody id="current-table">{{range .pages}}<tr id="row-{{.ID}}"><td data-signals="{ID: {{.ID}} }">{{.ID}}</td><td data-signals="{Name: '{{.Name}}'}">{{.Name}}</td><td><button data-on:click="confirm('Are you sure?') && @delete('/customers?id={{.ID}}')">Delete</button></td><td><button data-on:click="@get('/customers/edit?id={{.ID}}&name={{.Name}}')">Edit</button></td></tr>{{end}}</tbody>`)
		if err != nil {
      		fmt.Errorf("error: ", err)
	    }
	    var builder strings.Builder
	    res := map[string][]Customer{"pages": page}
	    t.Execute(&builder, res)
	    searchResult := builder.String()
	    sse := datastar.NewSSE(w, r)
		sse.PatchElements(searchResult)
	}
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
				fmt.Errorf("error: ", err)
			}
			defer rows.Close()
			var product_data []Product
			for rows.Next() {
		      	var product_field Product
		       	if err := rows.Scan(&product_field.ID, &product_field.Name, &product_field.Description); err != nil {
		        	fmt.Errorf("error: ", err)
		        }
		        product_data = append(product_data, product_field)
		   	}
		   	if err := rows.Err(); err != nil {
		        fmt.Errorf("error: ", err)
		    }
			err = rows.Close()
				if err != nil {
					fmt.Errorf("error: ", err)
				}
			t := template.Must(template.ParseFiles("templates/base.html", "templates/products.html"))
			ipp := max(5, signals.ItemsPerPage)
			pn := max(1, signals.PageNumber)
			page := product_data[ipp*(min(pn - 1, len(product_data)/ipp)):min(pn*ipp, len(product_data))]
			t.Execute(w, map[string]interface{}{
				"data": page,
			})
		case "POST":
			result, err := db.Exec("INSERT INTO products (product_name, product_description) VALUES (?,?)", r.FormValue("product-name"), r.FormValue("product-description"))
		    if err != nil {
		        fmt.Errorf("error: ", err)
		    }
		    id, err := result.LastInsertId()
		    if err != nil {
		        fmt.Errorf("error: %v, %v", id, err)
		    }
			http.Redirect(w, r, "/products", http.StatusFound)
		case "DELETE":
			result, err := db.Exec("DELETE FROM products WHERE id=?", r.FormValue("id"))
			    if err != nil {
			        fmt.Errorf("error: ", err)
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
			result, err := db.Exec("UPDATE products SET product_name=?, product_description=? WHERE id=?", signals.NewName, signals.NewDescription, r.FormValue("id"))
			    if err != nil {
			        fmt.Errorf("error: ", err)
			    }
			id, err := result.LastInsertId()
			    if err != nil {
			        fmt.Errorf("error: %v, %v", id, err)
			    }
			sse.Redirect("/products")
	}
}

func editProducts(w http.ResponseWriter, r *http.Request) {
	signals := &Product{}
	if err := datastar.ReadSignals(r, signals); err != nil {
	    http.Error(w, err.Error(), http.StatusBadRequest)
	    return
	}
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(fmt.Sprintf(`<tr id="row-%v"><td>%v</td><td><input name="product-name" type="text" value="%v" data-bind:new-name required><td><input name="product-description" type="text" value="%v" data-bind:new-description required></td><td><button data-on:click="@get('/products')">Cancel</button></td><td><button data-on:click="@patch('/products?id=%v')">Update</button></td></form></tr>`, signals.ID, signals.ID, signals.Name, signals.ID))
}

func searchProducts(w http.ResponseWriter, r *http.Request){
	signals := &Search{}
	if err := datastar.ReadSignals(r, signals); err != nil {
	    http.Error(w, err.Error(), http.StatusBadRequest)
	    return
	}
	fmt.Println("signals: ", signals)
	switch signals.SearchBy {
		case "Name":
		pattern := "%" + signals.Query + "%"
		rows, err := db.Query("SELECT * FROM products WHERE product_name LIKE ?", pattern)
	    if err != nil {
	        fmt.Errorf("error: ", err)
	    }
		defer rows.Close()
		var results []Product
		for rows.Next() {
	      	var product_field Product
	       	if err := rows.Scan(&product_field.ID, &product_field.Name, product_field.Description); err != nil {
	        	fmt.Errorf("error: ", err)
	        }
	        results = append(results, product_field)
	   	}
	   	if err := rows.Err(); err != nil {
	        fmt.Errorf("error: ", err)
	    }
		err = rows.Close()
		if err != nil {
			fmt.Errorf("error: ", err)
		}
		ipp := max(1, signals.ItemsPerPage)
		pn := max(1, signals.PageNumber)
		page := results[ipp*(min(pn - 1, len(results)/ipp)):min(pn*ipp, len(results))]
		t, err := template.New("results").Parse(`<tbody id="current-table">{{range .pages}}<tr id="row-{{.ID}}"><td data-signals="{ID: {{.ID}} }">{{.ID}}</td><td data-bind:product-name>{{.Name}}</td><td><button data-on:click="confirm('Are you sure?') && @delete('/customers?id={{.ID}}')">Delete</button></td><td><button data-on:click="@get('/customers/edit?id={{.ID}}&name={{.Name}}')">Edit</button></td></tr>{{end}}</tbody>`)
		if err != nil {
      		fmt.Errorf("error: ", err)
	    }
	    var builder strings.Builder
	    res := map[string][]Product{"pages": page}
	    t.Execute(&builder, res)
	    searchResult := builder.String()
	    sse := datastar.NewSSE(w, r)
		sse.PatchElements(searchResult)
	}
}

func employees (w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			t := template.Must(template.ParseFiles("templates/base.html", "templates/employees.html"))
	}
}

func orders (w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			t := template.Must(template.ParseFiles("templates/base.html", "templates/orders.html"))
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
	    // var err error
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
	http.HandleFunc("/customers/edit", editCustomers)

	http.HandleFunc("/products", products)
	http.HandleFunc("/products/search", products)
	http.HandleFunc("/products/edit", products)

	http.HandleFunc("/orders", orders)
	http.HandleFunc("/employees", employees)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
