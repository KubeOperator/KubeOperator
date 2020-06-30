import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {PlanCreateComponent} from './plan-create/plan-create.component';
import {PlanDeleteComponent} from './plan-delete/plan-delete.component';
import {PlanListComponent} from './plan-list/plan-list.component';
import {PlanComponent} from './plan.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';


@NgModule({
    declarations: [PlanCreateComponent, PlanDeleteComponent, PlanListComponent, PlanComponent],
    exports: [
        PlanListComponent
    ],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class PlanModule {
}
