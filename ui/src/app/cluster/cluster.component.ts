import {Component, OnInit, ViewChild} from '@angular/core';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';

@Component({
  selector: 'app-cluster',
  templateUrl: './cluster.component.html',
  styleUrls: ['./cluster.component.css']
})
export class ClusterComponent implements OnInit {

  @ViewChild(ClusterCreateComponent, {static: true})
  creationCluster: ClusterCreateComponent;

  @ViewChild(ClusterListComponent, {static: true})
  listCluster: ClusterListComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openModal(): void {
    this.creationCluster.newCluster();
  }

  createCluster(created: boolean) {
    if (created) {
      this.listCluster.listCluster();
    }
  }
}
