import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NodeService} from '../node.service';
import {Node} from '../node';

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {

  loading = true;
  nodes: Node[] = [];
  selectedRow: Node[] = [];
  @Input() clusterId: string;
  @Output() addNode = new EventEmitter();

  constructor(private nodeService: NodeService) {
  }

  ngOnInit() {
    this.listNodes();
  }

  listNodes() {
    // this.nodeService.listNodes(this.clusterId).subscribe(data => {
    //   this.nodes = data;
    //   this.loading = false;
    // }, error => {
    //   this.loading = false;
    // });
    const node: Node = new Node();
    node.name = 'master-1';
    node.ip = '172.101.1.1';
    node.roles = 'master';
    node.status = 'running';
    this.nodes.push(node);
    this.loading = false;

    const node2: Node = new Node();
    node2.name = 'master-1';
    node2.ip = '172.101.1.1';
    node2.roles = 'master';
    node2.status = 'running';
    this.nodes.push(node2);
  }

  refresh() {
    this.listNodes();
  }

  addNewNode() {
    this.addNode.emit();
  }
}
