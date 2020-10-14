import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {VmConfigComponent} from './vm-config.component';
import { VmConfigListComponent } from './vm-config-list/vm-config-list.component';
import { VmConfigCreateComponent } from './vm-config-create/vm-config-create.component';
import { VmConfigDeleteComponent } from './vm-config-delete/vm-config-delete.component';
import { VmConfigUpdateComponent } from './vm-config-update/vm-config-update.component';
import {CoreModule} from '../../../core/core.module';
import {SharedModule} from '../../../shared/shared.module';


@NgModule({
    declarations: [VmConfigComponent, VmConfigListComponent, VmConfigCreateComponent, VmConfigDeleteComponent, VmConfigUpdateComponent],
    imports: [
        CommonModule,
        CoreModule,
        SharedModule
    ]
})
export class VmConfigModule {
}
