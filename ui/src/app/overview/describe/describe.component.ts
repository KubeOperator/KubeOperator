import {Component, Input, OnInit} from '@angular/core';
import {Cluster, Operation} from '../../cluster/cluster';
import {PackageService} from '../../package/package.service';
import {ClusterInfo, Portal, Template} from '../../package/package';
import {ClusterService} from '../../cluster/cluster.service';
import {OverviewService} from '../overview.service';
import {OperaterService} from '../../deploy/component/operater/operater.service';
import {Router} from '@angular/router';
import {ClusterStatus} from './class/describe';

@Component({
  selector: 'app-describe',
  templateUrl: './describe.component.html',
  styleUrls: ['./describe.component.css']
})
export class DescribeComponent implements OnInit {

  @Input() currentCluster: Cluster;
  clusterInfos: ClusterInfo[] = [];
  operations: Operation[] = [];

  constructor(private packageService: PackageService, private clusterService: ClusterService,
              private overviewService: OverviewService, private operaterService: OperaterService,
              private router: Router) {
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
    this.overviewService.getClusterToken(this.currentCluster);
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

  redirect(cluster_name: string, url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', 'cluster', cluster_name, url];
      this.router.navigate(linkUrl);
    }
  }

  getStatus(): ClusterStatus {
    return this.getStatusDescribe(this.currentCluster);
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

}
