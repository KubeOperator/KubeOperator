provider "vsphere" {
  user = "administrator@vsphere.local"
  password = "Calong@2015"
  vsphere_server = "172.16.10.20"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = "Datacenter"
}


data "vsphere_resource_pool" "az-1" {
  
  
   name          = "vSAN-Cluster/Resources/az-1"
  
   datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "az-1" {
  name = "VM Network"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "az-1" {
  name = "vsanDatastore"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}


data "vsphere_virtual_machine" "template" {
  name = "kubeoperator/kubeoperator_centos_7.6.1810"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}


resource "vsphere_virtual_machine" "" {
  name = "master1.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-1.id}"
  datastore_id = "${data.vsphere_datastore.az-1.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-1.id}"
  }

  disk {
    label            = "disk0"
    size             = "${data.vsphere_virtual_machine.template.disks.0.size}"
    eagerly_scrub    = "${data.vsphere_virtual_machine.template.disks.0.eagerly_scrub}"
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }


  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = ""
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.149"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "" {
  name = "daemon1.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-1.id}"
  datastore_id = "${data.vsphere_datastore.az-1.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-1.id}"
  }

  disk {
    label            = "disk0"
    size             = "${data.vsphere_virtual_machine.template.disks.0.size}"
    eagerly_scrub    = "${data.vsphere_virtual_machine.template.disks.0.eagerly_scrub}"
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }


  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = ""
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.148"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "" {
  name = "worker1.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-1.id}"
  datastore_id = "${data.vsphere_datastore.az-1.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-1.id}"
  }

  disk {
    label            = "disk0"
    size             = "${data.vsphere_virtual_machine.template.disks.0.size}"
    eagerly_scrub    = "${data.vsphere_virtual_machine.template.disks.0.eagerly_scrub}"
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }


  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = ""
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.147"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "" {
  name = "worker2.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-1.id}"
  datastore_id = "${data.vsphere_datastore.az-1.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-1.id}"
  }

  disk {
    label            = "disk0"
    size             = "${data.vsphere_virtual_machine.template.disks.0.size}"
    eagerly_scrub    = "${data.vsphere_virtual_machine.template.disks.0.eagerly_scrub}"
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }


  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = ""
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.146"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "" {
  name = "worker3.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-1.id}"
  datastore_id = "${data.vsphere_datastore.az-1.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-1.id}"
  }

  disk {
    label            = "disk0"
    size             = "${data.vsphere_virtual_machine.template.disks.0.size}"
    eagerly_scrub    = "${data.vsphere_virtual_machine.template.disks.0.eagerly_scrub}"
    thin_provisioned = "${data.vsphere_virtual_machine.template.disks.0.thin_provisioned}"
  }


  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = ""
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.145"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}
