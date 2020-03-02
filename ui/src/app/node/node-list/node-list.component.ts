import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {NodeService} from '../node.service';
import {Node} from '../node';
import {Cluster} from '../../cluster/cluster';
import {NodeDetailComponent} from '../node-detail/node-detail.component';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {DashboardService} from '../../dashboard/dashboard.service';
import {SessionService} from "../../shared/session.service";

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
  clusterData;
  permission;

  constructor(private nodeService: NodeService, private alertService: CommonAlertService,
              private dashboardService: DashboardService, private sessionService: SessionService) {
  }

  ngOnInit() {
    this.permission = this.sessionService.getItemPermission(this.currentCluster.item_name);
    this.listNodes();
  }

  listNodes() {
    this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
      this.nodes = data.filter((node) => {
        return node.name !== 'localhost' && node.name !== '127.0.0.1' && node.name !== '::1';
      });
      this.dashboardService.getDashboard(this.currentCluster.name, this.currentCluster.item_name).subscribe(res => {
        this.clusterData = res.data;
        const nodeList = JSON.parse(res.data[0])['nodes'];
        nodeList.forEach(n => {
          this.nodes.forEach(node => {
            if (n.name === node.name) {
              node.cpu_usage = Number(n.cpu_usage) * 100;
              node.mem_usage = Number(n.mem_usage) * 100;
            }
          });
        });
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

  checkNodes() {
    this.nodeService.checkNodes(this.currentCluster.name).subscribe(data => {
      this.alertService.showAlert('同步成功', AlertLevels.SUCCESS);
      this.refresh();
    }, error1 => {
      this.alertService.showAlert('同步失败', AlertLevels.ERROR);
    });
  }
}
