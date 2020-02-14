import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ItemResourceService} from '../item-resource.service';
import {ActivatedRoute} from '@angular/router';
import {AlertLevels} from "../../base/header/components/common-alert/alert";
import {BackupStorageService} from "../../setting/backup-storage-setting/backup-storage.service";
import {CommonAlertService} from "../../base/header/common-alert.service";

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
  showDelete = false;
  resourceTypeName = '资源';
  isSubmitGoing = false;
  resourceType;

  constructor(private itemResourceService: ItemResourceService, private route: ActivatedRoute, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
    this.itemId = this.route.snapshot.queryParams['id'];
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
      promises.push(this.itemResourceService.deleteItemResource(item.resource_id).toPromise());
    });

    Promise.all(promises).then(data => {
      this.alert.showAlert('删除成功', AlertLevels.SUCCESS);
    }, res => {
      this.alert.showAlert('删除失败' + res.error.msg, AlertLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.selected = [];
        this.getItemResources();
      }
    );
  }

  cancelDelete() {
    this.showDelete = false;
  }

  openDeleteModal(resourceType) {
    this.resourceType = resourceType;
    this.showDelete = true;
  }
}
