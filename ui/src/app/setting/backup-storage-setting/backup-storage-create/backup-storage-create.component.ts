import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BackupStorage} from '../backup-storage';
import {NgForm} from '@angular/forms';
import {BackupStorageService} from '../backup-storage.service';
import {StorageCredential} from '../storage-credential';
import {CommonAlertService} from '../../../base/header/common-alert.service';
import {AlertLevels} from '../../../base/header/components/common-alert/alert';
import {ModalAlertComponent} from '../../../shared/common-component/modal-alert/modal-alert.component';

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
  message = '';
  buckets = [];


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
      this.tipShow = true;
    } else {
      this.item.credentials = this.credential;
      this.backupStorageService.checkBackupStorageConfig(this.item).subscribe(data => {
        // @ts-ignore
        this.invalid = !data.success;
        // @ts-ignore
        this.message = data.message;
        this.tipShow = true;
        // @ts-ignore
        if (data.success) {
          this.postItem(this.item);
        } else {
          this.isSubmitGoing = false;
        }
      }, err => {
        this.invalid = true;
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
      this.alertService.showAlert('新增成功!', AlertLevels.SUCCESS);
      this.tipShow = false;
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

  closeTip() {
    this.tipShow = false;
  }

  changeType() {
    this.buckets = [];
    this.closeTip();
  }

  getBuckets(credential) {
    this.item.credentials = credential;
    this.backupStorageService.getBuckets(this.item).subscribe(rep => {
      // @ts-ignore
      this.invalid = !rep.success;
      // @ts-ignore
      if (rep.success) {
        // @ts-ignore
        this.buckets = rep.data;
        this.message = '查询成功';
      } else {
        this.tipShow = true;
        // @ts-ignore
        this.message = '查询失败';
        this.buckets = [];
      }
    }, err => {
      this.invalid = true;
      this.tipShow = true;
      this.message = '查询失败';
      this.buckets = [];
    });
  }
}
