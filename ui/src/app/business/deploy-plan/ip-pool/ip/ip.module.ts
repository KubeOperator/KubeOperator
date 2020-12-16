import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {IpComponent} from './ip.component';
import {IpCreateComponent} from './ip-create/ip-create.component';
import {IpDeleteComponent} from './ip-delete/ip-delete.component';
import {IpListComponent} from './ip-list/ip-list.component';
import {CoreModule} from '../../../../core/core.module';


@NgModule({
    declarations: [IpComponent, IpCreateComponent, IpDeleteComponent, IpListComponent],
    exports: [
        IpListComponent
    ],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class IpModule {
}
