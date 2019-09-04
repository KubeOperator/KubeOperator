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


data "vsphere_resource_pool" "az-2" {
   name          = "az-2"
   datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "az-2" {
  name = "VM Network"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "az-2" {
  name = "vsanDatastore"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_resource_pool" "az-3" {
   name          = "az-3"
   datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "az-3" {
  name = "VM Network"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "az-3" {
  name = "vsanDatastore"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_resource_pool" "az-1" {
   name          = "az-1"
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
  name = "kubeoperator_centos_7.6.1810"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}


resource "vsphere_virtual_machine" "master1" {
  name = "master1.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-3.id}"
  datastore_id = "${data.vsphere_datastore.az-3.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-3.id}"
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
        host_name = "master1"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.98"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "daemon1" {
  name = "daemon1.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-3.id}"
  datastore_id = "${data.vsphere_datastore.az-3.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-3.id}"
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
        host_name = "daemon1"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.97"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "master2" {
  name = "master2.cluster.f2c.com"
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
        host_name = "master2"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.73"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "master3" {
  name = "master3.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-2.id}"
  datastore_id = "${data.vsphere_datastore.az-2.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-2.id}"
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
        host_name = "master3"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.87"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "worker2" {
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
        host_name = "worker2"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.72"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "worker1" {
  name = "worker1.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-3.id}"
  datastore_id = "${data.vsphere_datastore.az-3.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-3.id}"
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
        host_name = "worker1"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.96"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "worker3" {
  name = "worker3.cluster.f2c.com"
  folder = "kubeoperator"
  resource_pool_id = "${data.vsphere_resource_pool.az-2.id}"
  datastore_id = "${data.vsphere_datastore.az-2.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos6_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.az-2.id}"
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
        host_name = "worker3"
        domain = "cluster.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.86"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}
