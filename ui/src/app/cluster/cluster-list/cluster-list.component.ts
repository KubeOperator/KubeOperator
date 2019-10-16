import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster, Operation} from '../cluster';
import {ClusterService} from '../cluster.service';
import {Router} from '@angular/router';
import {SettingService} from '../../setting/setting.service';
import {PackageLogoService} from '../../package/package-logo.service';
import {ClusterStatusService} from '../cluster-status.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

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
              private alertService: CommonAlertService, private settingService: SettingService,
              private packageLogoService: PackageLogoService,
              private clusterStatusService: ClusterStatusService) {
  }

  ngOnInit() {
    this.listCluster();
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
      this.listCluster();
      this.alertService.showAlert('删除集群成功！', AlertLevels.SUCCESS);
    }).finally(() => {
      this.deleteModal = false;
      this.selectedClusters = [];
    });
  }


  addNewCluster() {
    this.addCluster.emit();
  }


  getStatusComment(status: string): string {
    return this.clusterStatusService.getComment(status);
  }

  getDeployTypeComment(type: string): string {
    return this.clusterStatusService.getDeployType(type);
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


}
