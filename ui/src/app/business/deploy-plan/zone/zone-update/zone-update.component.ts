import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Zone, ZoneUpdateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import * as ipaddr from 'ipaddr.js';
import {AlertLevels} from '../../../../layout/common-alert/alert';


@Component({
    selector: 'app-zone-update',
    templateUrl: './zone-update.component.html',
    styleUrls: ['./zone-update.component.css']
})
export class ZoneUpdateComponent extends BaseModelComponent<Zone> implements OnInit {

    opened = false;
    item: ZoneUpdateRequest = new ZoneUpdateRequest();
    networkError = [];
    @Output() updated = new EventEmitter();

    constructor(private zoneService: ZoneService, private regionService: RegionService, private modalAlertService: ModalAlertService,
                private translateService: TranslateService, private commonAlertService: CommonAlertService) {
        super(zoneService);
    }

    ngOnInit(): void {
    }

    open(item) {
        this.item = item;
        this.item.cloudVars = JSON.parse(item.vars);
        this.opened = true;
    }

    checkIp() {
        this.networkError = [];
        let result = true;
        const ipStart = this.item.cloudVars['ipStart'];
        const ipEnd = this.item.cloudVars['ipEnd'];
        if (!ipaddr.isValid(ipStart)) {
            result = false;
            this.networkError.push(this.translateService.instant('APP_IP_START_INVALID'));
        }
        if (!ipaddr.isValid(ipEnd)) {
            result = false;
            this.networkError.push(this.translateService.instant('APP_IP_END_INVALID'));
        }
        if (ipaddr.isValid(ipStart) && ipaddr.isValid(ipEnd)) {
            const start = ipaddr.parse(ipStart).toByteArray();
            const end = ipaddr.parse(ipEnd).toByteArray();
            for (let i = 0; i < 4; i++) {
                if (start[i] > end[i]) {
                    result = false;
                    this.networkError.push(this.translateService.instant('APP_IP_START_MUST'));
                    break;
                }
            }
        }
        return result;
    }

    onCancel() {
        this.opened = false;
    }

    onConfirm() {
        this.zoneService.update(this.item.name, this.item).subscribe(res => {
            this.onCancel();
            this.updated.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.onCancel();
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
