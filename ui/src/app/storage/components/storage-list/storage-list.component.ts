import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Storage} from '../../models/storage';
import {StorageService} from '../../services/storage.service';
import {StorageTemplate} from '../../models/storage-template';
import {StorageTemplateService} from '../../services/storage-template.service';

@Component({
  selector: 'app-storage-list',
  templateUrl: './storage-list.component.html',
  styleUrls: ['./storage-list.component.css']
})
export class StorageListComponent implements OnInit {

  constructor(private storageService: StorageService, private storageTemplateService: StorageTemplateService) {
  }

  loading = true;
  items: Storage[] = [];
  selectedItems: Storage[] = [];
  storageTemplates: StorageTemplate[] = [];

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
    this.storageService.getStorage().subscribe(data => {
      this.items = data;
      this.loading = false;
    }, error => {
      this.loading = false;
    });
  }

  addNewItem() {
    this.addItem.emit();
  }

}
