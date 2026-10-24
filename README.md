# ecommerce-tracker
A simple business tool to track data related to customers, inventory, and finances. 

It is a local web server designed to be simple, lightweight, and teach myself more about accounting.
## Installstion
1. Download and extract the program and supporting files from the <a href="https://github.com/EnochMwanzo/ecommerce-tracker/releases/tag/v.0.1.0-alpha">releases page</a>.
2. Ecommerce Tracker uses MySQL. <a href="https://dev.mysql.com/downloads/mysql/">You can download it here.</a> There are post-installation steps you may need to complete depending on whether you are on <a href="https://dev.mysql.com/doc/refman/9.7/en/windows-postinstallation.html">Windows</a> or <a href="https://dev.mysql.com/doc/refman/9.7/en/macos-installation-launchd.html">macOS</a>.
3. Open a terminal and log in to the server with `mysql -u root -p`. You may want to <a href="https://www.digitalocean.com/community/tutorials/how-to-create-a-new-user-and-grant-permissions-in-mysql"> create a non root user</a>.
4. At the `mysql` command line create a new database
```sql
CREATE DATABASE ecommerce_tracker
```
5. Then run `source /path/to/schema.sql` to to create the tables the program is expecting.

(You may also be able to complete 2-4 using <a href="https://dev.mysql.com/doc/workbench/en/wb-intro.html">MySQL Workbench</a>)

6.  Modify `sample-config.ini` to match your MySQL address and login info, and save it as `config.ini`. By default the web UI can be accessed at `localhost:8080`.
## Sections
The program has six main sections.
### Customers
<p align="center">
<img src="/screenshots/Customers.png" style="width:75%;" />
</p>
View details about customers, as well as RFM (Recency Frequency Monetary Value) analysis.

* ID: every customer has a unique ID.
* Customer Name: The name of the customer, can e first name only, or first and last name.
* Cohort: A way to categorize customers. Options are `'engaged prospect','lapsed prospect','one purchase','two or more purchases','VIP','churned'`.
* Subscriber: Whether or not the customer is subscried to a service provided by the company. 
* Signup Date: The day the customer was entered into the database.
* Recency (Days): How any days have passed since the customer placed an order. It is calculated by taking the difference etween the unix timestamps of the date or the most recent purchase and the current date.
* Order Frequency Days: The average number of days between the customer's orders.
* Lifetime Value: Lifetime Value is based on the sum of totals in the orders that match their customer ID on the orders page. New customers do not appear in the marketing analysis table until after order data has been added for them.
* Number Of Orders: The sum of orders that the customer has placed. Each order is only one product, so if a customer was getting multiple products at once, they would be counted as multiple separate orders.
* Net Profit: The difference between the lifetime value of the customer and the su of the cost of goods sold of their orders.
* Score: A number to compare customers to understand which ones are most profitable to the business.

The 'Add a customer' form allow a new customer to be created by entering a name, subscriber status, cohort. The sign up date entered automatically. Other fields are derived from the customer's order data.

### Products
<p align="center">
<img src="/screenshots/Products.png" style="width:75%;">
</p>
View company product information and each piece of inventory. 

* Id: Each product has a unique ID
* Product Name: The name of the product.
* Description: Description of the product that may be used for the product description page on an ecommerce platform.
* Stock: The quantity of the product available for sale.
* Minimum Order Quantity: For products that are purchased from a supplier, this is the minimum order size that can be placed.
* Retail Price: The price which consumers pay for the product.
* Source: Whether the product is manufactured by the copany or purchased from another supplierr.

In the Add a Product form, a new product can be added by entering the product's name, description, retail pricec, minium order quantity, and source. The stock is updated when more product inventory is purchased.

