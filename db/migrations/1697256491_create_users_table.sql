CREATE TABLE `users` (
  `id` varchar(100) PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `external_id` varchar(100) NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `unique_user_identifier` (`email`,`external_id`),
  INDEX (`email`, `external_id`)
);
