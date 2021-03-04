import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Zone} from '../zone';
import {ZoneService} from '../zone.service';
import {Region} from '../../region/region';
import {SystemService} from '../../../setting/system.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-zone-list',
    templateUrl: './zone-list.component.html',
    styleUrls: ['./zone-list.component.css']
})
export class ZoneListComponent extends BaseModelDirective<Zone> implements OnInit {

    @Output() detailEvent = new EventEmitter<Region>();

    constructor(
        private zoneService: ZoneService,
        private settingService: SystemService,
        private commonAlert: CommonAlertService,
        private translateService: TranslateService) {
        super(zoneService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onCreate() {
        this.settingService.getRegistry().subscribe(data => {
            if (data.items === null) {
                this.commonAlert.showAlert(this.translateService.instant('APP_NOT_SET_SYSTEM_IP_X86'), AlertLevels.ERROR);
                return
            }
            let isRepoExit: boolean = false;
            for (const repo of data.items) {
                if (repo.architecture === 'x86_64') {
                    isRepoExit = true;
                    break;
                }
            }
            if (!isRepoExit) {
                this.commonAlert.showAlert(this.translateService.instant('APP_NOT_SET_SYSTEM_IP_X86'), AlertLevels.ERROR);
                return
            }
            this.createEvent.emit();
        });
    }

    onDetail(item) {
        this.detailEvent.emit(item);
    }
}
