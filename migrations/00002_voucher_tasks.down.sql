-- Drop voucher_task_items table (first, due to foreign key constraints)
DROP TABLE IF EXISTS voucher_task_items;

-- Drop voucher_tasks table (second, due to foreign key constraints)  
DROP TABLE IF EXISTS voucher_tasks;

-- Drop voucher_task_types table (last, as it's referenced by voucher_tasks)
DROP TABLE IF EXISTS voucher_task_types;
