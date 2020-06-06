import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CoreModule} from '../core/core.module';
import {ModalAlertComponent} from './common-component/modal-alert/modal-alert.component';
import { K8sPaginationComponent } from './common-component/k8s-pagination/k8s-pagination.component';


@NgModule({
    declarations: [ModalAlertComponent, K8sPaginationComponent,],
    exports: [
        ModalAlertComponent,
        K8sPaginationComponent
    ],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class SharedModule {
}
