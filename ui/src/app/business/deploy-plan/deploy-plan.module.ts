import {NgModule} from '@angular/core';
import {RegionModule} from './region/region.module';
import {RouterModule} from '@angular/router';
import {DeployPlanComponent} from './deploy-plan.component';
import {CoreModule} from '../../core/core.module';
import {ZoneModule} from './zone/zone.module';
import {PlanModule} from './plan/plan.module';
import {VmConfigModule} from './vm-config/vm-config.module';
import {IpPoolModule} from './ip-pool/ip-pool.module';


@NgModule({
    declarations: [DeployPlanComponent],
    imports: [
        ZoneModule,
        RegionModule,
        RouterModule,
        CoreModule,
        PlanModule,
        VmConfigModule,
        IpPoolModule
    ]
})
export class DeployPlanModule {
}
