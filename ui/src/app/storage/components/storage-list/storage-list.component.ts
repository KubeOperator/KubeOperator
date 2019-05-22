import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Storage} from '../../models/storage';
import {StorageService} from '../../services/storage.service';
import {StorageTemplate} from '../../models/storage-template';
import {StorageTemplateService} from '../../services/storage-template.service';
import {TipService} from '../../../tip/tip.service';
import {TipLevels} from '../../../tip/tipLevels';
import {StorageDetailComponent} from '../storage-detail/storage-detail.component';

@Component({
  selector: 'app-storage-list',
  templateUrl: './storage-list.component.html',
  styleUrls: ['./storage-list.component.css']
})
export class StorageListComponent implements OnInit {

  constructor(private storageService: StorageService, private storageTemplateService: StorageTemplateService,
              private tipService: TipService) {
  }

  resourceTypeName = '存储';
  loading = true;
  items: Storage[] = [];
  selectedItems: Storage[] = [];
  storageTemplates: StorageTemplate[] = [];
  showDelete = false;
  showDetail = false;
  @ViewChild(StorageDetailComponent)
  child: StorageDetailComponent;
  @Output() addItem = new EventEmitter<void>();

  ngOnInit() {
    this.getStorageTemplates();
    this.listStorage();
  }

  getStorageTemplates() {
    this.storageTemplateService.listStorageTemplates().subscribe(data => {
      this.storageTemplates = data;
    });
  }

  listStorage() {
    this.loading = true;
    this.storageService.listStorage().subscribe(data => {
      this.items = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  onDetail(item) {
    this.showDetail = true;
    this.child.item = item;
    this.child.loadTemplate();
  }

  onDeleted() {
    this.showDelete = true;
  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selectedItems.forEach(item => {
        promises.push(this.storageService.deleteStorage(item.name).toPromise());
      }
    );
    Promise.all(promises).then(data => {
      this.tipService.showTip('删除成功', TipLevels.SUCCESS);
    }, error => {
      this.tipService.showTip('删除失败' + error.toString(), TipLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.listStorage();
      }
    );
  }

  addNewItem() {
    this.addItem.emit();
  }
}
