ALTER TABLE `clients`
ADD COLUMN `cognito_id` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
ADD UNIQUE KEY `unique_cognito_id` (`cognito_id`);
