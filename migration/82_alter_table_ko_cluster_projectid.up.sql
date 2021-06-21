ALTER TABLE `ko`.`ko_cluster` ADD COLUMN `project_id` varchar(255) NULL AFTER `source`;

UPDATE ko_cluster 
JOIN ko_project_resource ON ko_project_resource.resource_id = ko_cluster.id AND ko_project_resource.resource_type = 'CLUSTER' 
JOIN ko_project ON ko_project.ID = ko_project_resource.project_id 
SET ko_cluster.project_id = ko_project.ID;