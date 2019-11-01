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
      this.loading = false;
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

  search() {
    this.loading = true;
    if (this.dashboardSearch.cluster === 'all') {
      this.listCluster();
    } else {
      this.getCluster();
    }
  }

  refresh() {
    // this.search();
    this.dashboardService.getDashboard(this.clusters[1].name).subscribe(data => {
      console.log('1111');
    });
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