#### Inventory
* Id: Each item has a unique ID
* Product ID: The product the inventory item is.
* Employee ID: The employee who made the product if manufactured, else 0 if it was purchased.
* Date Created: The date the inventory was manufactured or purchased.
* Order ID: The order that the item was assigned to, 0 if no order yet.
* Location: The location of the ite. Options are `'warehouse','manufacturing facility','order delivery'`
* Damaged: Whether or not the inventory is damaged. Items can be marked as damaged and the will not be assigned to an order.
* Sold: Whether or not the ite was sold. After an it is assigned to an order it s marked as sold so it cannot be assgned to an order again.
* Cost: The cost of purchasing or manufacturing the item. If the product is purchased, this is the price that was paid for this individual item. If manufactured, this is the varible cost + a fraction of overhead.
* Turnover (Days): The difference in days between creation and order.
* Notes: Any additioal notes about the inventory.

The purchase inventory form is for updating the inventory cout of products. The for takes the product, the location, the cost, and notes.
### Orders
<p align="center">
<img src="/screenshots/Orders.png" style="width:75%;">
</p>
When a new order is added, a two journal entries are created to account for the sale and inventory reduction. The cost of goods sold column of the inventory costing table is updated as well. You can edit the progress column to update the current status of orders.

* Order ID: Every order has a unique ID. If a customer buys the same product multiple ties, each order of that product will be considered as a unique order.
* Customer ID: The ID of the customer making the order. The ID matches the ID in the customers table.
* Product ID: The ID of the product being ordered. It corresponds to the ID in the products table.
* Quantity: The number of the product the customer orders.
* Total: The the price of the product being ordered multiplied by the quantity.
* Progress: The stage that the order is in. Options are `'received','in transport','delivered'`
* Order Date: The date the order was entered into the database.
* Cost of Goods Sold: The cost of the items that are assigned to the customer's order.
* Order Profit Margin: The difference between the total and the costs of goods sold. 
### Employees
<p align="center">
<img src="/screenshots/Employees.png" style="width:75%;">
</p>
You can store information about employees here.

* ID: Every employee has a unique ID. It is used to identify the employee that manufactured each inventory item and the employees listed in payroll table.
* Employee Name: The name of the employee. It can be First name and last name, or other format.
* Job Title: The name of the role the employee is assigned.
* Department: The department the employee is part of.
* Start Date: The date the employee was added to the database.
* Phone: The phone number of the employee. Any text format can be used.
* Email: The company email address of the employee.
* Salary: The yearly salary of the employee in dollars.
* Withholding Percentage: If the employee witholds for taxes, this is the percentage that will be withheld each pay period.

The Add an Employee form allows a new employee to be added to the database by entering the employee name, job title, department, phone, email salary, and their start date.
#### Payroll
The pay period which is displayed can be selected in the dropdown menu.

* Employee: The employee the pay is for.
* Jobs Completed: The number of hours the employee worked during the pay period.
* Gross Wages: The employee's hourly wage multiplied by the hours completed.
* Deductions: The sum of all benefits and other deductions for the employee.
* Withholding Amount: The amount that is withheld from the employee for taxes.
* Net Pay: What the employee will be paid on the pay date.

### Reviews
<p align="center">
<img src="/screenshots/Reviews.png" style="width:75%;">
</p>
Save customer ratings and what they say about the their orders here.

* Order ID: This matches the Order ID in the order tale that the review is for.
* Rating: A number from 1 to 5 that represents the customer rates their order.
* Review: Detailed text from the customer about their order

The add a review form allow a new review by entering the order the review is for, the rating, and review.
### Finances
<p align="center">
<img src="/screenshots/Finances.png" style="width:75%; border: 5px solid black;"> 
</p>
This page has the chart of accounts, general ledger, the financial statements, and journal entries, and trial balance.

#### Chart of Accounts
* Account Code: The 4 digit code of the account. By default the first digit in the code represents the account type, and the second digit represents the financial statement. The last two digits are unique ID for accounts of the same statement and type.
* Account Name: The name of the account.
* Account Type: Whether the account is `1. 'Asset', 2. 'Liability', 3. 'Revenue', 4. 'Expense', 5. 'Equity'`.
* Financial Statement The financial statement that uses the account. `1. 'Balance Sheet', 2. 'Income Statement', 3. 'Cash Flow Statement'`.

