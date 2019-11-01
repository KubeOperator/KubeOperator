import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BackupStorageService} from '../backup-storage.service';
import {BackupStorage} from '../backup-storage';
import {StorageCredential} from '../storage-credential';
import {CommonAlertService} from '../../../base/header/common-alert.service';
import {AlertLevels} from '../../../base/header/components/common-alert/alert';


@Component({
  selector: 'app-backup-storage-list',
  templateUrl: './backup-storage-list.component.html',
  styleUrls: ['./backup-storage-list.component.scss']
})
export class BackupStorageListComponent implements OnInit {

  loading = true;
  showDelete = false;
  items: BackupStorage[] = [];
  selected: BackupStorage[] = [];
  resourceTypeName: '备份';
  @Output() add = new EventEmitter();
  credential = new StorageCredential();

  constructor(private backupStorageService: BackupStorageService, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.listItems();
  }

  listItems() {
    this.loading = true;
    this.backupStorageService.listBackupStorage().subscribe(data => {
      this.items = data;
      this.loading = false;
    });

  }

  delete() {
    const promises: Promise<{}>[] = [];
    this.selected.forEach(item => {
      promises.push(this.backupStorageService.deleteBackupStorage(item.name).toPromise());
    });

    Promise.all(promises).then(data => {
      this.alert.showAlert('删除成功', AlertLevels.SUCCESS);
    }, res => {
      this.alert.showAlert('删除失败' + res.error.msg , AlertLevels.ERROR);
    }).finally(
      () => {
        this.showDelete = false;
        this.selected = [];
        this.listItems();
      }
    );
  }

  refresh() {
    this.listItems();
  }

  addItem() {
    this.add.emit();
  }

  getBucket(item) {
    if (item.type === 'AZURE') {
      return item.credentials.container;
    } else {
      return item.credentials.bucket;
    }
  }
}
