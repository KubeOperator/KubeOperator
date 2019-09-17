import {Component, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NodeService} from './node.service';
import {NodeCreateComponent} from './node-create/node-create.component';
import {NodeListComponent} from './node-list/node-list.component';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../cluster/cluster';

@Component({
  selector: 'app-node',
  templateUrl: './node.component.html',
  styleUrls: ['./node.component.css']
})
export class NodeComponent implements OnInit {

  @ViewChild(NodeCreateComponent, { static: true })
  creationNode: NodeCreateComponent;

  @ViewChild(NodeListComponent, { static: true })
  listNode: NodeListComponent;

  public currentCluster: Cluster;

  constructor(private route: ActivatedRoute) {
  }

  openModal() {
    this.creationNode.newNode();
  }

  createNode(created: boolean) {
    if (created) {
      this.refresh();
    }
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
