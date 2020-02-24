import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../cluster/cluster.service';
import {Cluster} from '../cluster/cluster';
import {Router} from '@angular/router';
import {DashboardSearch} from './dashboardSearch';
import {DashboardService} from './dashboard.service';
import {ItemService} from '../item/item.service';

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
  cpuUsage = 0;
  memUsage = 0;
  cpuTotal = 0;
  memTotal = 0;
  showPodDetail = false;
  showContainerDetail = false;
  showClusterUsageDetail = false;
  showErrorLokiContainerDetail = false;
  showErrorPodDetail = false;
  nodes = [];
  timer;
  maxPodCount = 0;
  containerUsage = 0;
  showWindow = true;
  items = [];

  constructor(private clusterService: ClusterService, private router: Router,
              private dashboardService: DashboardService, private itemService: ItemService) {
  }

  ngOnInit() {
    this.dashboardSearch.cluster = 'all';
    this.dashboardSearch.item = 'all';
    this.dashboardSearch.dateLimit = 1;
    this.getItems();
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

  dataInit() {
    this.clusterData = [];
    this.podCount = 0;
    this.nodeCount = 0;
    this.namespaceCount = 0;
    this.deploymentCount = 0;
    this.containerCount = 0;
    this.restartPods = [];
    this.warnContainers = [];
    this.errorPods = [];
    this.cpuUsage = 0;
    this.memUsage = 0;
    this.cpuTotal = 0;
    this.memTotal = 0;
    this.showPodDetail = false;
    this.showContainerDetail = false;
    this.showClusterUsageDetail = false;
    this.showErrorLokiContainerDetail = false;
    this.showErrorPodDetail = false;
    this.nodes = [];
    this.maxPodCount = 0;
    this.containerUsage = 0;
  }

  listCluster() {
    this.clusterService.listItemClusters(this.dashboardSearch.item).subscribe(data => {
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
    this.dataInit();
    this.dashboardService.getDashboard(this.dashboardSearch.cluster, this.dashboardSearch.item).subscribe(data => {
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
        this.cpuTotal = this.cpuTotal + d['cpu_total'];
        this.memTotal = this.memTotal + d['mem_total'];
        this.cpuUsage = this.cpuUsage + d['cpu_usage'];
        this.memUsage = this.memUsage + d['mem_usage'];
        if (d['cpu_total'] === 0 && d['mem_total'] === 0) {
          count--;
        }
        this.nodes = this.nodes.concat(d['nodes']);

        for (const cluster of this.clusters) {
          if (cluster.name === d['name']) {
            let max_pod = cluster.configs['MAX_PODS'];
            if (max_pod === undefined) {
              max_pod = 110;
            }
            const all_max_pod = max_pod * cluster.nodes.length;
            this.maxPodCount = this.maxPodCount + all_max_pod;
          }
        }
      }
      if (count > 0) {
        this.cpuUsage = this.cpuUsage / count * 100;
        this.memUsage = this.memUsage / count * 100;
      }
      this.containerUsage = this.podCount / this.maxPodCount * 100;
      this.loading = false;
    }, error1 => {
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
      const linkUrl = ['', url];
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

  getItems() {
    this.itemService.listItem().subscribe(data => {
      this.items = data;
      this.dashboardSearch.item = this.items[0].name;
      this.search();
    });
  }

  changeItem() {
    this.dashboardSearch.cluster = 'all';
    this.listCluster();
  }
}
