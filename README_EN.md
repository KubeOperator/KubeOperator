# KubeOperator – Hop onto the sailing of Kubernetes

[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeOperatpr/KubeOperatpr/blob/master/LICENSE)
[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

## What is KubeOperator?

KubeOperator is an open-source light-weighted Kubernetes distribution that focuses on helping enterprises plan, deploy, and operate production-grade Kubernetes clusters in an offline network environment. It has a graphic Web UI that fasten up the process of software lifecycle in this current rapid cloud age.

## How it works?

KubeOperator uses Terraform to auto-build infrastructure on LaaS platform (vSphere, OpenStack, user can also use their resources, e.g. VMs or On-premise). It also implements automated deployment and allows changing operation through Ansible, supporting Kubernetes clusters a full life-cycle self-defined management from Day 0 planning, Day 1 deployment, to Day 2 operating. KubeOperator really achieves the idea of Kubernetes as a service.

The following image is the whole architecture of KubeOperator:

![Architecture](https://kubeoperator.io/images/screenshot/ko-framework.svg) (or the new translated one)

> Note: KubeOperator passed the [Certified Kubernetes Conformance Program] (https://landscape.cncf.io/selected=kube-operator) provided by CNCF (Cloud Native Computing Foundation)

## Screenshot of the Web UI

![Web UI](https://kubeoperator.io/images/screenshot/6.png)

> [More features screenshot of KubeOperator](https://docs.kubeoperator.io/kubeoperator-v2.2/screenshot)

## Technology Advantages

- Easy to Use: Using a visible Web UI that significantly lower down the difficulty of K8s deployment and management, built-in with Webkubectl.
- Offline Support: Continue updating Kubernetes and common components of the offline pack
- Build by demand: Calling cloud platform API, build and deploy Kubernetes cluster in just a click
- Scale by demand: Swiftly scale Kubernetes clusters and improve resources utilization
- Patch by demand: rapid update, patch Kubernetes cluster and being up to date with the community version
- Self Repair: Through rebuilding malfunction node to confirm the usability of the cluster
- Full-Stack Monitoring: Full record of events, monitoring, warning and journaling from node, pod, to cluster
- Multi-AZ Support: Master nodes are distributed in different failure domain to make sure the high usability
- Marketplace: built-in with KubeApps Plus marketplace. Able to deploy and manage common apps quickly
- GPU Support: Support with GPU nodes which help operating high computation applications such as machine learning

## Features List

<table class=”subscription-level-table”>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-first-td-background-style” rowspan=”15”>Day 0 Planning</td>
        </td>
        <td class=”features-third-td-background-style” rowspan=”2”>Cluster Model
        </td>
        <td class=”features-third-td-background-style”>1 master node with N number of worker nodes : suitable for develop testing purpose
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style”>3 master nodes with N number of worker nodes : suitable for production-grade purpose
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style” rowspan=”3”>Calculation Scheme
        </td>
        <td class=”features-third-td-background-style”>Independent Host : support self-prepared VMs, public clouds or physical machines
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style”>vSphere Platform : Support auto-build host (using Terraform)
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style”>OpenStack Platform : Support auto-build host (using Terraform)
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style” rowspan=”3”>Storage Scheme
        </td>
        <td class=”features-third-td-background-style”>Independent host : Support NFS / Ceph RBD / Local Volume
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style”>vSphere Platform : Support vSphere Datastore (Centralized storage that compatible with vSAN & vSphere)
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style”>OpenStack Platform : Support OpenStack Cinder (Centralized storage that compatible with Ceph & Cinder)
        </td>
    </tr>
    <tr class=”subscription-level-tr-border”>
        <td class=”features-third-td-background-style” rowspan=”4”>Network Scheme
        </td>
        <td class="features-third-td-background-style">Support Flannel / Calico Network Plug-ins
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style"> Support internet expose service through F5 Big IP
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Support Traefik / Ingress-Nginx
        </td>
    </tr>    
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Support CoreDNS
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">GPU Scheme
        </td>
        <td class="features-third-td-background-style">Support NVIDIA GPU
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Operating System
        </td>
        <td class="features-third-td-background-style">Support RHEL/CentOS 7.4+
        </td>
    </tr>  
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Running on Container
        </td>
        <td class="features-third-td-background-style">Support Docker / containerd
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-first-td-background-style" rowspan="3">Day 1 Deploying
        </td>
        <td class="features-third-td-background-style" rowspan="3">Deployment
        </td>  
        <td class="features-third-td-background-style">Provide full installation package in an offline environment
        </td>         
    </tr>
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Support a visible screen of the deploying process
        </td>
    </tr>
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Support one-click automation deployment (using Ansible)
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-first-td-background-style" rowspan="21">Day 2 Operating
        </td>
        <td class="features-third-td-background-style" rowspan="9">Management
        </td>  
        <td class="features-third-td-background-style"> Support project-centralized hierarchical authorization management
        </td>         
    </tr>
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">3 roles: system admin, project admin and read-access user
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">Support docking with LDAP/AD
        </td>
    </tr>    
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">Expose with REST API
        </td>
    </tr>    
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style"> Install K8s Dashboard Management app through Kubeapps+
        </td>
    </tr>     
     <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style"> Install Weave Scope Management app through Kubeapps+
        </td>
    </tr>  
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">Support Web Kubectl UI
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">Built-in with Helm 
        </td>
    </tr>   
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">Constant updating certificate
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="4">Observable
        </td>
         <td class="features-third-td-background-style">Built-in with Prometheus, support fully monitoring & alarming of clusters, pods, nodes, and container
        </td>
    </tr>
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Built-in with Loki log system
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Built-in with Grafana for monitoring & logs display
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style"> Support notification center, signaling various cluster unusual events through DingTalk or WeChat
        </td>
    </tr>      
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Upgrade
        </td>
         <td class="features-third-td-background-style">Support whole cluster promotion
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Scale
        </td>
         <td class="features-third-td-background-style">Support flexible number of worker nodes
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Backup
        </td>
         <td class="features-third-td-background-style">Support periodical backup for etcd
        </td>
    </tr>  
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style"  rowspan="2">Safety compliance
        </td>
         <td class="features-third-td-background-style">Support score system for cluster’s health condition
        </td>
    </tr>   
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Support CSI Safe Scan
        </td>
    </tr>    
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="2">Kubeapps+
        </td>
         <td class="features-third-td-background-style">Support CI/CD tools, e.g. GitLab, Jenkins, Harbor, Argo CD 
        </td>
    </tr> 
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Support Machine Learning/AI applications like TensorFlow
        </td>
    </tr>    
 </table>

Complete Roadmap of versions history : [Roadmap](https://github.com/KubeOperator/KubeOperator/blob/master/ROADMAP.md)

## PowerPoint, Demo video, and Docs

-  [:books: PowerPoint]( ./KubeOperator_Intro.pdf) : Introduction of KubeOperator (should change the link to translated version)
-  [:tv: A short 8 mins demo]( https://kubeoperator-1256577600.file.myqcloud.com/video/KubeOperator2.1.mp4)
: Clearly explain the functionality of KubeOperator (maybe someone makes an English version?)
-  [:books: Installation & Doc](https://docs.kubeoperator.io/) : Docs about How to install KubeOperator, How to use, Features, Common Q&A etc

## Kubernetes offline installation package

KubeOperator provides a complete offline installation package (includes Kubernetes, Docker, etcd, Dashboard, Prometheus, OS patches, etc). Each installation package will be built as an independent container image for use by KubeOperator. Details : [k8s-package](https://github.com/KubeOperator/k8s-package)

## Community (Currently it’s only supported in Chinese)

- technical exchange QQ 群：825046920
- WeChat group： Search for WeChat ID wzk727912172，add friend，Remarks（City-GitHub username）, join in after verification
- technical consulting：https://jinshuju.net/f/QrZJpt

## Thanks to (Credits)

- [Terraform](https://github.com/hashicorp/terraform): Allowing to auto-build VMs；
- [Clarity](https://github.com/vmware/clarity/): Using as a front-end web framework；
- [Ansible](https://github.com/ansible/ansible): Using as an automated deployment tool；
- [kubeasz](https://github.com/easzlab/kubeasz): Providing various Kubernetes Ansible scripts；
- [Kubeapps](https://github.com/kubeapps/kubeapps): Creating a marketplace based on Kubeapps；

## License

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubeOperator is licensed under the Apache License, Version 2.0.

