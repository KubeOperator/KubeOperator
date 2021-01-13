import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Notice} from '../notice';
import {NoticeService} from '../notice.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {SessionService} from '../../../../shared/auth/session.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-mailbox-list',
    templateUrl: './mailbox-list.component.html',
    styleUrls: ['./mailbox-list.component.css']
})
export class MailboxListComponent extends BaseModelDirective<Notice> implements OnInit {

    @Output() detailEvent = new EventEmitter<Notice>();
    items: Notice[] = [];
    user;
    loading = false;

    constructor(private noticeService: NoticeService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(noticeService);
    }

    ngOnInit(): void {
        this.listByUsername();
    }

    listByUsername() {
        this.loading = true;
        this.noticeService.pageBy(this.page, this.size).subscribe(res => {
            this.items = res.items;
            this.total = res.total;
            this.loading = false;
        }, error => {
            this.loading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    onDetail(item: Notice) {
        this.detailEvent.emit(item);
        if (item.readStatus === 'UNREAD') {
            item.readStatus = 'READ';
            const readItems = [];
            readItems.push(item);
            this.service.batch('update', readItems).subscribe(data => {
            }, error => {
                this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            });
        }
    }

    markAsRead() {
        this.service.batch('update', this.selected).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            this.selected = [];
            this.listByUsername();
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
