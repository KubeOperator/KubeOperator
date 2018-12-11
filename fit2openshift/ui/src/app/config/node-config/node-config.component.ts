import {Component, Input, OnInit} from '@angular/core';
import {NodeService} from '../../node/node.service';
import {Node} from '../../node/node';

@Component({
  selector: 'app-node-config',
  templateUrl: './node-config.component.html',
  styleUrls: ['./node-config.component.css']
})
export class NodeConfigComponent implements OnInit {

  @Input() clusterId: string;
  nodes: Node[] = [];

  constructor(private nodeService: NodeService) {
  }

  ngOnInit() {
    this.listNode();
  }

  listNode() {
    // this.nodeService.listNodes(this.clusterId).subscribe(data => this.nodes = data);
    const node: Node = new Node();
    node.name = 'master-1';
    node.ip = '172.101.1.1';
    node.roles = 'master';
    node.status = 'running';
    this.nodes.push(node);

    const node2: Node = new Node();
    node2.name = 'master-1';
    node2.ip = '172.101.1.1';
    node2.roles = 'master';
    node2.status = 'running';
    this.nodes.push(node2);
  }

}
