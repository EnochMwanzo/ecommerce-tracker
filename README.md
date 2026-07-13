# ecommerce-tracker
A business tool to track data related to customers, orders, and finances. The program has six main sections.
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
<img src="/screenshots/customers.png" style="width:75%;" />
</p>
View details about your customers. Their Lifetime Value is based on the sum of totals in the orders that match their customer ID on the orders page. New customers do not appear in the marketing analysis table until after order data has been added for them.

### Products
<p align="center">
<img src="/screenshots/products.png" style="width:75%;">
</p>
View your product information and inventory costing. The program currently uses weighted average to calculate cost of goods sold.

### Orders
<p align="center">
<img src="/screenshots/orders.png" style="width:75%;">
</p>
When a new order is added, a journal entry is created. The inventory costing table is updated once the progress is `delivered`.

### Employees
<p align="center">
<img src="/screenshots/employees.png" style="width:75%;">
</p>
You can store information about employees here.

### Reviews
<p align="center">
<img src="/screenshots/reviews.png" style="width:75%;">
</p>
Save what customer rating and what they say about the their orders here.

### Finances
<p align="center">
<img src="/screenshots/chart-of-accounts.png" style="width:45%; border: 5px solid black;"> <img src="/screenshots/journal-entries.png" style="width:45%;">
</p>
This page has the chart of accounts, general ledger, the financial statements, and journal entries, and trial balance.

You can add journal entries by entering the codes and the amount for the accounts being debited and credited.

The general ledger is created by going through all the journal entries and keeping track of debits and credits for each account.

The journal entries are used to calculate the balances for the general ledger and trial balance. The income and balance sheet financial statements are based on the adjusted trial balance. The cash flow statment is focused on operating cash flow, and calculated from the income statement and balance sheet.

You can generate financial statements to prit or save as PDF 
<div style="display: grid; grid-template-columns: 1fr 1fr;">
<div align="center" style="width:45%;">
  <img src="/screenshots/cash-flow-statement.png" style="width:30%;"> <img src="/screenshots/income-statement.png" style="width:30%;"><img src="/screenshots/balance-sheet.png" style="width:30%;">
<div align="center" style="width:45%;"> 
  <img src="/screenshots/general-ledger.png" style="width:30%;"></div></div>
</div>
