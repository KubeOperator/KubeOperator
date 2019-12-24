import {Component, OnInit} from '@angular/core';
import {DatePipe, DecimalPipe} from '@angular/common';
import {ClusterHealthService} from './cluster-health.service';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterHealth} from './cluster-health';

@Component({
  selector: 'app-cluster-health',
  templateUrl: './cluster-health.component.html',
  styleUrls: ['./cluster-health.component.css'],
  providers: [DatePipe, DecimalPipe]
})
export class ClusterHealthComponent implements OnInit {

  constructor(private route: ActivatedRoute, private decimalPipe: DecimalPipe,
              private datePipe: DatePipe, private clusterHealthService: ClusterHealthService) {
  }

  options: {};
  time: any;
  currentCluster: Cluster;
  projectName = '';
  projectId = '';
  clusterHealth: ClusterHealth = new ClusterHealth();
  loading = true;
  error = false;
  timer;
  componentData = [];
  healthData = [];
  errorMessage = '';
  namespaces = [];
  namespace = 'all';

  ngOnInit() {
    this.clusterHealth.data = [];
    this.clusterHealth.rate = 100;
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectName = this.currentCluster.name;
      this.projectId = this.currentCluster.id;
      this.getClusterHealth();
      this.getClusterNamespace();
      this.getComponent();
    });
    this.timer = setInterval(() => {
      this.getClusterHealth();
    }, 300000);
  }

  // tslint:disable-next-line:use-lifecycle-interface
  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  }

  getClusterHealth() {
    this.loading = true;
    this.clusterHealthService.listClusterHealth(this.projectName, this.namespace).subscribe(res => {
      this.healthData = res.pod_data;
      this.loading = false;
      if (res.message !== '') {
        this.errorMessage = res.message;
        this.error = true;
      }
    }, error1 => {
      this.clusterHealth.data = [];
      this.clusterHealth.rate = 0;
      this.loading = false;
      this.error = true;
    });
  }

  getClusterNamespace() {
    this.clusterHealthService.listNamespace(this.projectName).subscribe(res => {
      this.namespaces = res;
    }, error1 => {

    });
  }

  getComponent() {
    this.clusterHealthService.listComponent(this.projectName).subscribe(res => {
      this.componentData = res;
    }, error1 => {

    });
  }
}
