import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {StorageProvisioner, ProvisionerSync} from '../storage-provisioner';
import {CommonAlertService} from '../../../../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';
import { StorageProvisionerService } from '../storage-provisioner.service';

@Component({
    selector: 'app-storage-provisioner-sync',
    templateUrl: './storage-provisioner-sync.component.html',
    styleUrls: ['./storage-provisioner-sync.component.css']
})
export class StorageProvisionerSyncComponent implements OnInit {
    constructor(
        private service: StorageProvisionerService,
        private commonAlertService: CommonAlertService,
        private translateService: TranslateService,
    ){}

    opened = false;
    provisioners: StorageProvisioner[] = [];
    provisionerList: ProvisionerSync[] = [];
    isSubmitGoing = false;

    @Output() synced = new EventEmitter();
    @Input() clusterName: string;

    ngOnInit(): void {
    }

    open(items) {
        this.provisioners = items;
        this.provisionerList = this.provisioners.map(function (item) {
            let hostItem = new ProvisionerSync;
            hostItem.name = item.name;
            hostItem.type = item.type;
            hostItem.status = item.status;
            return hostItem;
        })
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.service.syncList(this.clusterName, this.provisionerList).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.synced.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_PROVISIONER_SYNC_START_SUCCESS'), AlertLevels.SUCCESS);
        })
    }
}
