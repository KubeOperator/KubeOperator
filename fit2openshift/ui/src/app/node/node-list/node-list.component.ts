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
    this.nodeService.listNodes(this.clusterId).subscribe(data => {
      this.nodes = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  refresh() {
    this.listNodes();
  }

  addNewNode() {
    this.addNode.emit();
  }
}
