import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../cluster/cluster.service';
import {Cluster} from '../cluster/cluster';
import {Router} from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  loading = true;
  clusters: Cluster[] = [];

  constructor(private clusterService: ClusterService, private router: Router) {
  }

  ngOnInit() {
    this.listCluster();
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      this.loading = false;
    }, error => {
      this.loading = false;
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
