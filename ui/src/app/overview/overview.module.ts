import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {OverviewComponent} from './overview.component';
import {CoreModule} from '../core/core.module';
import {DescribeComponent} from './describe/describe.component';
import {SharedModule} from '../shared/shared.module';
import {ClusterStatusComponent} from './cluster-status/cluster-status.component';
import { UpgradeComponent } from './upgrade/upgrade.component';

@NgModule({
  declarations: [OverviewComponent, DescribeComponent, ClusterStatusComponent, UpgradeComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule
  ]
})
export class OverviewModule {
}
