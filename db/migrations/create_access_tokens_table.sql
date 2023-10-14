CREATE TABLE `access_tokens` (
  `id` varchar(100) PRIMARY KEY,
  `user_id` varchar(100) NOT NULL,
  `token_string` varchar(255) NOT NULL,
  `expired_at` datetime,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  INDEX (`token_string`)
);
