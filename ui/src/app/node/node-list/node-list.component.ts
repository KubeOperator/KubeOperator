import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NodeService} from '../node.service';
import {Node} from '../node';
import {Cluster} from '../../cluster/cluster';
import {NodeDetailComponent} from '../node-detail/node-detail.component';
import {HostInfoComponent} from '../../host/host-info/host-info.component';

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {

  loading = true;
  nodes: Node[] = [];
  showDetail = false;
  @ViewChild(NodeDetailComponent, {static: true}) child: NodeDetailComponent;
  @Input() currentCluster: Cluster;
  timeResult;
  openView = false;
  loadingTime = false;

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

  openInfo(node: Node) {
    this.showDetail = true;
    this.child.node = node;
  }

  toGrafana() {
    const url = 'http://grafana.apps.' + this.currentCluster.name + '.' + this.currentCluster.cluster_doamin_suffix + '/explore';
    window.open(url, '_blank');
  }

  syncTime() {
    this.loadingTime = true;
    this.openView = true;
    this.nodeService.syncHostTime(this.currentCluster.name).subscribe(data => {
      this.timeResult = data;
      this.loadingTime = false;
    }, error1 => {
      this.loadingTime = false;
    });
  }
}
