import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {ItemResourceService} from '../item-resource.service';
import {ActivatedRoute} from '@angular/router';
import {ItemResource, ItemResourceDTO} from '../item-resource';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';

@Component({
  selector: 'app-item-resource-create',
  templateUrl: './item-resource-create.component.html',
  styleUrls: ['./item-resource-create.component.css']
})
export class ItemResourceCreateComponent implements OnInit {

  createOpened = false;
  itemName;
  resources: ItemResourceDTO[] = [];
  isSubmitGoing = false;
  itemId;
  resourceType;
  @ViewChild('resourceAlert', {static: true}) resourceAlert;
  loading = false;
  @Output() create = new EventEmitter<boolean>();


  constructor(private itemResourceService: ItemResourceService, private route: ActivatedRoute, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.itemName = this.route.snapshot.queryParams['name'];
    this.itemId = this.route.snapshot.queryParams['id'];
  }

  createItemResource(resourceType) {
    this.resourceType = resourceType;
    this.itemResourceService.getResources(this.itemName, resourceType).subscribe(res => {
      this.resources = res;
      if (this.resources.length > 0) {
        this.createOpened = true;
      } else {
        this.alert.showAlert('没有资源可授权', AlertLevels.ERROR);
      }
    });
  }

  onCancel() {
    this.createOpened = false;
  }

  onSubmit() {
    const itemResources = [];
    for (const resource of this.resources) {
      if (resource.checked) {
        const itemResource = new ItemResource();
        itemResource.resource_id = resource.resource_id;
        itemResource.resource_type = resource.resource_type;
        itemResource.item_id = this.itemId;
        itemResources.push(itemResource);
      }
    }
    if (itemResources.length === 0) {
      this.resourceAlert.showTip(true, '至少选择一个集群');
      this.sleep(1000).then(function (this) {
        this.resourceAlert.closeTip();
      });
      return;
    }

    this.isSubmitGoing = true;
    this.loading = true;

    this.itemResourceService.createItemResources(this.itemName, itemResources).subscribe(res => {
      this.alert.showAlert('授权成功', AlertLevels.SUCCESS);
      this.isSubmitGoing = false;
      this.createOpened = false;
      this.loading = false;
      this.create.emit(true);
    }, error => {
      this.alert.showAlert('授权失败', AlertLevels.ERROR);
      this.isSubmitGoing = false;
      this.loading = false;
      this.create.emit(true);
    });
  }

  sleep(ms) {
    return new Promise(
      (resolve) => {
        setTimeout(resolve, ms);
      });
  }

}
