import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Plan} from '../plan';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {PlanService} from '../plan.service';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-plan-delete',
    templateUrl: './plan-delete.component.html',
    styleUrls: ['./plan-delete.component.css']
})
export class PlanDeleteComponent extends BaseModelDirective<Plan> implements OnInit {

    opened = false;
    items: Plan[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private planService: PlanService, private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService, private translateService: TranslateService) {
        super(planService);
    }

    ngOnInit(): void {
    }

    open(items) {
        this.items = items;
        this.opened = true;
    }


    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.planService.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.opened = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
