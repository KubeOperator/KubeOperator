import {Component, OnInit, ViewChild} from '@angular/core';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {Cluster} from './cluster';

@Component({
  selector: 'app-cluster',
  templateUrl: './cluster.component.html',
  styleUrls: ['./cluster.component.css']
})
export class ClusterComponent implements OnInit {

  @ViewChild(ClusterCreateComponent)
  creationCluster: ClusterCreateComponent;

  @ViewChild(ClusterListComponent)
  listCluster: ClusterListComponent;

  loading = true;

  constructor() {
  }

  ngOnInit() {
  }

  openModal(): void {
    this.creationCluster.newCluster();
    Cluster
  }
}
