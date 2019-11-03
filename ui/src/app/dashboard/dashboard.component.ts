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
  dashboardSearch: DashboardSearch = new DashboardSearch();
  clusterData = [];
  podCount = 0;
  nodeCount = 0;
  namespaceCount = 0;
  deploymentCount = 0;
  containerCount = 0;
  restartPods = [];

  constructor(private clusterService: ClusterService, private router: Router, private dashboardService: DashboardService) {
  }

  ngOnInit() {
    this.dashboardSearch.cluster = 'all';
    this.dashboardSearch.dateLimit = 1;
    this.search();
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      this.getClusterData();
    }, error => {
      this.loading = false;
    });
  }

  getCluster() {
    this.clusterService.getCluster(this.dashboardSearch.cluster).subscribe(data => {
      this.clusters = [];
      this.clusters.push(data);
      this.loading = false;
    });
  }

  getClusterData() {
    this.dashboardService.getDashboard(this.clusters[1].name).subscribe(data => {
      this.clusterData = JSON.parse(data.data);
      for (const d of this.clusterData) {
        this.podCount = this.podCount + d['pods'].length;
        this.namespaceCount = this.namespaceCount + d['namespaces'].length;
        this.deploymentCount = this.deploymentCount + d['deployments'].length;
        this.nodeCount = this.nodeCount + d['nodes'].length;
        for (const p of d['pods']) {
          this.containerCount = this.containerCount + p['containers'].length;
          if (p['restart_count'] > 0) {
            this.restartPods.push(p);
          }
        }
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
}
