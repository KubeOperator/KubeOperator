import {NgModule} from '@angular/core';
import {RegionModule} from './region/region.module';
import {RouterModule} from '@angular/router';
import {DeployPlanComponent} from './deploy-plan.component';
import {CoreModule} from '../../core/core.module';


@NgModule({
    declarations: [DeployPlanComponent],
    imports: [
        RegionModule,
        RouterModule,
        CoreModule,
    ]
})
export class DeployPlanModule {
}
