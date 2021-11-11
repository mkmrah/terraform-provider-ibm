# cluster config file path
output "cluster_config_file_path" {
  value = data.ibm_container_cluster_config.cluster_config.config_file_path
}


output cluster_id {
    description = "Cluster ID, will await multizone attachment before returning"
    value       = ibm_container_vpc_cluster.cluster.id
}

