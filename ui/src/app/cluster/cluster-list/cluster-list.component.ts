import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Cluster, Operation} from '../cluster';
import {ClusterService} from '../cluster.service';
import {Router} from '@angular/router';
import {TipService} from '../../tip/tip.service';
import {TipLevels} from '../../tip/tipLevels';
import {MessageService} from '../../base/message.service';
import {MessageLevels} from '../../base/message/message-level';
import {SettingService} from '../../setting/setting.service';
import {PackageLogoService} from '../../package/package-logo.service';
import {ClusterStatusService} from '../cluster-status.service';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {el} from '@angular/platform-browser/testing/src/browser_util';

@Component({
  selector: 'app-cluster-list',
  templateUrl: './cluster-list.component.html',
  styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent implements OnInit {
  loading = true;
  clusters: Cluster[] = [];
  deleteModal = false;
  selectedClusters: Cluster[] = [];

  @Output() addCluster = new EventEmitter<void>();

  constructor(private clusterService: ClusterService, private router: Router,
              private tipService: TipService, private messageService: MessageService, private settingService: SettingService,
              private packageLogoService: PackageLogoService,
              private clusterStatusService: ClusterStatusService, private operaterService: OperaterService) {
  }

  ngOnInit() {
    this.listCluster();
    this.checkSetting();
  }

  checkSetting() {
    // this.settingService.getSetting('local_hostname').subscribe(data => {
    //   if (!data.value || data.value === '127.0.0.1') {
    //     this.messageService.announceMessage('部署前请先设置主机IP,否则部署将造成失败！', MessageLevels.WARN);
    //   }
    // });
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }


  onDeleted() {
    this.deleteModal = true;
  }

  confirmDelete() {
    const promises: Promise<{}>[] = [];
    this.selectedClusters.forEach(cluster => {
      promises.push(this.clusterService.deleteCluster(cluster.name).toPromise());
    });
    Promise.all(promises).then(() => {
      this.deleteModal = false;
      this.listCluster();
      this.tipService.showTip('删除集群成功！', TipLevels.SUCCESS);
    }, (error) => {
      this.tipService.showTip('删除集群失败:' + error, TipLevels.ERROR);
    });
  }


  addNewCluster() {
    this.addCluster.emit();
  }

  goToLink(clusterName: string) {
    const linkUrl = ['kubeOperator', 'cluster', clusterName, 'overview'];
    this.router.navigate(linkUrl);
  }

  getLogo(resource: string) {
    return this.packageLogoService.getLogo(resource);
  }

  getStatusComment(status: string): string {
    return this.clusterStatusService.getComment(status);
  }

  redirect(cluster_name: string, url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', 'cluster', cluster_name, url];
      this.router.navigate(linkUrl);
    }
  }

  showBtn(cluster: Cluster, opt: Operation): boolean {
    let result = true;
    if (opt.display_on) {
      if (!opt.display_on.includes(cluster.status)) {
        result = false;
      }
    }
    return result;
  }


  handleEvent(cluster_name: string, opt: Operation) {
    if (opt.event) {
      this.operaterService.executeOperate(cluster_name, opt.event).subscribe(() => {
        this.redirect(cluster_name, opt.redirect);
      });
    } else if (opt.redirect) {
      this.redirect(cluster_name, opt.redirect);
    }
  }
}
