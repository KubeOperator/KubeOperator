import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {NodeService} from '../node.service';
import {Node} from '../node';
import {Cluster} from '../../cluster/cluster';

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {

  loading = true;
  nodes: Node[] = [];
  selectedRow: Node[] = [];
  @Input() currentCluster: Cluster;
  @Output() addNode = new EventEmitter();

  constructor(private nodeService: NodeService) {
  }

  ngOnInit() {
    this.listNodes();
  }

  listNodes() {
    this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
      this.nodes = data.filter((node) => {
        return node.name !== 'localhost' && node.name !== '127.0.0.1' && node.name !== '::1';
      });
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
