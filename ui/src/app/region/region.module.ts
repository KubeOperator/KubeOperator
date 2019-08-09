import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {RegionListComponent} from './region-list/region-list.component';
import {RegionDetailComponent} from './region-detail/region-detail.component';
import {RegionCreateComponent} from './region-create/region-create.component';
import {CoreModule} from '../core/core.module';
import {TipModule} from '../tip/tip.module';
import {RegionComponent} from './region.component';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [RegionComponent, RegionListComponent, RegionDetailComponent, RegionCreateComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule,
    SharedModule
  ],
})
export class RegionModule {
}
