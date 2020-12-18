import {Component, OnInit} from '@angular/core';
import {BaseModelDirective} from '../../../../../shared/class/BaseModelDirective';
import {Ip, IpUpdate} from '../ip';
import {IpService} from '../ip.service';
import {ActivatedRoute, Router} from '@angular/router';
import {IpPool} from '../../ip-pool';
import {ModalAlertService} from '../../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../layout/common-alert/alert';

@Component({
    selector: 'app-ip-list',
    templateUrl: './ip-list.component.html',
    styleUrls: ['./ip-list.component.css']
})
export class IpListComponent extends BaseModelDirective<Ip> implements OnInit {

    ipPoolName: string;
    ipPool: IpPool;

    constructor(private ipService: IpService,
                private router: Router,
                private route: ActivatedRoute,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(ipService);
    }

    ngOnInit(): void {
        this.route.data.subscribe(data => {
            this.ipPool = data.ipPool;
            this.ipPoolName = data.ipPool.name;
            this.refresh();
        });
    }

    refresh() {
        this.loading = true;
        this.ipService.page(this.page, this.size, this.ipPoolName).subscribe(data => {
            this.items = data.items;
            this.total = data.total;
            this.loading = false;
        });
    }

    update(item: Ip, operation: string) {
        const update = new IpUpdate();
        update.address = item.address;
        update.operation = operation;
        this.ipService.update(update.address, update, this.ipPoolName).subscribe(data => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
            this.refresh();
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
