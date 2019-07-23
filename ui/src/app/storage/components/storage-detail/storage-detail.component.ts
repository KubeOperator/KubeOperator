import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {StorageTemplateService} from '../../services/storage-template.service';
import {Storage} from '../../models/storage';
import {StorageTemplate} from '../../models/storage-template';
import {StorageService} from '../../services/storage.service';

@Component({
  selector: 'app-storage-detail',
  templateUrl: './storage-detail.component.html',
  styleUrls: ['./storage-detail.component.css']
})
export class StorageDetailComponent implements OnInit {

  @Input()
  opened = false;
  @Output()
  openedChange = new EventEmitter();
  item: Storage;
  storageTemplate: StorageTemplate;
  onUpdating = false;

  constructor(private storageTempateService: StorageTemplateService, private storageService: StorageService) {
  }

  ngOnInit() {
  }

  loadTemplate() {
    this.storageTempateService.getStorageTemplate(this.item.template).subscribe(data => {
      this.storageTemplate = data;
    });
  }


  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  update() {
    if (this.onUpdating) {
      return;
    }
    this.onUpdating = true;
    this.storageService.updateStorage(this.item.name, this.item).subscribe(data => {
      this.item = data;
      this.onUpdating = false;
    });
  }

  getStatus(item: Storage) {
    return this.storageService.getStorageStatus(item);
  }

}
