import {Component, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NodeService} from './node.service';
import {NodeCreateComponent} from './node-create/node-create.component';
import {NodeListComponent} from './node-list/node-list.component';

@Component({
  selector: 'app-node',
  templateUrl: './node.component.html',
  styleUrls: ['./node.component.css']
})
export class NodeComponent implements OnInit {

  @ViewChild(NodeCreateComponent)
  creationNode: NodeCreateComponent;

  @ViewChild(NodeListComponent)
  listNode: NodeListComponent;

  constructor() {
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
  }


}
