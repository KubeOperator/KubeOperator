import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ZoneCreateComponent} from './zone-create/zone-create.component';
import {ZoneDeleteComponent} from './zone-delete/zone-delete.component';
import {ZoneUpdateComponent} from './zone-update/zone-update.component';
import {ZoneListComponent} from './zone-list/zone-list.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';
import {ZoneComponent} from './zone.component';
import { ZoneDetailComponent } from './zone-detail/zone-detail.component';


@NgModule({
    declarations: [ZoneComponent, ZoneCreateComponent, ZoneDeleteComponent, ZoneUpdateComponent, ZoneListComponent, ZoneDetailComponent],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class ZoneModule {
}
