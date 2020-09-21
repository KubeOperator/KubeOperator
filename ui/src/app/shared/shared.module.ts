import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CoreModule} from '../core/core.module';
import {ModalAlertComponent} from './common-component/modal-alert/modal-alert.component';
import {K8sPaginationComponent} from './common-component/k8s-pagination/k8s-pagination.component';
import {ZoneStatusPipe} from './pipe/zone-status.pipe';
import {MenuAuthDirective} from './directive/menu-auth.directive';
import {OperateAuthDirective} from './directive/operate-auth.directive';
import {CommonStatusPipe} from './pipe/common-status.pipe';
import {BackupAccountStatusPipe} from './pipe/backup-account-status.pipe';
import { UserTypePipe } from './pipe/user-type.pipe';
import { MessageTypePipe } from './pipe/message-type.pipe';


@NgModule({
    declarations: [ModalAlertComponent, K8sPaginationComponent, ZoneStatusPipe, MenuAuthDirective,
        OperateAuthDirective, CommonStatusPipe, BackupAccountStatusPipe, UserTypePipe, MessageTypePipe
    ],
    exports: [
        ModalAlertComponent,
        K8sPaginationComponent,
        ZoneStatusPipe,
        MenuAuthDirective,
        OperateAuthDirective,
        CommonStatusPipe,
        BackupAccountStatusPipe,
        UserTypePipe,
        MessageTypePipe
    ],
    imports: [
        CommonModule,
        CoreModule,
    ]
})
export class SharedModule {
}
