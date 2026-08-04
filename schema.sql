DROP TABLE IF EXISTS `chart_of_accounts`;
CREATE TABLE `chart_of_accounts` (
  `account_code` int DEFAULT NULL,
  `account_name` text,
  `account_type` enum('Asset','Liability','Revenue','Expense','Equity') DEFAULT NULL,
  `financial_statement` enum('Balance Sheet','Income Statement','Cash Flow Statement') DEFAULT NULL
  );
INSERT INTO `chart_of_accounts` (`account_code`, `account_name`, `account_type`, `financial_statement`) VALUES
(2301, 'Sales', 'Revenue', 'Income Statement'),
(2401, 'Cost of Goods Sold', 'Expense', 'Income Statement'),
(2402, 'Promotion Expenses', 'Expense', 'Income Statement'),
(2403, 'Selling and General Administrative Expenses', 'Expense', 'Income Statement'),
(2404, 'Depreciation and Amoritization', 'Expense', 'Income Statement'),
(1101, 'Cash', 'Asset', 'Balance Sheet'),
(1102, 'Accounts Receivable', 'Asset', 'Balance Sheet'),
(1105, 'Prepaid Expenses', 'Asset', 'Balance Sheet'),
(1103, 'Property and Equipent', 'Asset', 'Balance Sheet'),
(1201, 'Accounts Payable', 'Liability', 'Balance Sheet'),
(1203, 'Accurued Expenses', 'Liability', 'Balance Sheet'),
(1204, 'Unearned Revenue', 'Liability', 'Balance Sheet'),
(1202, 'Long Term Debt', 'Liability', 'Balance Sheet'),
(1104, 'Inventory', 'Asset', 'Balance Sheet'),
(2405, 'Wages', 'Expense', 'Income Statement'),
(1501, 'Equity', 'Equity', 'Balance Sheet');
DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `customer_name` text,
  `subscriber` tinyint(1) DEFAULT '0',
  `cohort` enum('engaged prospect','lapsed prospect','one purchase','two or more purchases','VIP','churned') DEFAULT 'engaged prospect',
  `signup_date` date DEFAULT '2026-01-01',
  PRIMARY KEY (`id`)
);
DROP TABLE IF EXISTS `employees`;
CREATE TABLE `employees` (
  `id` int NOT NULL AUTO_INCREMENT,
  `employee_name` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `job_title` text,
  `department` text,
  `start_date` date DEFAULT NULL,
  `phone` text,
  `company_email` text,
  `salary` int DEFAULT NULL,
  PRIMARY KEY (`id`)
);
DROP TABLE IF EXISTS `journal_entries`;
CREATE TABLE `journal_entries` (
  `id` int NOT NULL AUTO_INCREMENT,
  `entry_date` date DEFAULT '2026-01-01',
  `debit_account` int DEFAULT '0',
  `credit_account` int DEFAULT '0',
  `debit` int DEFAULT '0',
  `credit` int DEFAULT '0',
  `notes` text,
  PRIMARY KEY (`id`)
);
DROP TABLE IF EXISTS `costs`;
CREATE TABLE `costs` (
  `cost_name` text,
  `classification` enum('direct labor','direct material','indirect labor','indirect material','overhead') DEFAULT NULL,
  `variable_or_fixed` enum('variable cost','fixed cost') DEFAULT NULL,
  `amount` int DEFAULT NULL,
  `output_range` text
);
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `product_name` text,
  `description` text,
  `price` int DEFAULT '0',
  `minimum_order_quantity` int DEFAULT '0',
  PRIMARY KEY (`id`)
);
DROP TABLE IF EXISTS `inventory_costing`;
CREATE TABLE `inventory_costing` (
  `product_id` int DEFAULT NULL,
  `time_period` date DEFAULT NULL,
  `inventory_type` enum('Purchase','Cost of Goods Sold','Inventory') DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  `unit_cost` int DEFAULT NULL,
  `total` int DEFAULT NULL,
  KEY `product_id` (`product_id`),
  FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
);
CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `customer_id` int NOT NULL,
  `product_id` int NOT NULL,
  `quantity` int NOT NULL,
  `total` int NOT NULL DEFAULT '0',
  `progress` enum('received','in transport','delivered') NOT NULL,
  `order_date` date DEFAULT NULL,
  `cost_of_goods_sold` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `customer_id` (`customer_id`),
  KEY `product_id` (`product_id`),
  FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`),
  FOREIGN KEY (`product_id`) REFERENCES `products` (`id`)
);
DROP TABLE IF EXISTS `reviews`;
CREATE TABLE `reviews` (
  `order_id` int DEFAULT NULL,
  `rating` int DEFAULT NULL,
  `review` text,
  KEY `order_id` (`order_id`),
  FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`)
)
