import {Component, Input, OnInit} from '@angular/core';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';
import {Cluster} from '../../cluster/cluster';
import {Role} from '../../node/role';

@Component({
  selector: 'app-node-config',
  templateUrl: './node-config.component.html',
  styleUrls: ['./node-config.component.css']
})
export class NodeConfigComponent implements OnInit {

  @Input() currentCluster: Cluster;
  nodes: Node[] = [];
  roles: Role[] = [];

  constructor(private nodeService: NodeService) {
  }

  ngOnInit() {
    this.listNode();
    this.listRole();
  }

  listNode() {
    this.nodeService.listNodes(this.currentCluster.name)
      .subscribe(data => this.nodes = data);
  }

  listRole() {
    this.nodeService.listRoles(this.currentCluster.name)
      .subscribe(data => this.roles = data);
  }
}
