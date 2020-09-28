import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Notice } from '../notice';
import { BaseModelDirective } from '../../../../shared/class/BaseModelDirective';
import { NoticeService } from '../notice.service';
import { ModalAlertService } from '../../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
  selector: 'app-mailbox-detail',
  templateUrl: './mailbox-detail.component.html',
  styleUrls: ['./mailbox-detail.component.css']
})
export class MailboxDetailComponent extends BaseModelDirective<Notice> implements OnInit {

  opened = false;
  item: Notice = new Notice();
  @Output() detail = new EventEmitter();

  constructor(private noticeService: NoticeService, private modalAlertService: ModalAlertService) {
    super(noticeService);
  }

  ngOnInit(): void {
  }

  onCancel() {
    this.item = new Notice();
    this.opened = false;
  }

  open(item: Notice) {
    this.opened = true;
    this.item = item;
  }


}
