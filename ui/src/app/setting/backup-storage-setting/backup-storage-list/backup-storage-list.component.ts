import {Component, OnInit} from '@angular/core';
import {BackupStorageService} from '../backup-storage.service';
import {BackupStorage} from '../backup-storage';
import {TipService} from '../../../tip/tip.service';
import {TipLevels} from '../../../tip/tipLevels';

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

  constructor(private backupStorageService: BackupStorageService, private tipService: TipService) {
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
      this.tipService.showTip('删除成功', TipLevels.SUCCESS);
    }, error => {
      this.tipService.showTip('删除失败', TipLevels.ERROR);
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
}
