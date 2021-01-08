import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../../shared/class/BaseModelDirective';
import {Ip, IpCreate} from '../ip';
import {IpService} from '../ip.service';
import {ModalAlertService} from '../../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {ActivatedRoute, Router} from '@angular/router';
import {AlertLevels} from '../../../../../layout/common-alert/alert';
import {IpPool} from '../../ip-pool';
import * as ipaddr from 'ipaddr.js';


@Component({
    selector: 'app-ip-create',
    templateUrl: './ip-create.component.html',
    styleUrls: ['./ip-create.component.css']
})
export class IpCreateComponent extends BaseModelDirective<Ip> implements OnInit {

    opened = false;
    isSubmitGoing = false;
    ipPool: IpPool = new IpPool();
    item: IpCreate = new IpCreate();
    @Output() created = new EventEmitter();
    networkValid = false;

    constructor(private ipService: IpService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService,
                private router: Router,
                private route: ActivatedRoute) {
        super(ipService);
    }

    ngOnInit(): void {
        this.route.data.subscribe(data => {
            this.ipPool = data.ipPool;
        });
    }

    open() {
        this.opened = true;
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
        this.item.ipPoolName = this.ipPool.name;
        this.item.subnet = this.ipPool.subnet;
        this.ipService.create(this.item, this.ipPool.name).subscribe(res => {
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
