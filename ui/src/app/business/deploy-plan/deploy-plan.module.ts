import {NgModule} from '@angular/core';
import {RegionModule} from './region/region.module';
import {RouterModule} from '@angular/router';
import {DeployPlanComponent} from './deploy-plan.component';
import {CoreModule} from '../../core/core.module';
import { ZoneComponent } from './zone/zone.component';
import {ZoneModule} from './zone/zone.module';


@NgModule({
    declarations: [DeployPlanComponent, ZoneComponent],
    imports: [
        RegionModule,
        RouterModule,
        CoreModule,
        ZoneModule,
    ]
})
export class DeployPlanModule {
}
