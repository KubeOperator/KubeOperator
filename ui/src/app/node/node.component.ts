import {Component, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NodeListComponent} from './node-list/node-list.component';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../cluster/cluster';

@Component({
  selector: 'app-node',
  templateUrl: './node.component.html',
  styleUrls: ['./node.component.css']
})
export class NodeComponent implements OnInit {


  @ViewChild(NodeListComponent, { static: true })
  listNode: NodeListComponent;

  public currentCluster: Cluster;

  constructor(private route: ActivatedRoute) {
  }


  refresh() {
    this.listNode.refresh();
  }

  ngOnInit() {
    this.route.parent.data.subscribe(data => {
      this.currentCluster = data['cluster'];
    });
  }


}
