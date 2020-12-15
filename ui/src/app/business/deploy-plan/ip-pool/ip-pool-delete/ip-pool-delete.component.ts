import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {IpPool} from '../ip-pool';
import {IpPoolService} from '../ip-pool.service';
import {Project} from '../../../project/project';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-ip-pool-delete',
    templateUrl: './ip-pool-delete.component.html',
    styleUrls: ['./ip-pool-delete.component.css']
})
export class IpPoolDeleteComponent extends BaseModelDirective<IpPool> implements OnInit {


    opened = false;
    isSubmitGoing = false;
    items: IpPool[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private ipPoolService: IpPoolService,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(ipPoolService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.opened = true;
        this.items = items;
    }

    onCancel() {
        this.opened = false;
        this.isSubmitGoing = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.ipPoolService.batch('delete', this.items).subscribe(res => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.deleted.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
            this.opened = false;
        });
    }
}
