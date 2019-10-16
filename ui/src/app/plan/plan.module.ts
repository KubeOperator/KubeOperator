import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {PlanComponent} from './plan.component';
import {PlanListComponent} from './plan-list/plan-list.component';
import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';

import { PlanDetailComponent } from './plan-detail/plan-detail.component';
import { PlanCreateComponent } from './plan-create/plan-create.component';

@NgModule({
  declarations: [PlanComponent, PlanListComponent, PlanDetailComponent, PlanCreateComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
  ]
})
export class PlanModule {
}
