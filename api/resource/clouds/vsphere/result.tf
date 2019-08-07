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
  name = "vSan-Cluster"
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


resource "vsphere_virtual_machine" "vm" {
  name = "master1.cluster"
  folder = "kubeops"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 4
  memory = 8096
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
        host_name = "master1.cluster"
        domain = "f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.201"
        ipv4_netmask = 24
      }
    }
  }
}

resource "vsphere_virtual_machine" "vm" {
  name = "daemon1.cluster"
  folder = "kubeops"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 4
  memory = 8096
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
        host_name = "daemon1.cluster"
        domain = "f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.202"
        ipv4_netmask = 24
      }
    }
  }
}

resource "vsphere_virtual_machine" "vm" {
  name = "worker1.cluster"
  folder = "kubeops"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 4
  memory = 8096
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
        host_name = "worker1.cluster"
        domain = "f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.203"
        ipv4_netmask = 24
      }
    }
  }
}

resource "vsphere_virtual_machine" "vm" {
  name = "worker2.cluster"
  folder = "kubeops"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 4
  memory = 8096
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
        host_name = "worker2.cluster"
        domain = "f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.204"
        ipv4_netmask = 24
      }
    }
  }
}

resource "vsphere_virtual_machine" "vm" {
  name = "worker3.cluster"
  folder = "kubeops"
  resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
  datastore_id = "${data.vsphere_datastore.datastore.id}"
  num_cpus = 4
  memory = 8096
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
        host_name = "worker3.cluster"
        domain = "f2c.com"
      }

      network_interface {
        ipv4_address = "172.16.10.205"
        ipv4_netmask = 24
      }
    }
  }
}
