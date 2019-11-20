import {Component, OnInit} from '@angular/core';
import {DatePipe, DecimalPipe} from '@angular/common';
import {ClusterHealthService} from './cluster-health.service';
import {Cluster} from '../cluster/cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterHealth, Data, HealthData} from './cluster-health';
import {ClusterHealthHistory} from './cluster-health-history';

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
  projectId = ''
  clusterHealth: ClusterHealth = new ClusterHealth();
  clusterHealthHistories: ClusterHealthHistory[] = [];
  loading = true;
  totalRate = 0;
  error = false;
  timer;
  componentData = [];
  kubeSystemData = [];
  kubeOperatorData = [];
  healthData = [];
  errorMessage = '';

  ngOnInit() {
    this.clusterHealth.data = [];
    this.clusterHealth.rate = 100;
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.projectName = this.currentCluster.name;
      this.projectId = this.currentCluster.id;
      this.getClusterHealth();
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
    this.clusterHealthService.listClusterHealth(this.projectName).subscribe(res => {

      this.componentData = res.component;
      this.kubeSystemData = res['kube-system'];
      this.kubeOperatorData = res.monitoring;
      this.healthData = this.kubeSystemData.concat(this.kubeOperatorData);
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

}
