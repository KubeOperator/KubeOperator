import { Component, OnInit } from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {ActivatedRoute, Router} from '@angular/router';
import {ClusterService} from '../../cluster/cluster.service';
import {OperaterService} from '../../deploy/component/operater/operater.service';


@Component({
  selector: 'app-cluster-backup-strategy',
  templateUrl: './cluster-backup-strategy.component.html',
  styleUrls: ['./cluster-backup-strategy.component.scss']
})
export class ClusterBackupStrategyComponent implements OnInit {

  currentCluster: Cluster;
    constructor(private route: ActivatedRoute, private clusterService: ClusterService,
              private router: Router, private operaterService: OperaterService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.refreshCluster();
    });
  }
   refreshCluster() {
    this.clusterService.getCluster(this.currentCluster.name).subscribe(cluster => {
      this.currentCluster = cluster;
    });
  }
}
