CREATE TABLE `activity_logs` (
  `id` varchar(100) PRIMARY KEY,
  `user_id` varchar(100) NOT NULL,
  `source_uri` varchar(255) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  INDEX (`user_id`, `source_uri`)
);
