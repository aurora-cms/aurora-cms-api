-- Create "user_tenants" table
CREATE TABLE `user_tenants` (
 `tenant_id` bigint unsigned NOT NULL,
 `user_id` bigint unsigned NOT NULL,
 PRIMARY KEY (`tenant_id`, `user_id`),
 INDEX `fk_user_tenants_user` (`user_id`),
 CONSTRAINT `fk_user_tenants_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
 CONSTRAINT `fk_user_tenants_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
