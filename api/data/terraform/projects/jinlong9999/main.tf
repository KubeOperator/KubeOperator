provider "openstack" {
  user_name   = "admin"
  tenant_id   = "ed2838ecd90a4ec5a1ef5cf305bef59c"
  password    = "Calong@2015"
  auth_url    = "http://openstack.fit2cloud.com/identity/v3"
  region      = "RegionOne"
  user_domain_name = "Default"
}




    resource "openstack_networking_port_v2" "worker1" {
      name           = "worker1"
      admin_state_up = "true"
      network_id = "be6d9615-b273-4f4c-ba27-e65e7e66aa63"
    }

    resource "openstack_compute_instance_v2" "worker1" {
      name            = "worker1"
      image_name        = "kubeoperator_centos_7.6.1810"
      
       flavor_name          = "m1.medium"
      
      
        

      security_groups = ["default"]
      network {
        port = "${openstack_networking_port_v2.worker1.id}"
      }
    }

    resource "openstack_networking_floatingip_v2" "worker1" {
      pool = "658a6c61-957c-4294-ad40-bc7fedaee665"
      address = "172.18.22.250"
    }

    resource "openstack_compute_floatingip_associate_v2" "worker1" {
      floating_ip = "${openstack_networking_floatingip_v2.worker1.address}"
      instance_id = "${openstack_compute_instance_v2.worker1.id}"
      fixed_ip    = "${openstack_compute_instance_v2.worker1.network.0.fixed_ip_v4}"
    }




    resource "openstack_networking_port_v2" "daemon1" {
      name           = "daemon1"
      admin_state_up = "true"
      network_id = "be6d9615-b273-4f4c-ba27-e65e7e66aa63"
    }

    resource "openstack_compute_instance_v2" "daemon1" {
      name            = "daemon1"
      image_name        = "kubeoperator_centos_7.6.1810"
      
      
        
       flavor_name          = "m1.medium"
      

      security_groups = ["default"]
      network {
        port = "${openstack_networking_port_v2.daemon1.id}"
      }
    }

    resource "openstack_networking_floatingip_v2" "daemon1" {
      pool = "658a6c61-957c-4294-ad40-bc7fedaee665"
      address = "172.18.22.251"
    }

    resource "openstack_compute_floatingip_associate_v2" "daemon1" {
      floating_ip = "${openstack_networking_floatingip_v2.daemon1.address}"
      instance_id = "${openstack_compute_instance_v2.daemon1.id}"
      fixed_ip    = "${openstack_compute_instance_v2.daemon1.network.0.fixed_ip_v4}"
    }




    resource "openstack_networking_port_v2" "worker3" {
      name           = "worker3"
      admin_state_up = "true"
      network_id = "be6d9615-b273-4f4c-ba27-e65e7e66aa63"
    }

    resource "openstack_compute_instance_v2" "worker3" {
      name            = "worker3"
      image_name        = "kubeoperator_centos_7.6.1810"
      
       flavor_name          = "m1.medium"
      
      
        

      security_groups = ["default"]
      network {
        port = "${openstack_networking_port_v2.worker3.id}"
      }
    }

    resource "openstack_networking_floatingip_v2" "worker3" {
      pool = "658a6c61-957c-4294-ad40-bc7fedaee665"
      address = "172.18.22.248"
    }

    resource "openstack_compute_floatingip_associate_v2" "worker3" {
      floating_ip = "${openstack_networking_floatingip_v2.worker3.address}"
      instance_id = "${openstack_compute_instance_v2.worker3.id}"
      fixed_ip    = "${openstack_compute_instance_v2.worker3.network.0.fixed_ip_v4}"
    }




    resource "openstack_networking_port_v2" "worker2" {
      name           = "worker2"
      admin_state_up = "true"
      network_id = "be6d9615-b273-4f4c-ba27-e65e7e66aa63"
    }

    resource "openstack_compute_instance_v2" "worker2" {
      name            = "worker2"
      image_name        = "kubeoperator_centos_7.6.1810"
      
       flavor_name          = "m1.medium"
      
      
        

      security_groups = ["default"]
      network {
        port = "${openstack_networking_port_v2.worker2.id}"
      }
    }

    resource "openstack_networking_floatingip_v2" "worker2" {
      pool = "658a6c61-957c-4294-ad40-bc7fedaee665"
      address = "172.18.22.249"
    }

    resource "openstack_compute_floatingip_associate_v2" "worker2" {
      floating_ip = "${openstack_networking_floatingip_v2.worker2.address}"
      instance_id = "${openstack_compute_instance_v2.worker2.id}"
      fixed_ip    = "${openstack_compute_instance_v2.worker2.network.0.fixed_ip_v4}"
    }




    resource "openstack_networking_port_v2" "master1" {
      name           = "master1"
      admin_state_up = "true"
      network_id = "be6d9615-b273-4f4c-ba27-e65e7e66aa63"
    }

    resource "openstack_compute_instance_v2" "master1" {
      name            = "master1"
      image_name        = "kubeoperator_centos_7.6.1810"
      
      
       flavor_name          = "m1.medium"
      
        

      security_groups = ["default"]
      network {
        port = "${openstack_networking_port_v2.master1.id}"
      }
    }

    resource "openstack_networking_floatingip_v2" "master1" {
      pool = "658a6c61-957c-4294-ad40-bc7fedaee665"
      address = "172.18.22.252"
    }

    resource "openstack_compute_floatingip_associate_v2" "master1" {
      floating_ip = "${openstack_networking_floatingip_v2.master1.address}"
      instance_id = "${openstack_compute_instance_v2.master1.id}"
      fixed_ip    = "${openstack_compute_instance_v2.master1.network.0.fixed_ip_v4}"
    }


