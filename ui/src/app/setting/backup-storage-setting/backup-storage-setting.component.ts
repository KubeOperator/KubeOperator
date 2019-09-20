import {Component, OnInit, ViewChild} from '@angular/core';
import {BackupStorageCreateComponent} from './backup-storage-create/backup-storage-create.component';
import {BackupStorageListComponent} from './backup-storage-list/backup-storage-list.component';

@Component({
  selector: 'app-backup-storage-setting',
  templateUrl: './backup-storage-setting.component.html',
  styleUrls: ['./backup-storage-setting.component.scss']
})
export class BackupStorageSettingComponent implements OnInit {

  @ViewChild(BackupStorageCreateComponent, {static: true})
  creation: BackupStorageCreateComponent;

  @ViewChild(BackupStorageListComponent, {static: true})
  listBackupStorage: BackupStorageListComponent;


  constructor() { }

  ngOnInit() {
  }

  openModal() {
    this.creation.newItem();
  }

  create(created: boolean) {
    if (created) {
      this.refresh();
    }
  }

  refresh() {
    this.listBackupStorage.refresh();
  }
}