The edit button can be used to change the name of an account.

The 'Add an account to Chart of Accounts' form allow new accounts to be added by entering the financial statement, category, and account name.

#### General Ledger
The general ledger is not stored in the database. It is dervied from the journal entries. The account code of the account the ledger is displaying can be chaged in the dropdown.

* Journal Entry ID: The ID of the journal entry
* Credit or Debit: Whether the transaction was a credit or debit.
* Balance: What remains in the account after adding the debit or subtracting the credit.

#### Journal Entries

* Journal Entry ID: Each journal entry has a unique ID.
* Date: The data the journal entry was entered into the database.
* Debit Account: The code of the account being debited.
* Credit Account: The code of the account being credited.
* Debit: The amount that was debited.
* Credit: The amount that was cridited.
* Notes: Any additional notes about the transaction.

You can add journal entries manually in the 'Add a journal entry' form by entering the codes and the amount for the accounts being debited and credited.

#### Trial Balance

* Account Code: THe code of the account.
* Account Name: Name of the account.
* Unadjusted Balance: The account's balance in dollars before making adjusted entries.
* Adjusting Entry: The ID of the journal entry for the adjustment.
* Adjusted Balance: The balance of the accout after the adjustment is applied.

The income and balance sheet are based on the adjusted trial balance. The cash flow statment is focused on operating cash flow, and calculated from the income statement and balance sheet.

#### Income Statement

* Time Period: The first and last day that the income statement is for.
* Total Sales: The sum of sales from orders.
* Cost of Goods Sold: The price paid for purchased inventory, and the cost of goods manufactured for manufactured products.
* Profit: The difference between total sales and Cost of goods sold.
* Prootion Expenses: The sum of debits to the prootio account on the chart of accouts.
* Selling & General Administrative Expeses: The sum of debits to the Selling & General Administrative account on the chart of accouts.
* Depreciation and Amoritization: The sum of debits to the depreciation and amoritzation account on the chart of accouts.

You can generate income statements in a new tab to print or save as PDF 

#### Balance Sheet

* Time Period: The first and last day that the balance sheet is for.
* Cash: The balance of the cash account after trial balance.
* Accounts Receivable: The balance of Accounts Receivable after trial balance.
* Prepaid Expenses: The balance of the Prepaid Expenses account after trial balance.
* Inventory: The balance of the Inventory account after trial balance.
* Property And Equipment: The balance of the Property And Equipment account after trial balance.
* Accounts Payable: The balance of the Accounts Payable account after trial balance.
* Accrued Expenses; The balance of the Accrued Expenses account after trial balance.
* Unearned Revenue: The balance of the Unearned Revenue account after trial balance.
* Long Term Debt: The balance of the Long Term Debt account after trial balance.

You can generate balance sheets in a new tab to print or save as PDF 

#### Fixed Assets

* ID: Each asset has a unique ID
* Asset Name: The name of the asset.
* Starting Date: The start date for when depreciation begins.
* Starting Value: The dollar value of the asset at the beginning of depreciation.
* Depreciation Method: How the asset loses value over time, such as straight-line depreciation.
* Asset Lifetime: The amount of time until the asset depreciation is over.
* Current Value: The difference between the starting value and accumulated depreciation.

#### Finacial Ratios

* Current Ratio: The ratio of current assets to current liabilities.
* Acid Test Ratio: The acid test ratio is similar to the current ratio but only uses cash for assets.
* Gross Profit Percentage: The total from the income statement as a percentage.
### Manufacturing
<p align="center">
<img src="/screenshots/Manufacturing.png" style="width:75%; border: 5px solid black;">
</p>
The manufacturing page has tables for materials and cost of goods manufactured.

#### Bill of Materials
* Material: The name of the material.
* Product ID: The ID of the product in the products table that the material is used for.
* Supplier: The supplier that the material is purchased from.
* Unit Price: The price the supplier chargges for one unit of the material.
* Metric: The unit of measurement for the material, such as square feet or pounds.
* Unit Volume: How much material is in one unit.

