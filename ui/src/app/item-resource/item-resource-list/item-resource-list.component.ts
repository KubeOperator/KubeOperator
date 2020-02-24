import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ItemResourceService} from '../item-resource.service';
import {ActivatedRoute} from '@angular/router';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {ClusterStatusService} from "../../cluster/cluster-status.service";
import {ClusterHealthService} from "../../cluster-health/cluster-health.service";
import {SessionService} from "../../shared/session.service";

@Component({
  selector: 'app-item-resource-list',
  templateUrl: './item-resource-list.component.html',
  styleUrls: ['./item-resource-list.component.css']
})
export class ItemResourceListComponent implements OnInit {

  loading = false;
  @Output() add = new EventEmitter();
  itemName;
  itemId;
  itemResources;
  selected = [];
  clusterSelected = [];
  hostSelected = [];
  storageSelected = [];
  planSelected = [];
  backupStorageSelected = [];
  showDelete = false;
  resourceTypeName = '资源';
  isSubmitGoing = false;
  resourceType;
  permission;

  constructor(private itemResourceService: ItemResourceService, private route: ActivatedRoute,
              private alert: CommonAlertService, private clusterStatusService: ClusterStatusService,
              private sessionService: SessionService) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
    this.itemId = this.route.snapshot.queryParams['id'];
    this.permission = this.sessionService.getItemPermission(this.itemName);
    this.getItemResources();
  }

  createItemResource(resourceType) {
    this.add.emit(resourceType);
  }

  getItemResources() {
    this.loading = true;
    this.itemResourceService.getItemResources(this.itemName).subscribe(res => {
      this.itemResources = res;
      this.loading = false;
    });
  }

  deleteItemResource() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
      promises.push(this.itemResourceService.deleteItemResource(this.itemName, this.resourceType, item.resource_id).toPromise());
    });

    Promise.all(promises).then(data => {
      this.alert.showAlert('删除成功', AlertLevels.SUCCESS);
    }, res => {
      this.alert.showAlert('删除失败' + res.error.msg, AlertLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.selected = [];
        this.hostSelected = [];
        this.storageSelected = [];
        this.planSelected = [];
        this.backupStorageSelected = [];
        this.getItemResources();
      }
    );
  }

  cancelDelete() {
    this.showDelete = false;
  }

  openDeleteModal(selected, resourceType) {
    if (selected.length === 0) {
      this.alert.showAlert('请至少选择一行数据', AlertLevels.ERROR);
      return;
    }

    this.resourceType = resourceType;
    this.selected = selected;
    this.showDelete = true;
  }


  getDeployName(name: string) {
    switch (name) {
      case 'SINGLE':
        return '一主多节点';
      case 'MULTIPLE':
        return '多主多节点';
      default:
        return '无';
    }
  }
}
