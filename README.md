# ecommerce-tracker
A business tool to track data related to customers, orders, and finances. The program has six main sections.
## Installstion
Download the binary from the releases page.
Ecommerce Tracker uses MySQL. Make sure you have a server and run the `schema.sql` script to create the database. 
Modify `sample-config.ini` to match your MySQL address and login info, and save it as `config.ini`.
By default the web UI can be accessed at `localhost:8080`.
## Customers
View details about your customers. Their LTV is based on their order volume from the orders page. New customers do not appear in the table until after order data has been added for them.
## Products
View your product information and inventory costing. The lead time and reorder point take order data into account. The interest rate used to calculate holdig cost for the ecomomic order quantity can be changed in `config.ini`. The unit price can be calculated using FIFO or  
## Orders
When a new order is added, a journal entry is created. The inventory costing table is updated once the progress is `delivered`.
## Employees
The salary is yearly and added to the accrued expenses account.
## Reviews
View what customers say about the their orders as well as their rating.
## Finances
This page has the chart of accounts and journal entries. 

You can add journal entries by entering the codes and the amount for the accounts being debited and credited.

The journal entries are used to calculate the balances for the general ledger and trial balance. The income and balance sheet financial statements are based on the adjusted trial balance. The cash flow statment is focused on operating cash flow, and calculated from the income statement and balance sheet.
