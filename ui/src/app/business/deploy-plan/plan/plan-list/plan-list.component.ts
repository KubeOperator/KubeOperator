import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Plan} from '../plan';
import {PlanService} from '../plan.service';

@Component({
    selector: 'app-plan-list',
    templateUrl: './plan-list.component.html',
    styleUrls: ['./plan-list.component.css']
})
export class PlanListComponent extends BaseModelComponent<Plan> implements OnInit {

    constructor(private planService: PlanService) {
        super(planService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
