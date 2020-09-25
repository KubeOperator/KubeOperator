import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { BaseModelComponent } from '../../../../shared/class/BaseModelComponent';
import { Notice } from '../notice';
import { NoticeService } from '../notice.service';
import { CommonAlertService } from '../../../../layout/common-alert/common-alert.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-mailbox-list',
  templateUrl: './mailbox-list.component.html',
  styleUrls: ['./mailbox-list.component.css']
})
export class MailboxListComponent extends BaseModelComponent<Notice> implements OnInit {

  // readColor = 'hsl(198, 100%, 32%)'; // normal color is #666666
  @Output() detailEvent = new EventEmitter<Notice>();
  @Output() onRead = new EventEmitter();

  unread: Notice[] = [];
  unreadAlert: number;
  unreadInfo: number;

  constructor(private noticeService: NoticeService,
              private commonAlertService: CommonAlertService,
              private translateService: TranslateService) {
    super(noticeService);
  }

  ngOnInit(): void {
    super.ngOnInit();
  }

  onDetail(item) {
    this.detailEvent.emit(item);
  }

  markAsRead() {
    this.onRead.emit();
  }

  readColor(value: boolean) {
    if (value) {
      return '#666666';
    } else { return 'hsl(198, 100%, 32%)'; }
  }

  checkUnread() {
    this.unreadAlert = 0;
    this.unreadInfo = 0;
    for (const item of this.items) {
      if (item.isRead === false) {
        if (item.level === 'info') {
          this.unreadInfo += 1;
        } else if (item.level === 'alert') {
          this.unreadAlert += 1;
        } else { console.log('your counting is wrong!'); }
      }
    }
    // this.noticeService.changeUnread({unreadInfo: this.unreadInfo, unreadAlert: this.unreadAlert});
  }


}
