-- Create "templates" table
CREATE TABLE `templates` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `name` varchar(255) NOT NULL,
 `description` longtext NULL,
 `file_path` varchar(255) NOT NULL,
 `enabled` bool NOT NULL DEFAULT 1,
 PRIMARY KEY (`id`),
 INDEX `idx_templates_deleted_at` (`deleted_at`),
 UNIQUE INDEX `idx_templates_file_path` (`file_path`),
 UNIQUE INDEX `idx_templates_name` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "tenants" table
CREATE TABLE `tenants` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `name` varchar(100) NOT NULL,
 `description` varchar(255) NOT NULL,
 `is_active` bool NULL DEFAULT 1,
 `is_billing_enabled` bool NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 INDEX `idx_tenants_deleted_at` (`deleted_at`),
 UNIQUE INDEX `uni_tenants_name` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "users" table
CREATE TABLE `users` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `keycloak_id` char(36) NULL,
 `role` longtext NULL,
 PRIMARY KEY (`id`),
 INDEX `idx_users_deleted_at` (`deleted_at`),
 UNIQUE INDEX `idx_users_keycloak_id` (`keycloak_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "sites" table
CREATE TABLE `sites` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `name` varchar(255) NOT NULL,
 `description` longtext NULL,
 `domain` varchar(255) NOT NULL,
 `enabled` bool NOT NULL DEFAULT 1,
 `template_id` bigint unsigned NOT NULL,
 PRIMARY KEY (`id`),
 UNIQUE INDEX `idx_site_template` (`template_id`),
 INDEX `idx_sites_deleted_at` (`deleted_at`),
 UNIQUE INDEX `idx_sites_domain` (`domain`),
 UNIQUE INDEX `idx_sites_name` (`name`),
 INDEX `idx_sites_template_id` (`template_id`),
 CONSTRAINT `fk_sites_template` FOREIGN KEY (`template_id`) REFERENCES `templates` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "template_settings" table
CREATE TABLE `template_settings` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `template_id` bigint unsigned NOT NULL,
 `setting_key` varchar(100) NOT NULL,
 `setting_value` longtext NOT NULL,
 `can_override` bool NOT NULL DEFAULT 0,
 PRIMARY KEY (`id`),
 UNIQUE INDEX `idx_template_setting` (`template_id`, `setting_key`),
 INDEX `idx_template_settings_deleted_at` (`deleted_at`),
 INDEX `idx_template_settings_setting_key` (`setting_key`),
 INDEX `idx_template_settings_template_id` (`template_id`),
 CONSTRAINT `fk_templates_settings` FOREIGN KEY (`template_id`) REFERENCES `templates` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "template_setting_overrides" table
CREATE TABLE `template_setting_overrides` (
 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
 `created_at` datetime(3) NULL,
 `updated_at` datetime(3) NULL,
 `deleted_at` datetime(3) NULL,
 `site_id` bigint unsigned NOT NULL,
 `template_setting_id` bigint unsigned NOT NULL,
 `setting_value` longtext NOT NULL,
 PRIMARY KEY (`id`),
 UNIQUE INDEX `idx_template_setting_override` (`site_id`, `template_setting_id`),
 INDEX `idx_template_setting_overrides_deleted_at` (`deleted_at`),
 INDEX `idx_template_setting_overrides_site_id` (`site_id`),
 INDEX `idx_template_setting_overrides_template_setting_id` (`template_setting_id`),
 CONSTRAINT `fk_sites_setting_overrides` FOREIGN KEY (`site_id`) REFERENCES `sites` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
 CONSTRAINT `fk_template_setting_overrides_template_setting` FOREIGN KEY (`template_setting_id`) REFERENCES `template_settings` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
