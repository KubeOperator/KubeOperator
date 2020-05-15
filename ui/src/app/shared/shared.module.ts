import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {DatagridComponent} from './component/datagrid/datagrid.component';
import {CoreModule} from '../core/core.module';


@NgModule({
    declarations: [DatagridComponent],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class SharedModule {
}
