import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Notice } from '../notice';
import { BaseModelComponent } from '../../../../shared/class/BaseModelComponent';
import { NoticeService } from '../notice.service';
import { ModalAlertService } from '../../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
  selector: 'app-mailbox-detail',
  templateUrl: './mailbox-detail.component.html',
  styleUrls: ['./mailbox-detail.component.css']
})
export class MailboxDetailComponent extends BaseModelComponent<Notice> implements OnInit {

  opened = false;
  item: Notice = new Notice();
  loading = false;
  @Output() detail = new EventEmitter();

  constructor(private noticeService: NoticeService, private modalAlertService: ModalAlertService) {
    super(noticeService);
  }

  ngOnInit(): void {
  }

  onCancel() {
    this.item = new Notice();
    this.opened = false;
    this.loading = false;
    this.detail.emit();
  }

  open(item: Notice) {
    this.opened = true;
    this.item = item;
  }

  onSync() {
    console.log('You click onSync button!');
  }

}
