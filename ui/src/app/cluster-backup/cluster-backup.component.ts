import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../cluster/cluster';
import {ClusterBackupStrategyComponent} from './cluster-backup-strategy/cluster-backup-strategy.component';
import {ClusterBackupListComponent} from './cluster-backup-list/cluster-backup-list.component';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../cluster/cluster.service';

@Component({
  selector: 'app-cluster-backup',
  templateUrl: './cluster-backup.component.html',
  styleUrls: ['./cluster-backup.component.css']
})
export class ClusterBackupComponent implements OnInit {

  @ViewChild(ClusterBackupStrategyComponent, {static: true})
  creation: ClusterBackupStrategyComponent;

  @ViewChild(ClusterBackupListComponent, {static: true})
  listClusterBackup: ClusterBackupListComponent;

  currentCluster: Cluster;


  constructor(private route: ActivatedRoute, private clusterService: ClusterService) {
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
      this.clusterService.getCluster(this.currentCluster.name).subscribe((d) => {
        this.currentCluster = d;
      });
    });
  }

}
