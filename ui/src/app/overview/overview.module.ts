import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {OverviewComponent} from './overview.component';
import {CoreModule} from '../core/core.module';
import {DescribeComponent} from './describe/describe.component';
import {ClusterStatusComponent} from './cluster-status/cluster-status.component';
import {UpgradeComponent} from './upgrade/upgrade.component';
import {ScaleComponent} from './scale/scale.component';
import {SharedModule} from '../shared/shared.module';
import {WebkubectlComponent} from './webkubectl/webkubectl.component';
import {AddWorkerComponent} from './add-worker/add-worker.component';
import {RemoveWorkerComponent} from './remove-worker/remove-worker.component';
import {NgCircleProgressModule} from "ng-circle-progress";

@NgModule({
  declarations: [OverviewComponent, DescribeComponent, ClusterStatusComponent, UpgradeComponent, ScaleComponent,
    WebkubectlComponent, AddWorkerComponent, RemoveWorkerComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
    NgCircleProgressModule
  ]
})
export class OverviewModule {
}
