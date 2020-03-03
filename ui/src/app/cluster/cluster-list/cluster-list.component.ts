import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster, Operation} from '../cluster';
import {ClusterService} from '../cluster.service';
import {ActivatedRoute, Router} from '@angular/router';
import {SettingService} from '../../setting/setting.service';
import {PackageLogoService} from '../../package/package-logo.service';
import {ClusterStatusService} from '../cluster-status.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {SessionService} from '../../shared/session.service';
import {ItemResourceService} from '../../item-resource/item-resource.service';
import {ItemResourceDTO} from '../../item-resource/item-resource';
import {SessionUser} from '../../shared/session-user';

@Component({
  selector: 'app-cluster-list',
  templateUrl: './cluster-list.component.html',
  styleUrls: ['./cluster-list.component.css']
})
export class ClusterListComponent implements OnInit {
  loading = true;
  clusters: Cluster[] = [];
  deleteModal = false;
  hasHostname = false;
  selectedClusters: Cluster[] = [];
  @Output() addCluster = new EventEmitter<void>();
  itemName = '';
  permission: string;
  itemClusters: ItemResourceDTO[] = [];
  showItem = false;
  user: SessionUser;
  canCreate = false;


  constructor(private clusterService: ClusterService, private router: Router,
              private alertService: CommonAlertService, private settingService: SettingService,
              private packageLogoService: PackageLogoService, private route: ActivatedRoute,
              private clusterStatusService: ClusterStatusService, private sessionService: SessionService,
              private itemResourceService: ItemResourceService) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
    if (this.itemName !== undefined) {
      this.permission = this.sessionService.getItemPermission(this.itemName);
      this.listItemCluster();
    } else {
      this.showItem = true;
      this.listItemAndCluster();
    }
    this.getProfile();
    this.checkSetting();
  }

  checkSetting() {
    this.settingService.getSettings().subscribe(data => {
      this.hasHostname = !!data['local_hostname'];
    });
  }

  listItemCluster() {
    this.clusterService.listItemClusters(this.itemName).subscribe(data => {
      this.clusters = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      for (const itemCluster of this.itemClusters) {
        for (const i in this.clusters) {
          if (itemCluster.resource_id === this.clusters[i].id) {
            this.clusters[i].item_name = itemCluster.item_name;
          }
        }
      }
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
    }, res => {
      this.alertService.showAlert('删除失败' + res.error.msg, AlertLevels.ERROR);
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

  listItemAndCluster() {
    this.itemResourceService.getClusters().subscribe(res => {
      this.itemClusters = res['data'];
      this.listCluster();
    });
  }

  getProfile() {
    const profile = this.sessionService.getCacheProfile();
    this.user = profile.user;
    if (this.user.is_superuser) {
      this.canCreate = true;
    } else {
      for (const item of profile.items) {
        for (const rm of profile.item_role_mappings) {
          if (item.name === rm.item_name && rm.role !== 'VIEWER') {
            this.canCreate = true;
            break;
          }
        }
      }
    }
  }

}
