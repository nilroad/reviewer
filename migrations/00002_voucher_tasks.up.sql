-- Create voucher_task_types table
CREATE TABLE voucher_task_types (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name enum('VOUCHER_CSV_ASSIGN', 'VOUCHER_BATCH_ASSIGN') NOT NULL,
    is_active        tinyint(1)  not null default 1,
    sample_file_path text NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    
    UNIQUE KEY idx_voucher_task_types_name (name),
    KEY idx_voucher_task_types_is_active (is_active)
);

-- Create voucher_tasks table
CREATE TABLE voucher_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    voucher_task_type_id BIGINT UNSIGNED NOT NULL,
    voucher_rule_type_number varchar(255) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    file_path text NOT NULL,
    items_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
    status enum('PENDING', 'PROCESSING', 'DONE', 'CANCELED', 'FAILED') NOT NULL,
    started_at timestamp not null,
    finished_at timestamp null,
    created_at timestamp not null,
    updated_at timestamp not null,
    
    KEY idx_voucher_tasks_voucher_task_type_id (voucher_task_type_id),
    KEY idx_voucher_tasks_user_id (user_id),
    KEY idx_voucher_tasks_status (status),
    KEY idx_voucher_tasks_started_at (started_at),
    KEY idx_voucher_tasks_finished_at (finished_at),
    
    CONSTRAINT fk_voucher_tasks_voucher_task_type_id
        FOREIGN KEY (voucher_task_type_id) 
        REFERENCES voucher_task_types(id)
);

-- Create voucher_task_items table  
CREATE TABLE voucher_task_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_id BIGINT UNSIGNED NOT NULL,
    payload TEXT NOT NULL,
    status enum('PENDING', 'PROCESSING', 'PROCESSED', 'CANCELED', 'FAILED') NOT NULL,
    error_reason TEXT NULL,
    created_at timestamp not null,
    updated_at timestamp not null,
    
    KEY idx_voucher_task_items_task_id (task_id),
    KEY idx_voucher_task_items_status (status),
    KEY idx_voucher_task_items_created_at (created_at),
    
    CONSTRAINT fk_voucher_task_items_task_id
        FOREIGN KEY (task_id) 
        REFERENCES voucher_tasks(id)
);
