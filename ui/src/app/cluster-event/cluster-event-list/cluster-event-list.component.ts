import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from '../../cluster/cluster';

@Component({
  selector: 'app-cluster-event-list',
  templateUrl: './cluster-event-list.component.html',
  styleUrls: ['./cluster-event-list.component.css']
})
export class ClusterEventListComponent implements OnInit {
  @Input() currentCluster: Cluster;

  constructor() {
  }

  ngOnInit() {
  }

}
