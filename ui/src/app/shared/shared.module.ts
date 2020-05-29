import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CoreModule} from '../core/core.module';
import {ModalAlertComponent} from './common-component/modal-alert/modal-alert.component';
import { CommonAlertComponent } from './common-component/common-alert/common-alert.component';


@NgModule({
    declarations: [ModalAlertComponent, CommonAlertComponent],
    exports: [
        CommonAlertComponent
    ],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class SharedModule {
}
