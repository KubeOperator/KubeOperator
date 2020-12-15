import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {IpPool, IpPoolCreate} from '../ip-pool';
import {IpPoolService} from '../ip-pool.service';
import {NamePattern} from '../../../../constant/pattern';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

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
        this.item = new IpPoolCreate();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
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
