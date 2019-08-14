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

data "vsphere_compute_cluster" "cluster" {
  name = "vSAN-Cluster"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "network" {
  name = "VM Network"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "datastore" {
  name = "vsanDatastore"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}
data "vsphere_virtual_machine" "template" {
  name = "Centos7.6_vSan_template"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}


resource "vsphere_virtual_machine" "master1" {
  name = "master1.sdfsdf.f2c.com"
  folder = "kubeops5"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label = "disk0"
    size = 50
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = "master1"
        domain = "sdfsdf.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.240"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "daemon1" {
  name = "daemon1.sdfsdf.f2c.com"
  folder = "kubeops5"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label = "disk0"
    size = 50
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = "daemon1"
        domain = "sdfsdf.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.241"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "worker1" {
  name = "worker1.sdfsdf.f2c.com"
  folder = "kubeops5"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label = "disk0"
    size = 50
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = "worker1"
        domain = "sdfsdf.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.242"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "worker2" {
  name = "worker2.sdfsdf.f2c.com"
  folder = "kubeops5"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label = "disk0"
    size = 50
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = "worker2"
        domain = "sdfsdf.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.243"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}

resource "vsphere_virtual_machine" "worker3" {
  name = "worker3.sdfsdf.f2c.com"
  folder = "kubeops5"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 2
  memory = 8192
  guest_id = "centos7_64Guest"

  network_interface {
    network_id = "${data.vsphere_network.network.id}"
  }

  disk {
    label = "disk0"
    size = 50
  }

  clone {
    template_uuid = "${data.vsphere_virtual_machine.template.id}"
    customize {

      linux_options {
        host_name = "worker3"
        domain = "sdfsdf.f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.244"
        ipv4_netmask = 24
      }
      ipv4_gateway = "172.16.10.254"
    }
  }
}
