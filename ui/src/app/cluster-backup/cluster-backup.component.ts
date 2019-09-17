import { Component, OnInit } from '@angular/core';
import {Cluster} from '../cluster/cluster';

@Component({
  selector: 'app-cluster-backup',
  templateUrl: './cluster-backup.component.html',
  styleUrls: ['./cluster-backup.component.css']
})
export class ClusterBackupComponent implements OnInit {
  currentCluster: Cluster;

  constructor() { }

  ngOnInit() {
  }

}
