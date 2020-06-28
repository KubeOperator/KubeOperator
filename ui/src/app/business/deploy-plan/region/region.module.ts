import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RegionComponent} from './region.component';
import {RegionCreateComponent} from './region-create/region-create.component';
import {RegionListComponent} from './region-list/region-list.component';
import {RegionDeleteComponent} from './region-delete/region-delete.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';


@NgModule({
    declarations: [RegionComponent, RegionCreateComponent, RegionListComponent, RegionDeleteComponent],
    exports: [
        RegionListComponent
    ],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class RegionModule {
}
