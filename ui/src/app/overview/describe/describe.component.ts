import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster, Operation} from '../../cluster/cluster';
import {PackageService} from '../../package/package.service';
import {ClusterInfo, Portal, Template} from '../../package/package';
import {ClusterService} from '../../cluster/cluster.service';
import {OverviewService} from '../overview.service';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {Router} from '@angular/router';
import {ClusterStatus} from './class/describe';
import {ClusterStatusService} from '../../cluster/cluster-status.service';
import {ConfirmAlertComponent} from '../../shared/common-component/confirm-alert/confirm-alert.component';
import {ClusterListComponent} from '../../cluster/cluster-list/cluster-list.component';

@Component({
  selector: 'app-describe',
  templateUrl: './describe.component.html',
  styleUrls: ['./describe.component.css']
})
export class DescribeComponent implements OnInit {

  @Input() currentCluster: Cluster;
  clusterInfos: ClusterInfo[] = [];
  operations: Operation[] = [];
  openToken = false;
  token: string = null;
  event: string = null;
  @ViewChild(ConfirmAlertComponent, { static: true }) confirmAlert: ConfirmAlertComponent;

  constructor(private packageService: PackageService, private clusterService: ClusterService,
              private overviewService: OverviewService, private operaterService: OperaterService,
              private router: Router, private clusterStatusService: ClusterStatusService) {
  }

  ngOnInit() {
    this.packageService.getPackage(this.currentCluster.package).subscribe(pkg => {
      const infos = pkg.meta.cluster_infos;
      this.operations = pkg.meta.operations;
      this.clusterService.listClusterConfig(this.currentCluster.name).subscribe(configs => {
        infos.forEach(info => {
          configs.forEach(cfg => {
            if (cfg.key === info.key) {
              info.value = cfg.value;
            }
          });
        });
        this.clusterInfos = infos;
      });
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

  testOpen() {
    this.confirmAlert.opened = true;
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

  handleEvent() {
    this.operaterService.executeOperate(this.currentCluster.name, this.event).subscribe(() => {
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


  getStatusDescribe(cluster: Cluster): ClusterStatus {
    const result = new ClusterStatus();
    switch (cluster.status) {
      case 'READY':
        result.color = 'red';
        result.alias = '准备安装';
        break;
      case 'ERROR':
        result.color = 'red';
        result.alias = '错误';
        break;
      case 'RUNNING':
        result.color = 'green';
        result.alias = '运行中';
        break;
      default :
        result.color = 'blue';
        result.alias = '执行中';
    }
    return result;
  }

  getDeployTypeComment(type: string): string {
    return this.clusterStatusService.getDeployType(type);
  }

}
