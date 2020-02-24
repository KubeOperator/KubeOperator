import {Component, OnInit, ViewChild} from '@angular/core';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {ActivatedRoute} from '@angular/router';
import {SessionService} from '../shared/session.service';

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

  itemName: string;
  permission: string;

  constructor(private route: ActivatedRoute, private sessionService: SessionService) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
    this.permission = this.sessionService.getItemPermission(this.itemName);
  }

  openModal(): void {
    this.creationCluster.newCluster(this.itemName);
  }

  createCluster(created: boolean) {
    if (created) {
      this.listCluster.listCluster();
    }
  }
}
