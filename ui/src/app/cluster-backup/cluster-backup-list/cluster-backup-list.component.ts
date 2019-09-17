import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';

@Component({
  selector: 'app-cluster-backup-list',
  templateUrl: './cluster-backup-list.component.html',
  styleUrls: ['./cluster-backup-list.component.scss']
})
export class ClusterBackupListComponent implements OnInit {
  loading = true;
  @Input() currentCluster: Cluster;

  constructor() { }

  ngOnInit() {
  }

}
