import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ZoneComponent} from './zone.component';
import {CoreModule} from '../core/core.module';
import {TipModule} from '../tip/tip.module';
import {SharedModule} from '../shared/shared.module';
import { ZoneListComponent } from './zone-list/zone-list.component';
import { ZoneCreateComponent } from './zone-create/zone-create.component';

@NgModule({
  declarations: [ZoneComponent, ZoneListComponent, ZoneCreateComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule,
    SharedModule
  ]
})
export class ZoneModule {
}
