import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Zone, ZoneUpdateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {NgForm} from '@angular/forms';
import {IpPool} from '../../ip-pool/ip-pool';
import {IpPoolService} from '../../ip-pool/ip-pool.service';


@Component({
    selector: 'app-zone-update',
    templateUrl: './zone-update.component.html',
    styleUrls: ['./zone-update.component.css']
})
export class ZoneUpdateComponent extends BaseModelDirective<Zone> implements OnInit {

    opened = false;
    item: ZoneUpdateRequest = new ZoneUpdateRequest();
    networkError = [];
    ipPools: IpPool[] = [];
    currentPool: IpPool = new IpPool();
    @Output() updated = new EventEmitter();
    @ViewChild('editForm') editForm: NgForm;

    constructor(private zoneService: ZoneService,
                private regionService: RegionService,
                private modalAlertService: ModalAlertService,
                private translateService: TranslateService,
                private commonAlertService: CommonAlertService,
                private ipPoolService: IpPoolService) {
        super(zoneService);
    }

    ngOnInit(): void {
    }

    open(item) {
        Object.assign(this.item, item);
        this.item.cloudVars = JSON.parse(item.vars);
        this.opened = true;
        this.ipPoolService.list().subscribe(res => {
            this.ipPools = res.items;
            this.changeIpPool(this.item.ipPoolName);
        }, error => {
        });
    }


    onCancel() {
        this.opened = false;
        this.networkError = [];
        this.editForm.resetForm(this.networkError);
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

    changeIpPool(ipPoolName) {
        for (const p of this.ipPools) {
            if (ipPoolName === p.name) {
                this.currentPool = p;
                break;
            }
        }
    }
}
