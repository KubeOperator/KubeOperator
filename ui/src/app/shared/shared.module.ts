import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CoreModule} from '../core/core.module';
import {ModalAlertComponent} from './common-component/modal-alert/modal-alert.component';
import {K8sPaginationComponent} from './common-component/k8s-pagination/k8s-pagination.component';
import { FilterComponent } from './common-component/filter/filter.component';
import {ZoneStatusPipe} from './pipe/zone-status.pipe';
import {CommonStatusPipe} from './pipe/common-status.pipe';
import {BackupAccountStatusPipe} from './pipe/backup-account-status.pipe';
import { UserTypePipe } from './pipe/user-type.pipe';
import { EmailShowPipe } from './pipe/email-show.pipe';
import { MessageTypePipe } from './pipe/message-type.pipe';


@NgModule({
    declarations: [ModalAlertComponent, K8sPaginationComponent, ZoneStatusPipe, CommonStatusPipe, BackupAccountStatusPipe, UserTypePipe,EmailShowPipe, MessageTypePipe,FilterComponent
    ],
    exports: [
        ModalAlertComponent,
        K8sPaginationComponent,
        ZoneStatusPipe,
        CommonStatusPipe,
        BackupAccountStatusPipe,
        UserTypePipe,
        EmailShowPipe,
        MessageTypePipe,
        FilterComponent
    ],
    imports: [
        CommonModule,
        CoreModule,
    ]
})
export class SharedModule {
}
