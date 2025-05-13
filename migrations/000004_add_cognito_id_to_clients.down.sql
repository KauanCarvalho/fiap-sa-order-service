ALTER TABLE `clients`
DROP INDEX `unique_cognito_id`,
DROP COLUMN `cognito_id`;
