import { Component, EventEmitter, OnInit, Output, ViewChild } from '@angular/core';
import { BaseModelComponent } from '../../../../shared/class/BaseModelComponent';
import { Notice } from '../notice';
import { NoticeService } from '../notice.service';
import { ModalAlertService } from '../../../../shared/common-component/modal-alert/modal-alert.service';
import { CommonAlertService } from '../../../../layout/common-alert/common-alert.service';
import { TranslateService } from '@ngx-translate/core';
import { NgForm } from '@angular/forms';
import { AlertLevels } from '../../../../layout/common-alert/alert';

@Component({
  selector: 'app-mailbox-create',
  templateUrl: './mailbox-create.component.html',
  styleUrls: ['./mailbox-create.component.css']
})
export class MailboxCreateComponent extends BaseModelComponent<Notice> implements OnInit {

  opened = false;
  isSubmitGoing = false;
  // item: NoticeCreateRequest = new NoticeCreateRequest();
  item: Notice = new Notice();
  @ViewChild('noticeForm') noticeForm: NgForm;
  @Output() created = new EventEmitter();

  constructor(private noticeService: NoticeService,
              private modalAlertService: ModalAlertService,
              private commonAlertService: CommonAlertService,
              private translateService: TranslateService) {
    super(noticeService);
  }

  ngOnInit(): void {
  }

  open() {
    this.opened = true;
    this.item = new Notice();
    this.noticeForm.resetForm();
  }

  onCancel() {
    this.opened = false;
  }

  check(value: boolean): boolean {
    if (value === true){ return true; }
    else { return false;}
  }

  onSubmit() {

    const currentItem: Notice = {
      content: this.item.content,
      type: this.item.type,
      level: this.item.level,
      isRead: this.check(this.item.isRead),
      createdAt: new Date().toTimeString(),
      updatedAt: this.item.updatedAt
    };

    this.isSubmitGoing = true;
    this.noticeService.addItem(currentItem);
    this.opened = false;
    this.isSubmitGoing = false;
    this.created.emit();
    this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);


    // this.noticeService.create(this.item).subscribe(data => {
    //   this.opened = false;
    //   this.isSubmitGoing = false;
    //   this.created.emit();
    //   this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
    // }, error => {
    //   this.isSubmitGoing = false;
    //   this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
    // });
  }

}
