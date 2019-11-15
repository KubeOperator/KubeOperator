import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../cluster/cluster.service';
import {Cluster} from '../cluster/cluster';
import {Router} from '@angular/router';
import {DashboardSearch} from './dashboardSearch';
import {DashboardService} from './dashboard.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  loading = true;
  clusters: Cluster[] = [];
  selectClusters: Cluster[] = [];
  dashboardSearch: DashboardSearch = new DashboardSearch();
  clusterData = [];
  podCount = 0;
  nodeCount = 0;
  namespaceCount = 0;
  deploymentCount = 0;
  containerCount = 0;
  restartPods = [];
  warnContainers = [];
  errorLokiContainers = [];
  errorPods = [];
  cpu_usage = 0;
  mem_usage = 0;
  cpu_total = 0;
  mem_total = 0;
  show_pod_detail = false;
  show_container_detail = false;
  show_cluster_usage_detail = false;
  show_error_loki_container_detail = false;
  show_error_pod_detail = false;
  nodes = [];
  timer;

  constructor(private clusterService: ClusterService, private router: Router, private dashboardService: DashboardService) {
  }

  ngOnInit() {
    this.dashboardSearch.cluster = 'all';
    this.dashboardSearch.dateLimit = 1;
    this.search();
    this.timer = setInterval(() => {
      this.search();
    }, 300000);
  }

  // tslint:disable-next-line:use-lifecycle-interface
  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  data_init() {
    this.clusterData = [];
    this.podCount = 0;
    this.nodeCount = 0;
    this.namespaceCount = 0;
    this.deploymentCount = 0;
    this.containerCount = 0;
    this.restartPods = [];
    this.warnContainers = [];
    this.errorPods = [];
    this.cpu_usage = 0;
    this.mem_usage = 0;
    this.cpu_total = 0;
    this.mem_total = 0;
    this.show_pod_detail = false;
    this.show_container_detail = false;
    this.show_cluster_usage_detail = false;
    this.show_error_loki_container_detail = false;
    this.nodes = [];
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      this.selectClusters = data;
      this.getClusterData();
    }, error => {
      this.loading = false;
    });
  }

  getCluster() {
    this.clusterService.getCluster(this.dashboardSearch.cluster).subscribe(data => {
      this.clusters = [];
      this.clusters.push(data);
      this.getClusterData();
    });
  }

  getClusterData() {
    this.data_init();
    this.dashboardService.getDashboard(this.dashboardSearch.cluster).subscribe(data => {
      this.clusterData = data.data;
      this.restartPods = data.restartPods;
      this.warnContainers = data.warnContainers;
      this.errorLokiContainers = data.errorLokiContainers;
      this.errorPods = data.errorPods;
      let count = this.clusterData.length;
      for (const cd of this.clusterData) {
        const d = JSON.parse(cd);
        this.podCount = this.podCount + d['pods'].length;
        this.namespaceCount = this.namespaceCount + d['namespaces'].length;
        this.deploymentCount = this.deploymentCount + d['deployments'].length;
        for (const p of d['pods']) {
          this.containerCount = this.containerCount + p['containers'].length;
        }
        this.nodeCount = this.nodeCount + d['nodes'].length;
        this.cpu_total = this.cpu_total + d['cpu_total'];
        this.mem_total = this.mem_total + d['mem_total'];
        this.cpu_usage = this.cpu_usage + d['cpu_usage'];
        this.mem_usage = this.mem_usage + d['mem_usage'];
        if (d['cpu_total'] === 0 && d['mem_total'] === 0) {
          count--;
        }
        this.nodes = this.nodes.concat(d['nodes']);
      }
      if (this.clusterData.length > 0) {
        this.cpu_usage = this.cpu_usage / count * 100;
        this.mem_usage = this.mem_usage / count * 100;
      }
      this.loading = false;
    });
  }

  search() {
    this.loading = true;
    if (this.dashboardSearch.cluster === 'all') {
      this.listCluster();
    } else {
      this.getCluster();
    }
  }

  refresh() {
    this.search();
  }

  toPage(url) {
    this.redirect(url);
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', url];
      this.router.navigate(linkUrl);
    }
  }

  toGrafana(cluster_name) {
    for (const cluster of this.clusters) {
      if (cluster_name === cluster.name) {
        const url = 'http://grafana.apps.' + cluster.name + '.' + cluster.cluster_doamin_suffix + '/explore';
        window.open(url, '_blank');
      }
    }
  }
}
