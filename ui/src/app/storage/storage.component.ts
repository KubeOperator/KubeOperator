import {Component, OnInit, ViewChild} from '@angular/core';
import {StorageCreateComponent} from './components/storage-create/storage-create.component';
import {StorageListComponent} from './components/storage-list/storage-list.component';

@Component({
  selector: 'app-storage',
  templateUrl: './storage.component.html',
  styleUrls: ['./storage.component.css']
})
export class StorageComponent implements OnInit {
  @ViewChild(StorageCreateComponent)
  creationStorage: StorageCreateComponent;

  @ViewChild(StorageListComponent)
  listStorage: StorageListComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openModal(): void {
    this.creationStorage.newItem();
  }

  createItem(created: boolean) {
    if (created) {
      this.listStorage.listStorage();
    }
  }
}
