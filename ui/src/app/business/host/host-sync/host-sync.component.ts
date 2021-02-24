import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Host, HostSync} from '../host';
import {CommonAlertService} from '../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';
import { HostService } from '../host.service';

@Component({
    selector: 'app-host-sync',
    templateUrl: './host-sync.component.html',
    styleUrls: ['./host-sync.component.css']
})
export class HostSyncComponent implements OnInit {
    constructor(
        private hostService: HostService,
        private commonAlertService: CommonAlertService,
        private translateService: TranslateService,
    ){}

    opened = false;
    hosts: Host[] = [];
    hostSyncList: HostSync[] = [];
    isSubmitGoing = false;

    @Output() sync = new EventEmitter();

    ngOnInit(): void {
    }

    open(items) {
        this.hosts = items;
        this.hostSyncList = this.hosts.map(function (item) {
            let hostItem = new HostSync;
            hostItem.HostName = item.name;
            hostItem.HostStatus = item.status;
            return hostItem;
        })
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.hostService.syncList(this.hostSyncList).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.sync.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_SYNC_START_SUCCESS'), AlertLevels.SUCCESS);
        })
    }
}
