-- Modify "users" table
ALTER TABLE `users` MODIFY COLUMN `role` varchar(255) NOT NULL DEFAULT "user";
-- Modify "sites" table
ALTER TABLE `sites` DROP FOREIGN KEY `fk_sites_tenant`, ADD CONSTRAINT `fk_tenants_sites` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "pages" table
CREATE TABLE `pages` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `key` varchar(255) NOT NULL,
 `path` varchar(255) NULL,
 `index` bigint NOT NULL DEFAULT 0,
 `parent_id` bigint unsigned NULL,
 `site_id` bigint unsigned NOT NULL,
 `type` varchar(32) NOT NULL DEFAULT "content",
 `link_url` varchar(255) NULL,
 `hard_link_page_id` bigint unsigned NULL,
 PRIMARY KEY (`id`),
 INDEX `idx_pages_deleted_at` (`deleted_at`),
 INDEX `idx_pages_hard_link_page_id` (`hard_link_page_id`),
 INDEX `idx_pages_parent_id` (`parent_id`),
 INDEX `idx_pages_site_id` (`site_id`),
 CONSTRAINT `fk_pages_children` FOREIGN KEY (`parent_id`) REFERENCES `pages` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
 CONSTRAINT `fk_sites_pages` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "page_versions" table
CREATE TABLE `page_versions` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `page_id` bigint unsigned NOT NULL,
 `version` bigint unsigned NOT NULL DEFAULT 1,
 `title` varchar(255) NOT NULL,
 `description` varchar(255) NULL,
 `is_published` bool NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 INDEX `idx_page_versions_deleted_at` (`deleted_at`),
 INDEX `idx_page_versions_page_id` (`page_id`),
 CONSTRAINT `fk_pages_versions` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "page_blocks" table
CREATE TABLE `page_blocks` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `page_version_id` bigint unsigned NOT NULL,
 `block_key` varchar(255) NOT NULL,
 `index` bigint NOT NULL DEFAULT 0,
 `content_type` varchar(255) NOT NULL,
 `content` longtext NOT NULL,
 PRIMARY KEY (`id`),
 INDEX `idx_page_blocks_deleted_at` (`deleted_at`),
 INDEX `idx_page_blocks_page_version_id` (`page_version_id`),
 CONSTRAINT `fk_page_versions_blocks` FOREIGN KEY (`page_version_id`) REFERENCES `page_versions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
