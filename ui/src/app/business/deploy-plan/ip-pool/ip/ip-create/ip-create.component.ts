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

    onSubmit() {
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
