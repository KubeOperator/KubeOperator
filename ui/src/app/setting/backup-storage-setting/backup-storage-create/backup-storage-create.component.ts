import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BackupStorage} from '../backup-storage';
import {NgForm} from '@angular/forms';
import {BackupStorageService} from '../backup-storage.service';
import {StorageCredential} from '../storage-credential';
import {CommonAlertService} from '../../../base/header/common-alert.service';
import {AlertLevels} from '../../../base/header/components/common-alert/alert';

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
  message = '';
  buckets = [];
  @ViewChild('alertModal', {static: true}) alertModal;


  constructor(private backupStorageService: BackupStorageService, private alertService: CommonAlertService) {
  }

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
    } else {
      this.item.credentials = this.credential;
      this.backupStorageService.checkBackupStorageConfig(this.item).subscribe(data => {
        this.invalid = !data.success;
        this.message = data.message;
        if (data.success) {
          this.postItem(this.item);
        } else {
          this.isSubmitGoing = false;
        }
        this.alertModal.showTip(this.invalid, this.message);
      }, err => {
        this.isSubmitGoing = false;
        this.alertModal.showTip(true, '校验失败!');
      });
    }
  }

  postItem(credentials) {
    this.backupStorageService.createBackupStorage(credentials).subscribe(data => {
      this.createOpened = false;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.alertService.showAlert('新增成功!', AlertLevels.SUCCESS);
    }, err => {
      this.createOpened = true;
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.loading = false;
      this.alertService.showAlert('新增失败!' + err.reson + 'state code:' + err.status, AlertLevels.ERROR);
    });
  }


  newItem() {
    this.item = new BackupStorage();
    this.item.type = 'OSS';
    this.credential = new StorageCredential();
    this.createOpened = true;
  }

  changeType() {
    this.buckets = [];
    this.alertModal.closeTip();
  }

  getBuckets(credential) {
    this.item.credentials = credential;
    this.backupStorageService.getBuckets(this.item).subscribe(rep => {
      this.invalid = !rep.success;
      if (rep.success) {
        this.buckets = rep.data;
        this.message = '查询成功';
      } else {
        this.message = '查询失败';
        this.buckets = [];
      }
      this.alertModal.showTip(this.invalid, this.message);
    }, err => {
      this.buckets = [];
      this.alertModal.showTip(true, '查询失败');
    });
  }
}
