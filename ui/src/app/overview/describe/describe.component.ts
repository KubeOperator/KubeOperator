import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../cluster/cluster';
import {PackageService} from '../../package/package.service';
import {ClusterService} from '../../cluster/cluster.service';
import {OverviewService} from '../overview.service';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {Router} from '@angular/router';
import {NodeService} from '../../node/node.service';
import {ClusterStatusService} from '../../cluster/cluster-status.service';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';
import {UpgradeComponent} from '../upgrade/upgrade.component';
import {WebkubectlComponent} from '../webkubectl/webkubectl.component';
import * as clipboard from 'clipboard-polyfill';
import {DashboardService} from '../../dashboard/dashboard.service';

@Component({
  selector: 'app-describe',
  templateUrl: './describe.component.html',
  styleUrls: ['./describe.component.css']
})
export class DescribeComponent implements OnInit {

  @Input() currentCluster: Cluster;
  openToken = false;
  status: string;
  openChangeStatus = false;
  token: string = null;
  event: string = null;
  openHost = false;
  openConfigs = false;
  workers = [];
  workerIp = '';
  @ViewChild(ConfirmAlertComponent, {static: true}) confirmAlert: ConfirmAlertComponent;
  @ViewChild(UpgradeComponent, {static: true}) upgrade: UpgradeComponent;
  @ViewChild(WebkubectlComponent, {static: true}) webKubeCtrl: WebkubectlComponent;
  @ViewChild('alertModal', {static: true}) alertModal;
  @ViewChild('tokenAlert', {static: true}) tokenAlert;
  nodeList = [];
  cpuUsage = 0;
  memUsage = 0;
  cpuTotal = 0;
  memTotal = 0;
  containerCount = 0;
  containerPercent = 0;
  podCount = 0;
  nodeCount = 0;
  namespaceCount = 0;
  deploymentCount = 0;
  baseRoute;


  constructor(private packageService: PackageService, private clusterService: ClusterService,
              private overviewService: OverviewService, private operaterService: OperaterService,
              private router: Router, private clusterStatusService: ClusterStatusService,
              private nodeService: NodeService, private dashboardService: DashboardService) {
  }

  ngOnInit() {
    this.nodeService.listNodes(this.currentCluster.name).subscribe(data => {
      this.workers = data.filter((node) => {
        return node.roles.includes('worker');
      });
      if (this.workers.length > 0) {
        this.workerIp = this.workers[0].ip;
      }
    });
    this.getClusterData();
    this.baseRoute = 'item/' + this.currentCluster.item_name + '/cluster/' + this.currentCluster.name;
  }


  onDownload() {
    this.overviewService.downLoad(this.currentCluster);
  }

  onGetToken() {
    this.openToken = true;
    this.overviewService.getClusterToken(this.currentCluster).subscribe(data => {
      this.token = data.token;
    });
  }

  openWebkubectl() {
    this.webKubeCtrl.loading = true;
    this.webKubeCtrl.opened = true;
    this.clusterService.getWebkubectlToken(this.currentCluster.id).subscribe(data => {
      this.webKubeCtrl.url = 'http://' + window.location.host + '/webkubectl/terminal/?token=' + data['token'];
      this.webKubeCtrl.open();
    });
  }


  onInstall() {
    this.confirmAlert.setTitle('确认安装');
    this.confirmAlert.setComment('安装即将开始，请确认所有配置已就绪');
    this.event = 'install';
    this.confirmAlert.opened = true;
  }

  onUninstall() {
    this.confirmAlert.setTitle('确认卸载');
    this.confirmAlert.setComment('卸载操作可能造成您的数据丢失，是否继续 ?');
    this.event = 'uninstall';
    this.confirmAlert.opened = true;
  }

  onUpgrade() {
    this.upgrade.opened = true;
    this.upgrade.currentPackageName = this.currentCluster.package;
    this.upgrade.reset();
    this.event = 'upgrade';
  }

  onCancel() {
    this.openChangeStatus = false;
  }

  onConfirm() {
    this.clusterService.changeStatus(this.status, this.currentCluster.name).subscribe(data => {
      this.currentCluster = data;
      this.openChangeStatus = false;
    });
  }

  handleUpgrade() {
    const params = {'package': this.upgrade.newPackage.name};
    this.handleEvent(params);
  }


  handleEvent(params?) {
    this.operaterService.executeOperate(this.currentCluster.name, this.event, params).subscribe(() => {
      this.router.navigate([this.baseRoute + '/deploy']);
    });
    this.confirmAlert.close();
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['cluster', this.currentCluster.name, url];
      this.router.navigate(linkUrl);
    }
  }

  getDeployTypeComment(type: string): string {
    return this.clusterStatusService.getDeployType(type);
  }

  copy(text: string) {
    const success = clipboard.writeText(text);
    if (success) {
      const alert = this.alertModal;
      const tokenAlert = this.tokenAlert;
      this.alertModal.showTip(false, '复制成功');
      this.tokenAlert.showTip(false, '复制成功');
      this.sleep(1000).then(function (this) {
        alert.closeTip();
        tokenAlert.closeTip();
      });
    }
  }

  sleep(ms) {
    return new Promise(
      (resolve) => {
        setTimeout(resolve, ms);
      });
  }

  getClusterData() {
    this.dashboardService.getDashboard(this.currentCluster.name).subscribe(res => {
      const clusterData = res.data;
      const data = JSON.parse(clusterData[0]);
      this.nodeList = data['nodes'];
      this.cpuTotal = data['cpu_total'];
      this.memTotal = data['mem_total'];
      this.cpuUsage = data['cpu_usage'] * 100;
      this.memUsage = data['mem_usage'] * 100;
      for (const p of data['pods']) {
        this.containerCount = this.containerCount + p['containers'].length;
      }
      let max_pod = this.currentCluster.configs['MAX_PODS'];
      if (max_pod === undefined) {
        max_pod = 110;
      }
      this.containerPercent = this.containerCount / (max_pod * this.nodeList.length) * 100;
      this.podCount = data['pods'].length;
      this.namespaceCount = data['namespaces'].length;
      this.deploymentCount = data['deployments'].length;
      this.nodeCount = data['nodes'].length;
    });
  }

  toApp(app) {
    const url = 'http://' + app + '.apps.' + this.currentCluster.name + '.' + this.currentCluster.cluster_doamin_suffix;
    window.open(url, '_blank');
  }
}
