import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Plan} from '../plan';
import {PlanService} from '../plan.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-plan-list',
    templateUrl: './plan-list.component.html',
    styleUrls: ['./plan-list.component.css']
})
export class PlanListComponent extends BaseModelDirective<Plan> implements OnInit {


    @Output() detailEvent = new EventEmitter<Plan>();


    constructor(private planService: PlanService, private translateService: TranslateService) {
        super(planService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onDetail(item) {
        this.detailEvent.emit(item);
    }

    getDeployName(name: string) {
        switch (name) {
            case 'SINGLE':
                return this.translateService.instant('APP_PLAN_DEPLOY_TEMPLATE_SINGLE');
            case 'MULTIPLE':
                return this.translateService.instant('APP_PLAN_DEPLOY_TEMPLATE_MULTIPLE');
            default:
                return 'None';
        }
    }
}
