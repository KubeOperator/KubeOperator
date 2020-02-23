import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ItemService} from '../item.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {Item} from '../item';

@Component({
  selector: 'app-item-list',
  templateUrl: './item-list.component.html',
  styleUrls: ['./item-list.component.css']
})
export class ItemListComponent implements OnInit {

  @Output() addItem = new EventEmitter<void>();

  constructor(private itemService: ItemService, private alertService: CommonAlertService) {
  }

  items: Item[] = [];
  loading = false;
  selectedItems: any = [];
  deleteModal = false;

  ngOnInit() {
    this.listItem();
  }

  listItem() {
    this.loading = true;
    this.itemService.listItem().subscribe(res => {
      this.items = res;
      this.loading = false;
    });
  }

  addNewItem() {
    this.addItem.emit();
  }

  onDeleted() {
    this.deleteModal = true;
  }

  confirmDelete() {
    const promises: Promise<{}>[] = [];
    this.selectedItems.forEach(cluster => {
      promises.push(this.itemService.deleteItem(cluster.name).toPromise());
    });
    Promise.all(promises).then(() => {
      this.listItem();
      this.alertService.showAlert('删除成功！', AlertLevels.SUCCESS);
    }, res => {
      this.alertService.showAlert('删除失败' + res.error.msg, AlertLevels.ERROR);
    }).finally(() => {
      this.deleteModal = false;
      this.selectedItems = [];
    });
  }
}
