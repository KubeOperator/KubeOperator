import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Plan} from '../plan';
import {PlanService} from '../plan.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-plan-detail',
    templateUrl: './plan-detail.component.html',
    styleUrls: ['./plan-detail.component.css']
})
export class PlanDetailComponent extends BaseModelDirective<Plan> implements OnInit {

    opened = false;
    item: Plan = new Plan();
    @Output() detail = new EventEmitter();

    constructor(private planService: PlanService, private translateService: TranslateService) {
        super(planService);
    }

    ngOnInit(): void {
    }

    open(item) {
        this.item = item;
        this.opened = true;
    }

    cancel() {
        this.opened = false;
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
