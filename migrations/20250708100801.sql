-- Modify "pages" table
ALTER TABLE `pages` ADD UNIQUE INDEX `unique_page_key` (`key`, `path`);
-- Modify "tenants" table
ALTER TABLE `tenants` RENAME INDEX `uni_tenants_name` TO `idx_tenants_name`;