The 'Add to Bill of Materials' form has fields for Material, Product, Supplier, Unit Price, and Unit Volume, and Volume for one product.

#### Materials Inventory
* Stock: The number of units of material available.
* How many units can be produced: This is the stock divided by volume for one product.
* Cash available to purchase more: True/False whether the cash account balance is enough to purchase more of the material.
#### Cost of Goods Manufactured.
* Time Period: The first and last day being displayed.
* Starting Work in Progress: The dollar value of the work in progress when the time period starts. It is the same as the Ending Work in Progress for the previous time period.
* Direct Material: The dollar value of direct materials used in the work in progress.
* Direct Labor: The dollar value of the labor needed to produce the work in progress.
* Manufacturing Overhead: The dollar value of indirect labor and materials, and other overhead costs. This is evenly divivded among the number of units in progress.
* Cost of Work in Progress: The sum of direct materials, direct labor, and manufacturing overhead for the work in progress.
* Ending Work in Progress: The dollar value of the work in progress at the end of the time period.
* Cost of Goods Manufactured: The difference between starting work in progress+ materials labor & overhead, and ending work in progress.

#### Overhead Costs
* Cost Object: What the cost is for.
* Variable or Fixed: The options are `'variable cost','fixed cost'`. Variable costs are multiplied by the number goods produced when calculating manufacturing cost.
* Amount: The dollar amount for the cost.
* Output Range: A description of what is gotten out of the cost, such as one month, or 100 uses.

Overhead costs can be added to the tale by entering the cost object, amount, output range, and whether the cost is variable of fixed.
### Master Budget
The master budget uses a projected number of units sold to estimate how money will be spent across different parts of the company throughout the year. The budgeted balance sheet and income statement are the same as ones on the fiancias page.

* Projected Units Sold: The expected number of units of each product to sell.
* Projected Sales: The projected un-its sold multiplied by the retail price of each product.
* Direct Materials: The value of the materials needed to manufacture each products.
* Direct Labor: The value of the labor needed to manufacture each product.
* Overhead: For variable overhead, it is the varible cost multiplied by the units sold. For fixed costs, they are equally divided among quarters.
* Manufactured Goods: The cost of goods manufactured for the quarter for each product that is manufactured.
* Purchased Goods: The cost of goods that were purchased for each product that is purchased.
* Salaries: The sum of the yearly salaries from employees divided among each quarter.
* Promotion: The dollar value expected to be spent on the promotion of each product. It is a percentage of sales.
#### Cash budget
* Beginning Cash Balance: The amount of cash in the cash account at the beginning of each quarter.
* Collections from Customers: Collections are the percentage of revenue that will be collected from customers during the quarter.
* Cost of Goods Sold Payments: Payments are the percentage of Manufacturing & Purchases that will be paid for during the quarter.
* Excess/Defiencency of Cash: The differece between the beginning cash balance plus collections, and selling & administrative plus payments. If the amount is negative, there is a deficiency of cash.
* Borrowings: The dollar value of money that was orrowed during the quarter.
* Repayments Including Interest: The dollar value of repayments made during the quarter, including interest.
* Ending Cash Balance: The excess/deficiency of cash, plus the difference between borrowings and repayments.
### Setings
<p align="center">
<img src="/screenshots/Settings.png" style="width:75%; border: 5px solid black;">
</p>
The settings for the program can be changed here. They match the settings in the `config.ini` file.

* User: The username used to log into the MySQL server.
* Passwd: The password to log into the MySQL server.
* Net: The networking protocol used to access the MySQL server. In the example config it is set to `tcp`.
* Addr: The IP address and port of the MySQL server.
* DBName: The name of the database to connect to. In the `schema.sql` it is called `ecommerce_tracker`, so that should be used unless the schema file or database name was modified.
* DefaultPage: The page that the program opens with when no page is specified in the web server URL.
## Roadmap
* Logistics: Determine how long it takes to move inventory between locations and what are the shipping costs.
* Compare company metrics to industry benchmarks.
* Process bank statements and scan receipts.
