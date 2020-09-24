import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Notice } from '../notice';
import { BaseModelComponent } from '../../../../shared/class/BaseModelComponent';
import { NoticeService } from '../notice.service';
import { ModalAlertService } from '../../../../shared/common-component/modal-alert/modal-alert.service';
import { CommonAlertService } from '../../../../layout/common-alert/common-alert.service';
import { TranslateService } from '@ngx-translate/core';
import { AlertLevels } from '../../../../layout/common-alert/alert';

@Component({
  selector: 'app-mailbox-delete',
  templateUrl: './mailbox-delete.component.html',
  styleUrls: ['./mailbox-delete.component.css']
})
export class MailboxDeleteComponent extends BaseModelComponent<Notice> implements OnInit {

  opened = false;
  items: Notice[] = [];
  @Output() deleted = new EventEmitter();

  constructor(private noticeService: NoticeService,
              private modalAlertService: ModalAlertService,
              private commonAlertService: CommonAlertService,
              private translateService: TranslateService) {
    super(noticeService);
  }

  ngOnInit(): void {
  }

  open(items: Notice[]) {
    this.opened = true;
    this.items = items;
  }

  onCancel() {
    this.opened = false;
  }

  onSubmit() {
    this.noticeService.deleteItems(this.items);
    console.log('delete successfully');
    this.opened = false;
    this.deleted.emit();
    this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);


    // this.noticeService.batch('delete', this.items).subscribe(data => {
    //   this.deleted.emit();
    //   this.opened = false;
    //   this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
    // }, error => {
    //   this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
    //   this.opened = false;
    // });
  }

}
