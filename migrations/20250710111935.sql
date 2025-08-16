-- Modify "sites" table
ALTER TABLE `sites` ADD COLUMN `title_template` varchar(255) NULL AFTER `domain`;
