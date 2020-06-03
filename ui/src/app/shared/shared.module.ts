import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CoreModule} from '../core/core.module';
import {ModalAlertComponent} from './common-component/modal-alert/modal-alert.component';


@NgModule({
    declarations: [ModalAlertComponent,],
    exports: [
        ModalAlertComponent
    ],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class SharedModule {
}
