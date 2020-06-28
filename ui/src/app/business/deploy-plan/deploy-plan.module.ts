import {NgModule} from '@angular/core';
import {RegionModule} from './region/region.module';
import {RouterModule} from '@angular/router';
import {DeployPlanComponent} from './deploy-plan.component';
import {CoreModule} from '../../core/core.module';
import {ZoneModule} from './zone/zone.module';


@NgModule({
    declarations: [DeployPlanComponent],
    imports: [
        ZoneModule,
        RegionModule,
        RouterModule,
        CoreModule,
    ]
})
export class DeployPlanModule {
}
