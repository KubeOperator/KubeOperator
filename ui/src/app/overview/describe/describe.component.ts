import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster, Operation} from '../../cluster/cluster';
import {PackageService} from '../../package/package.service';
import {ClusterInfo, Portal, Template} from '../../package/package';
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

@Component({
  selector: 'app-describe',
  templateUrl: './describe.component.html',
  styleUrls: ['./describe.component.css']
})
export class DescribeComponent implements OnInit {

  @Input() currentCluster: Cluster;
  clusterInfos: ClusterInfo[] = [];
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

  constructor(private packageService: PackageService, private clusterService: ClusterService,
              private overviewService: OverviewService, private operaterService: OperaterService,
              private router: Router, private clusterStatusService: ClusterStatusService,
              private nodeService: NodeService) {
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
      this.redirect('deploy');
    });
    this.confirmAlert.close();
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', 'cluster', this.currentCluster.name, url];
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
}
