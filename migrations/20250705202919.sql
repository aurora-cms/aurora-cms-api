-- Modify "sites" table
ALTER TABLE `sites` ADD COLUMN `tenant_id` bigint unsigned NOT NULL, ADD INDEX `idx_sites_tenant_id` (`tenant_id`), ADD CONSTRAINT `fk_sites_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON UPDATE CASCADE ON DELETE RESTRICT;
