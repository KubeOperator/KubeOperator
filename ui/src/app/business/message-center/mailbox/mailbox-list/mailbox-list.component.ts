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

  readColor = 'hsl(198, 100%, 32%)'; // normal color is #666666
  @Output() detailEvent = new EventEmitter<Notice>();

  constructor(private noticeService: NoticeService,
              private commonAlertService: CommonAlertService,
              private translateService: TranslateService) {
    super(noticeService);
  }

  ngOnInit(): void {
    super.ngOnInit();
    this.items = this.noticeService.getItems();
  }

  onDetail(item) {
    this.detailEvent.emit(item);
  }


}
