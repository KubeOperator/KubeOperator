import {NgModule} from '@angular/core';
import {SystemLogComponent} from './system-log.component';
import {CoreModule} from "../../core/core.module";
import {SharedModule} from "../../shared/shared.module";

@NgModule({
    declarations: [
        SystemLogComponent
    ],
    imports: [
        CoreModule,
        SharedModule
    ]
})
export class SystemLogModule {
}