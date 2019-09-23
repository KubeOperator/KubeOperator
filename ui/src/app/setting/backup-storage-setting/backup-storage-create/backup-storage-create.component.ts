import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BackupStorage} from '../backup-storage';
import {NgForm} from '@angular/forms';
import {BackupStorageService} from '../backup-storage.service';
import {TipService} from '../../../tip/tip.service';
import {TipLevels} from '../../../tip/tipLevels';
import {StorageCredential} from '../storage-credential';

@Component({
  selector: 'app-backup-storage-create',
  templateUrl: './backup-storage-create.component.html',
  styleUrls: ['./backup-storage-create.component.scss']
})
export class BackupStorageCreateComponent implements OnInit {

  @Output() create = new EventEmitter<boolean>();
  item: BackupStorage = new BackupStorage();
  createOpened: boolean;
  isSubmitGoing = false;
  loading = true;
  staticBackDrop = true;
  closable = false;
  @ViewChild('backupStorageForm', {static: true}) backupStorageForm: NgForm;
  credential = new StorageCredential();
  invalid = false;
  tipShow = false;


  constructor(private backupStorageService: BackupStorageService, private tipService: TipService ) { }

  ngOnInit() {
  }


  onCancel() {
    this.createOpened = false;
    this.backupStorageForm.resetForm();
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.loading = true;
    if (this.credential == null) {
          this.invalid = true;
          this.tipShow = true;
    } else {
        this.item.credentials = this.credential;
        this.backupStorageService.checkBackupStorageConfig(this.item).subscribe(data => {
          this.invalid = !data.success;
          this.tipShow = true;
          if (data.success) {
              this.postItem(this.item);
          } else {
              this.isSubmitGoing = false;
          }
        }, err => {
          this.invalid = false;
          this.tipShow = true;
          this.isSubmitGoing = false;
        });
    }
  }

  postItem(credentials) {
    this.backupStorageService.createBackupStorage(credentials).subscribe(data => {
      this.createOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.tipService.showTip('新增成功!', TipLevels.SUCCESS);
      this.tipShow = false;
    }, err => {
      this.createOpened = true;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.tipService.showTip('新增失败!' + err.reson + 'state code:' + err.status, TipLevels.ERROR);
    });
  }


  newItem() {
    this.item = new BackupStorage();
    this.item.type = 'OSS';
    this.credential = new StorageCredential();
    this.createOpened = true;
  }

  checkValid(credential) {

  }

  closeTip() {
    this.tipShow = false;
  }
}
