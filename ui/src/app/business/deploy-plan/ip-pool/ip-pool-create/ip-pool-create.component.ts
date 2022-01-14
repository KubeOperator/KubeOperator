import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {IpPool, IpPoolCreate} from '../ip-pool';
import {IpPoolService} from '../ip-pool.service';
import {NamePattern} from '../../../../constant/pattern';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import * as ipaddr from 'ipaddr.js';
import {IpPattern} from '../../../../constant/pattern';

@Component({
    selector: 'app-ip-pool-create',
    templateUrl: './ip-pool-create.component.html',
    styleUrls: ['./ip-pool-create.component.css']
})
export class IpPoolCreateComponent extends BaseModelDirective<IpPool> implements OnInit {


    @Output() created = new EventEmitter();
    opened = false;
    item: IpPoolCreate = new IpPoolCreate();
    namePattern = NamePattern;
    isSubmitGoing = false;
    networkValid = false;
    ipPattern = IpPattern;

    constructor(private ipPoolService: IpPoolService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(ipPoolService);
    }

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.networkValid = false;
        this.item = new IpPoolCreate();
    }

    onCancel() {
        this.opened = false;
    }


    checkNetwork() {
        const ipStart = this.item.ipStart;
        const ipEnd = this.item.ipEnd;
        if (!ipaddr.isValid(ipStart) || (!ipaddr.isValid(ipEnd))) {
            this.networkValid = false;
            this.modalAlertService.showAlert(this.translateService.instant('APP_IP_RANGE_INVALID'), AlertLevels.ERROR);
            return;
        }
        const ipStartAddr = ipaddr.IPv4.parse(ipStart);
        const ipEndAddr = ipaddr.IPv4.parse(ipEnd);
        const start = ipStartAddr.toByteArray();
        const end = ipEndAddr.toByteArray();
        for (let i = 0; i < 4; i++) {
            if (start[i] > end[i]) {
                this.networkValid = false;
                this.modalAlertService.showAlert(this.translateService.instant('APP_IP_RANGE_INVALID'), AlertLevels.ERROR);
                return;
            }
            if (i === 3 && (end[i] - start[i]) < 1) {
                this.networkValid = false;
                this.modalAlertService.showAlert(this.translateService.instant('APP_IP_RANGE_INVALID'), AlertLevels.ERROR);
                return;
            }
        }
        const subnet = this.item.subnet.split('/', 2);
        if (subnet.length !== 2) {
            this.networkValid = false;
            this.modalAlertService.showAlert(this.translateService.instant('APP_SUBNET_INVALID'), AlertLevels.ERROR);
            return;
        }
        if (!ipEndAddr.match(ipaddr.IPv4.parseCIDR(this.item.subnet))) {
            this.networkValid = false;
            this.modalAlertService.showAlert(this.translateService.instant('APP_IP_RANGE_INVALID'), AlertLevels.ERROR);
            return;
        }
        const gateway = this.item.gateway;
        if (!ipaddr.isValid(gateway)) {
            this.networkValid = false;
            this.modalAlertService.showAlert(this.translateService.instant('APP_GATEWAY_INVALID'), AlertLevels.ERROR);
            return;
        }
        const dns1 = this.item.dns1;
        const dns2 = this.item.dns2;
        if (!ipaddr.isValid(dns1) || (!ipaddr.isValid(dns2))) {
            this.networkValid = false;
            this.modalAlertService.showAlert(this.translateService.instant('APP_DNS_INVALID'), AlertLevels.ERROR);
            return;
        }
        this.networkValid = true;
    }

    onSubmit() {

        this.checkNetwork();
        if (this.networkValid === false) {
            return;
        }
        this.isSubmitGoing = true;
        this.ipPoolService.create(this.item).subscribe(res => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
